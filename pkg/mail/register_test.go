package mail

import (
	"errors"
	"io"
	"sort"
	"strings"
	"testing"

	"github.com/emersion/go-imap"
	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/kvs"
	"github.com/tauraamui/maildew/pkg/logging"
)

type mockRemoteConnection struct {
	selected          string
	mailboxes         map[string][]*imap.Message
	returnErrAfterNum int
	err               error
}

func sortedKeys(m map[string][]*imap.Message) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (mc mockRemoteConnection) List(ref, name string, ch chan *imap.MailboxInfo) error {
	defer close(ch)

	for i, name := range sortedKeys(mc.mailboxes) {
		if mc.returnErrAfterNum > 0 && i >= mc.returnErrAfterNum {
			return mc.err
		}
		ch <- &imap.MailboxInfo{Name: name}
	}
	return mc.err
}

func (mc *mockRemoteConnection) Select(name string, readOnly bool) (*imap.MailboxStatus, error) {
	mc.selected = name
	return &imap.MailboxStatus{
		Messages: uint32(len(mc.mailboxes[mc.selected])),
	}, nil
}

func (mc mockRemoteConnection) Fetch(seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
	for _, msg := range mc.mailboxes[mc.selected] {
		ch <- msg
	}
	close(ch)
	return nil
}

func (mc mockRemoteConnection) Close() error { return nil }

type mockAccountRepo struct{}

func (mar mockAccountRepo) Save(user Account) error {
	return nil
}

func (mar mockAccountRepo) DumpTo(w io.Writer) error {
	return nil
}

func (mar mockAccountRepo) Close() {}

type mockMailboxRepo struct {
	saved                  []savedMailbox
	returnErrAfterNumSaved int
	savedNum               int
	err                    error
}

type savedMailbox struct {
	owner kvs.UUID
	mb    Mailbox
}

func (mmr *mockMailboxRepo) Save(owner kvs.UUID, mailbox Mailbox) error {
	defer func() { mmr.savedNum++ }()
	if mmr.returnErrAfterNumSaved > 0 && mmr.savedNum >= mmr.returnErrAfterNumSaved {
		return mmr.err
	}
	mmr.saved = append(mmr.saved, savedMailbox{owner, mailbox})
	return mmr.err
}

func (mmr *mockMailboxRepo) FetchByOwner(owner kvs.UUID) ([]Mailbox, error) {
	return nil, nil
}

func (mmr *mockMailboxRepo) DumpTo(w io.Writer) error {
	return nil
}

func (mmr mockMailboxRepo) Close() {}

type mockMessageRepo struct {
	saved    []Message
	savedNum int
	err      error
}

func (mmsgr *mockMessageRepo) Save(owner kvs.UUID, msg Message) error {
	defer func() { mmsgr.savedNum++ }()
	mmsgr.saved = append(mmsgr.saved, msg)
	return mmsgr.err
}

func (mmsgr *mockMessageRepo) DumpTo(w io.Writer) error {
	return nil
}

func (mmsgr mockMessageRepo) Close() {}

func overloadAcquireClientConn(overload func(string, Account, bool) (RemoteConnection, error)) func() {
	acquireClientConnRef := acquireClientConn
	acquireClientConn = overload
	return func() {
		acquireClientConn = acquireClientConnRef
	}
}

func makeRemoteConnectionData() map[string][]*imap.Message {
	return map[string][]*imap.Message{
		"INBOX": {
			{
				Uid: 321,
				Envelope: &imap.Envelope{
					Subject: "Test inbox message",
				},
			},
			{
				Uid: 5940,
				Envelope: &imap.Envelope{
					Subject: "Car insurance ad",
				},
			},
			{
				Uid: 623943,
				Envelope: &imap.Envelope{
					Subject: "Order is 15 days late",
				},
			},
			{
				Uid: 65096,
				Envelope: &imap.Envelope{
					Subject: "Feel happy!",
				},
			},
		},
		"INBOX0":  {},
		"INBOX1":  {},
		"INBOX2":  {},
		"INBOX3":  {},
		"INBOX4":  {},
		"INBOX5":  {},
		"INBOX6":  {},
		"INBOX7":  {},
		"INBOX8":  {},
		"INBOX9":  {},
		"INBOX10": {},
		"INBOX11": {},
		"SPAM":    {},
	}
}

func TestRegisterAccountSuccessSyncedRemoteMailboxes(t *testing.T) {
	fakeStdout := strings.Builder{}
	log := logging.New(logging.Options{Level: logging.DEBUG, Writer: &fakeStdout})
	mconn := &mockRemoteConnection{
		mailboxes: makeRemoteConnectionData(),
	}
	reset := overloadAcquireClientConn(func(s string, a Account, ssl bool) (RemoteConnection, error) {
		return mconn, nil
	})
	defer reset()

	accRepo := mockAccountRepo{}
	mbRepo := mockMailboxRepo{}
	msgRepo := mockMessageRepo{}

	is := is.New(t)

	is.NoErr(RegisterAccount(log, "", accRepo, &mbRepo, &msgRepo, &Account{Username: "test@place.com", Password: "efewfweoifjio"}))

	is.Equal(len(mbRepo.saved), 14)
	is = is.NewRelaxed(t)
	is.Equal(mbRepo.saved[0].mb.Name, "INBOX")
	is.Equal(mbRepo.saved[1].mb.Name, "INBOX0")
	is.Equal(mbRepo.saved[2].mb.Name, "INBOX1")
	is.Equal(mbRepo.saved[3].mb.Name, "INBOX10")
	is.Equal(mbRepo.saved[4].mb.Name, "INBOX11")
	is.Equal(mbRepo.saved[5].mb.Name, "INBOX2")
	is.Equal(mbRepo.saved[6].mb.Name, "INBOX3")
	is.Equal(mbRepo.saved[7].mb.Name, "INBOX4")
	is.Equal(mbRepo.saved[8].mb.Name, "INBOX5")
	is.Equal(mbRepo.saved[9].mb.Name, "INBOX6")
	is.Equal(mbRepo.saved[10].mb.Name, "INBOX7")
	is.Equal(mbRepo.saved[11].mb.Name, "INBOX8")
	is.Equal(mbRepo.saved[12].mb.Name, "INBOX9")
	is.Equal(mbRepo.saved[13].mb.Name, "SPAM")

	is = is.New(t)
	is.Equal(len(msgRepo.saved), 4)
	is = is.NewRelaxed(t)
	is.Equal(msgRepo.saved[0].RemoteUID, uint32(321))
	is.Equal(msgRepo.saved[1].RemoteUID, uint32(5940))
	is.Equal(msgRepo.saved[2].RemoteUID, uint32(623943))
	is.Equal(msgRepo.saved[3].RemoteUID, uint32(65096))
}

func TestRegisterAccountErrorDuringListingMailboxes(t *testing.T) {
	fakeStdout := strings.Builder{}
	log := logging.New(logging.Options{Level: logging.DEBUG, Writer: &fakeStdout})

	mconn := &mockRemoteConnection{
		mailboxes:         makeRemoteConnectionData(),
		returnErrAfterNum: 1,
		err:               errors.New("failed to acquire next mailbox"),
	}
	reset := overloadAcquireClientConn(func(s string, a Account, ssl bool) (RemoteConnection, error) {
		return mconn, nil
	})
	defer reset()

	accRepo := mockAccountRepo{}
	mbRepo := mockMailboxRepo{}
	msgRepo := mockMessageRepo{}

	is := is.NewRelaxed(t)

	err := RegisterAccount(log, "", accRepo, &mbRepo, &msgRepo, &Account{Username: "test@place.com", Password: "efewfweoifjio"})
	is.Equal(err.Error(), "failed to acquire next mailbox")

	is = is.New(t)
	is.Equal(len(mbRepo.saved), 1)
	is.Equal(mbRepo.saved[0].mb.Name, "INBOX")
}

func TestRegisterAccountErrorDuringStoringMailboxes(t *testing.T) {
	fakeStdout := strings.Builder{}
	log := logging.New(logging.Options{Level: logging.DEBUG, Writer: &fakeStdout})

	mconn := &mockRemoteConnection{
		mailboxes: makeRemoteConnectionData(),
	}
	reset := overloadAcquireClientConn(func(s string, a Account, ssl bool) (RemoteConnection, error) {
		return mconn, nil
	})
	defer reset()

	accRepo := mockAccountRepo{}
	mbRepo := mockMailboxRepo{
		returnErrAfterNumSaved: 1,
		err:                    errors.New("failed to persist mailbox"),
	}
	msgRepo := mockMessageRepo{}

	is := is.NewRelaxed(t)

	err := RegisterAccount(log, "", accRepo, &mbRepo, &msgRepo, &Account{Username: "test@place.com", Password: "efewfweoifjio"})
	is.Equal(err.Error(), "failed to persist mailbox")

	is = is.New(t)
	is.Equal(len(mbRepo.saved), 1)

	is.Equal(mbRepo.saved[0].mb.Name, "INBOX")
}
