package mail

import "github.com/tauraamui/maildew/internal/kvs"

type Mailbox interface {
	Name() string
	FetchAllMessages() ([]Message, error)
	FetchAllMessageUIDs() ([]MessageUID, error)
}

type mailbox struct {
	db      kvs.DB
	mf      messageFetcher
	account Account
	name    string
}

func newMailbox(db kvs.DB, name string, owner Account, mf messageFetcher) Mailbox {
	return mailbox{db, mf, owner, name}
}

func (m mailbox) Name() string {
	return m.name
}

func (m mailbox) FetchAllMessages() ([]Message, error) {
	// TODO:(tauraamui) here we should store/cache mailboxes to
	//                  a prefix key set in the K/V DB
	return m.mf.fetchAllMessages(m.account, m)
}

func (m mailbox) FetchAllMessageUIDs() ([]MessageUID, error) {
	return m.mf.fetchAllMessageUIDs(m.account, m)
}
