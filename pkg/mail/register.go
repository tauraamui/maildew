package mail

import (
	"fmt"
	"strings"
	"sync"

	"github.com/emersion/go-imap"
	imapclient "github.com/emersion/go-imap/client"
	"github.com/google/uuid"
	"github.com/tauraamui/maildew/internal/kvs"
	"github.com/tauraamui/maildew/pkg/logging"
)

type RemoteConnection interface {
	RemoteMailboxLister
	RemoteMessagesFetcher
	Close() error
}

type RemoteMailboxLister interface {
	List(ref, name string, ch chan *imap.MailboxInfo) error
}

type RemoteMessagesFetcher interface {
	Select(name string, readOnly bool) (*imap.MailboxStatus, error)
	Fetch(seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error
}

type Account struct {
	UUID               kvs.UUID
	Username, Password string
}

type Mailbox struct {
	UUID kvs.UUID
	Name string
}

type Message struct {
	UUID      kvs.UUID
	RemoteUID uint32
}

func resolveAddressFromUsername(username string) string {
	parts := strings.Split(username, "@")
	if len(parts) > 1 {
		return fmt.Sprintf("imap.%s:%d", parts[1], 993)
	}
	return ""
}

var acquireClientConn = func(addr string, acc Account, useSSL bool) (RemoteConnection, error) {
	var cc *imapclient.Client
	var err error

	if useSSL {
		cc, err = imapclient.DialTLS(addr, nil)
	} else {
		cc, err = imapclient.Dial(addr)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to dial to address %s: %w", addr, err)
	}

	if err := cc.Login(acc.Username, acc.Password); err != nil {
		return nil, fmt.Errorf("failed to login to account: %w", err)
	}

	return cc, nil
}

func RegisterAccount(
	log logging.I,
	addr string,
	accRepo AccountRepo,
	mbRepo MailboxRepo,
	msgRepo MessageRepo,
	acc *Account,
) error {

	useSSL := false
	if len(addr) == 0 {
		addr = resolveAddressFromUsername(acc.Username)
		useSSL = true
	}

	log.Debug().Msgf("resolved addr to %s", addr)

	if err := persistAccount(accRepo, acc); err != nil {
		return err
	}

	log.Debug().Msg("attempting to login to account")
	cc, err := acquireClientConn(addr, *acc, useSSL)
	if err != nil {
		return err
	}
	defer cc.Close()
	log.Debug().Msg("logged into account")

	log.Debug().Msg("syncing mailboxes")
	if err := syncAccountMailboxes(cc, mbRepo, msgRepo, *acc); err != nil {
		return err
	}
	log.Debug().Msg("synced mailboxes")

	return nil
}

func syncAccountMailboxes(conn RemoteConnection, mbRepo MailboxRepo, msgRepo MessageRepo, acc Account) error {
	if err := persistMailboxes(conn, mbRepo, msgRepo, acc); err != nil {
		return err
	}

	return nil
}

func persistAccount(ar AccountRepo, acc *Account) error {
	acc.UUID = uuid.New()
	return ar.Save(*acc)
}

// FIX:(tauraamui)
// This function has a reading bug, in that it seems as though the channel is only
// read from once it's finished being filled to capacity, and it is never filled beyond that.
func persistMailboxes(conn RemoteConnection, mbRepo MailboxRepo, msgRepo MessageRepo, acc Account) error {
	done := make(chan error, 1)
	defer close(done)
	mailboxes := make(chan *imap.MailboxInfo, 10)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		done <- conn.List("", "*", mailboxes)
	}()

	for mb := range mailboxes {
		err := persistMailbox(conn, mbRepo, msgRepo, acc.UUID, Mailbox{Name: mb.Name})
		if err != nil {
			return err
		}
	}

	wg.Wait()

	if err := <-done; err != nil {
		return err
	}

	return nil
}

func persistMailbox(conn RemoteConnection, mbRepo MailboxRepo, msgRepo MessageRepo, owner kvs.UUID, mb Mailbox) error {
	mb.UUID = uuid.New()
	if err := mbRepo.Save(owner, mb); err != nil {
		return err
	}

	done := make(chan error)
	defer close(done)
	messages := make(chan *imap.Message)

	go fetchMailboxMessages(conn, mb.Name, messages, done)

	for msg := range messages {
		_, err := storeMessage(msgRepo, mb.UUID, Message{
			RemoteUID: msg.Uid,
		})
		if err != nil {
			return err
		}
	}

	if err := <-done; err != nil {
		return err
	}

	return nil
}

func storeMessage(msgr MessageRepo, owner kvs.UUID, msg Message) (kvs.UUID, error) {
	msg.UUID = uuid.New()
	return msg.UUID, msgr.Save(owner, msg)
}

func fetchMailboxMessages(fetcher RemoteMessagesFetcher, mbName string, dest chan *imap.Message, errch chan<- error) {
	mb, err := fetcher.Select(mbName, true)
	if err != nil {
		errch <- err
		return
	}

	from := uint32(1)
	to := mb.Messages
	seqset := imap.SeqSet{}
	seqset.AddRange(from, to)

	errch <- fetcher.Fetch(&seqset, []imap.FetchItem{imap.FetchEnvelope}, dest)
}
