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

func (mmsgr *mockMessageRepo) FetchByOwner(owner kvs.UUID) ([]Message, error) {
	return nil, nil
}

func (mmsgr *mockMessageRepo) DumpTo(w io.Writer) error {
	return mmsgr.err
}

func (mmsgr mockMessageRepo) Close() error { return nil }

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
		"WORK":     {},
		"SHOPPING": {},
		"SPAM":     {},
		"DRAFTS":   {},
		"OTHER":    {},
		"MISC":     {},
		"SCHOOL":   {},
		"LIBRARY":  {},
		"MISC2":    {},
		"MISC3":    {},
		"MISC4":    {},
	}
}

func TestRegisterAccountSuccessAgainstRealKVSInstance(t *testing.T) {
	fakeStdout := strings.Builder{}
	log := logging.New(logging.Options{Level: logging.DEBUG, Writer: &fakeStdout})
	is := is.New(t)

	mconn := &mockRemoteConnection{
		mailboxes: makeRemoteConnectionData(),
	}
	reset := overloadAcquireClientConn(func(s string, a Account, ssl bool) (RemoteConnection, error) {
		return mconn, nil
	})
	defer reset()

	db, err := kvs.NewMemDB()
	is.NoErr(err)

	accRepo := NewAccountRepo(db)
	mbRepo := NewMailboxRepo(db)

	acc := Account{Username: "test@place.com", Password: "efewfweoifjio"}
	is.NoErr(RegisterAccount(log, "", accRepo, mbRepo, &acc))
	mboxes, err := mbRepo.FetchByOwner(acc.UUID)
	is.NoErr(err)

	is.Equal(len(mboxes), len(mconn.mailboxes))
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

	is := is.New(t)

	is.NoErr(RegisterAccount(log, "", accRepo, &mbRepo, &Account{Username: "test@place.com", Password: "efewfweoifjio"}))

	is.Equal(len(mbRepo.saved), 12)
	is = is.NewRelaxed(t)
	is.Equal(mbRepo.saved[0].mb.Name, "DRAFTS")
	is.Equal(mbRepo.saved[1].mb.Name, "INBOX")
	is.Equal(mbRepo.saved[2].mb.Name, "LIBRARY")
	is.Equal(mbRepo.saved[4].mb.Name, "MISC2")
	is.Equal(mbRepo.saved[5].mb.Name, "MISC3")
	is.Equal(mbRepo.saved[6].mb.Name, "MISC4")
	is.Equal(mbRepo.saved[7].mb.Name, "OTHER")
	is.Equal(mbRepo.saved[8].mb.Name, "SCHOOL")
	is.Equal(mbRepo.saved[9].mb.Name, "SHOPPING")
	is.Equal(mbRepo.saved[10].mb.Name, "SPAM")
	is.Equal(mbRepo.saved[11].mb.Name, "WORK")
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

	is := is.NewRelaxed(t)

	err := RegisterAccount(log, "", accRepo, &mbRepo, &Account{Username: "test@place.com", Password: "efewfweoifjio"})
	is.Equal(err.Error(), "failed to acquire next mailbox")

	is = is.New(t)
	is.Equal(len(mbRepo.saved), 1)
	is.Equal(mbRepo.saved[0].mb.Name, "DRAFTS")
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

	is := is.New(t)

	err := RegisterAccount(log, "", accRepo, &mbRepo, &Account{Username: "test@place.com", Password: "efewfweoifjio"})
	is.True(err != nil)
	is.Equal(err.Error(), "failed to persist mailbox")

	is.Equal(len(mbRepo.saved), 1)

	is.Equal(mbRepo.saved[0].mb.Name, "DRAFTS")
}
