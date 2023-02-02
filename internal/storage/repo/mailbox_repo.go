package repo

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/maildew/internal/storage"
	"github.com/tauraamui/maildew/internal/storage/models"
)

const (
	mailboxesTableName = "mailboxes"
)

// TODO:(tauraamui): this mailboxes repo should also be in charge of saving and loading
//
//	messages, when saving a message it would just simply set the owner
//	ID to be the mailboxes' ID
type Mailboxes struct {
	DB  storage.DB
	seq *badger.Sequence
}

func (r *Mailboxes) Save(accountID uint32, mailbox *models.Mailbox) error {
	rowID, err := r.nextRowID()
	if err != nil {
		return err
	}

	entries := storage.ConvertToEntries(mailboxesTableName, accountID, rowID, *mailbox)
	for _, e := range entries {
		if err := storage.Store(r.DB, e); err != nil {
			return err
		}
	}

	mailbox.ID = rowID

	return nil
}

func (r *Mailboxes) GetByID(rowID uint32) (models.Mailbox, error) {
	mb := models.Mailbox{
		ID: rowID,
	}
	blankEntries := storage.ConvertToBlankEntries(r.tableName(), 0, rowID, mb)
	for _, e := range blankEntries {
		if err := storage.Get(r.DB, &e); err != nil {
			return mb, err
		}

		if err := storage.LoadEntry(&mb, e); err != nil {
			return mb, err
		}
	}
	return mb, nil
}

func (r *Mailboxes) GetAll(accountID uint32) ([]models.Mailbox, error) {
	return []models.Mailbox{}, nil
}

func (r *Mailboxes) tableName() string {
	return mailboxesTableName
}

func (r *Mailboxes) nextRowID() (uint32, error) {
	if r.seq == nil {
		seq, err := r.DB.GetSeq([]byte(r.tableName()), 100)
		if err != nil {
			return 0, err
		}
		r.seq = seq
	}

	s, err := r.seq.Next()
	if err != nil {
		return 0, err
	}

	return uint32(s), nil
}

func (r *Mailboxes) Close() {
	if r.seq == nil {
		return
	}
	r.seq.Release()
}
