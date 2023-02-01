package mail

import "github.com/tauraamui/maildew/internal/storage/models"

type Mailbox struct {
	mf      messageFetcher
	account models.Account
	Name    string
}

func (m Mailbox) FetchAllMessages() ([]Message, error) {
	// TODO:(tauraamui) here we should store/cache mailboxes to
	//                  a prefix key set in the K/V DB
	return m.mf.fetchAllMessages(m)
}

func (m Mailbox) FetchAllMessageUIDs() ([]MessageUID, error) {
	return m.mf.fetchAllMessageUIDs(m)
}
