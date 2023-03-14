package mail

import (
	"testing"

	"github.com/emersion/go-imap"
	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/kvs"
)

type mockRemoteConnection struct{}

func (mc mockRemoteConnection) List(ref, name string, ch chan *imap.MailboxInfo) error {
	close(ch)
	return nil
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

type mockMailboxRepo struct{}

func (mmr mockMailboxRepo) Save(owner kvs.UUID, mailbox Mailbox) error {
	return nil
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
	mconn := mockRemoteConnection{}
	reset := overloadAcquireClientConn(func(s string, a Account) (RemoteConnection, error) {
		return mconn, nil
	})
	defer reset()

	accRepo := mockAccountRepo{}
	mbRepo := mockMailboxRepo{}
	msgRepo := mockMessageRepo{}

	is := is.New(t)

	is.NoErr(RegisterAccount(accRepo, mbRepo, msgRepo, Account{Username: "test@place.com", Password: "efewfweoifjio"}))
}
