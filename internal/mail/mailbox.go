package mail

import "github.com/tauraamui/maildew/internal/storage/models"

type Mailbox interface {
	Name() string
	FetchAllMessages() ([]Message, error)
	FetchAllMessageUIDs() ([]MessageUID, error)
}

type mailbox struct {
	mf      MessageFetcher
	account models.Account
	name    string
}

func NewMailbox(name string, owner models.Account, mf MessageFetcher) Mailbox {
	return newMailbox(name, owner, mf)
}

func newMailbox(name string, owner models.Account, mf MessageFetcher) Mailbox {
	return mailbox{mf, owner, name}
}

func (m mailbox) Name() string {
	return m.name
}

func (m mailbox) FetchAllMessages() ([]Message, error) {
	// TODO:(tauraamui) here we should store/cache mailboxes to
	//                  a prefix key set in the K/V DB
	return m.mf.fetchAllMessages(m)
}

func (m mailbox) FetchAllMessageUIDs() ([]MessageUID, error) {
	return m.mf.fetchAllMessageUIDs(m)
}
