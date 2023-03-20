package mail

import (
	"strings"

	"github.com/emersion/go-imap"
	imapclient "github.com/emersion/go-imap/client"
	"github.com/google/uuid"
	"github.com/tauraamui/maildew/internal/kvs"
)

type RemoteConnection interface {
	RemoteMailboxLister
	RemoteMessagesFetcher
}

type listFunc func(ref, name string, ch chan *imap.MailboxInfo) error

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
		return parts[1]
	}
	return ""
}

var acquireClientConn = func(addr string, acc Account) (RemoteConnection, error) {

	cc, err := imapclient.Dial(addr)
	if err != nil {
		return nil, err
	}

	if err := cc.Login(acc.Username, acc.Password); err != nil {
		return nil, err
	}

	return cc, nil
}

func RegisterAccount(
	accRepo AccountRepo,
	mbRepo MailboxRepo,
	msgRepo MessageRepo,
	acc Account,
) error {
	addr := resolveAddressFromUsername(acc.Username)

	if err := persistAccount(accRepo, acc); err != nil {
		return err
	}

	cc, err := acquireClientConn(addr, acc)
	if err != nil {
		return err
	}

	if err := syncAccountMailboxes(cc, mbRepo, acc); err != nil {
		return err
	}

	return nil
}

func syncAccountMailboxes(conn RemoteConnection, mbRepo MailboxRepo, acc Account) error {
	mailboxes := make(chan *imap.MailboxInfo, 10)

	if err := persistMailboxes(conn.List, mailboxes, mbRepo, acc); err != nil {
		return err
	}

	return nil
}

func persistAccount(ar AccountRepo, acc Account) error {
	acc.UUID = uuid.New()
	return ar.Save(acc)
}

func persistMailboxes(lister listFunc, mailboxes chan *imap.MailboxInfo, mbRepo MailboxRepo, acc Account) error {
	done := make(chan error, 1)
	go func() {
		done <- lister("", "*", mailboxes)
	}()

	for mb := range mailboxes {
		persistMailbox(mbRepo, acc.UUID, Mailbox{Name: mb.Name})
	}

	if err := <-done; err != nil {
		return err
	}

	return nil
}

func persistMailbox(mr MailboxRepo, owner kvs.UUID, mb Mailbox) (kvs.UUID, error) {
	mb.UUID = uuid.New()
	return mb.UUID, mr.Save(owner, mb)
}

func storeMessage(msgr MessageRepo, owner kvs.UUID, msg Message) (kvs.UUID, error) {
	msg.UUID = uuid.New()
	return msg.UUID, msgr.Save(owner, msg)
}

func listMailboxes(lister RemoteMailboxLister, dest chan *imap.MailboxInfo, errch chan<- error) {
	errch <- lister.List("", "*", dest)
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
