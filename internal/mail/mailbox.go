package mail

import (
	"github.com/tauraamui/maildew/internal/storage/models"
	"github.com/tauraamui/maildew/internal/storage/repo"
)

type Mailbox interface {
	Name() string
	FetchAllMessages() ([]Message, error)
	FetchAllMessageUIDs() ([]MessageUID, error)
}

type mailbox struct {
	mr      repo.Messages
	mf      messageFetcher
	account models.Account
	name    string
}

func newMailbox(mr repo.Messages, name string, owner models.Account, mf messageFetcher) Mailbox {
	return mailbox{mr, mf, owner, name}
}

func (m mailbox) Name() string {
	return m.name
}

func (m mailbox) FetchAllMessages() ([]Message, error) {
	existingMessageMUIDs, err := m.mf.fetchAllMessageUIDs(m)
	if err != nil {
		return nil, err
	}

	existingMessageUIDs := convertMessageUIDsToUint32(existingMessageMUIDs)

	storedMessages, err := m.mr.GetAll(m.account.ID)
	if err != nil {
		return nil, err
	}
	storedMessageUIDs := extractEmailUIDs(storedMessages)

	new, deleted := ResolveAddedAndRemoved(existingMessageUIDs, storedMessageUIDs)

	return nil, nil
}

func convertMessageUIDsToUint32(msgUIDs []MessageUID) []uint32 {
	ids := make([]uint32, len(msgUIDs))
	for _, e := range msgUIDs {
		ids = append(ids, uint32(e))
	}
	return ids
}

func extractEmailUIDs(msgs []models.Email) []uint32 {
	ids := make([]uint32, len(msgs))
	for _, e := range msgs {
		ids = append(ids, e.ID)
	}
	return ids
}

func (m mailbox) FetchAllMessageUIDs() ([]MessageUID, error) {
	return m.mf.fetchAllMessageUIDs(m)
}
