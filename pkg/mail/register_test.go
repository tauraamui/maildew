package mail

import (
	"errors"
	"testing"

	"github.com/emersion/go-imap"
	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/kvs"
)

type mockRemoteConnection struct {
	mailboxes         []string
	returnErrAfterNum int
	err               error
}

func (mc mockRemoteConnection) List(ref, name string, ch chan *imap.MailboxInfo) error {
	defer close(ch)
	for i, name := range mc.mailboxes {
		if mc.returnErrAfterNum > 0 && i >= mc.returnErrAfterNum {
			return mc.err
		}
		ch <- &imap.MailboxInfo{Name: name}
	}
	return mc.err
}

func (mc mockRemoteConnection) Select(name string, readOnly bool) (*imap.MailboxStatus, error) {
	return nil, nil
}

func (mc mockRemoteConnection) Fetch(seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
	close(ch)
	return nil
}

type mockAccountRepo struct{}

func (mar mockAccountRepo) Save(user Account) error {
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

func (mmr mockMailboxRepo) Close() {}

type mockMessageRepo struct{}

func (mmsgr mockMessageRepo) Save(owner kvs.UUID, msg Message) error {
	return nil
}

func (mmsgr mockMessageRepo) Close() {}

func overloadAcquireClientConn(overload func(string, Account) (RemoteConnection, error)) func() {
	acquireClientConnRef := acquireClientConn
	acquireClientConn = overload
	return func() {
		acquireClientConn = acquireClientConnRef
	}
}

func TestRegisterAccountSuccessSyncedRemoteMailboxes(t *testing.T) {
	mconn := mockRemoteConnection{
		mailboxes: []string{"INBOX", "SPAM"},
	}
	reset := overloadAcquireClientConn(func(s string, a Account) (RemoteConnection, error) {
		return mconn, nil
	})
	defer reset()

	accRepo := mockAccountRepo{}
	mbRepo := mockMailboxRepo{}
	msgRepo := mockMessageRepo{}

	is := is.New(t)

	is.NoErr(RegisterAccount(accRepo, &mbRepo, msgRepo, Account{Username: "test@place.com", Password: "efewfweoifjio"}))

	is.Equal(len(mbRepo.saved), 2)
	is.Equal(mbRepo.saved[0].mb.Name, "INBOX")
	is.Equal(mbRepo.saved[1].mb.Name, "SPAM")
}

func TestRegisterAccountErrorDuringListingMailboxes(t *testing.T) {
	mconn := mockRemoteConnection{
		mailboxes:         []string{"INBOX", "SPAM"},
		returnErrAfterNum: 1,
		err:               errors.New("failed to acquire next mailbox"),
	}
	reset := overloadAcquireClientConn(func(s string, a Account) (RemoteConnection, error) {
		return mconn, nil
	})
	defer reset()

	accRepo := mockAccountRepo{}
	mbRepo := mockMailboxRepo{}
	msgRepo := mockMessageRepo{}

	is := is.NewRelaxed(t)

	err := RegisterAccount(accRepo, &mbRepo, msgRepo, Account{Username: "test@place.com", Password: "efewfweoifjio"})
	is.Equal(err.Error(), "failed to acquire next mailbox")

	is = is.New(t)
	is.Equal(len(mbRepo.saved), 1)
	is.Equal(mbRepo.saved[0].mb.Name, "INBOX")
}

func TestRegisterAccountErrorDuringStoringMailboxes(t *testing.T) {
	mconn := mockRemoteConnection{
		mailboxes: []string{"INBOX", "SPAM"},
	}
	reset := overloadAcquireClientConn(func(s string, a Account) (RemoteConnection, error) {
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

	err := RegisterAccount(accRepo, &mbRepo, msgRepo, Account{Username: "test@place.com", Password: "efewfweoifjio"})
	is.Equal(err.Error(), "failed to persist mailbox")

	is = is.New(t)
	is.Equal(len(mbRepo.saved), 1)
	is.Equal(mbRepo.saved[0].mb.Name, "INBOX")
}
