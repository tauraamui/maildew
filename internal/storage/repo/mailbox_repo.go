package repo

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/maildew/internal/kvs"
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
	DB  kvs.DB
	seq *badger.Sequence
}

func (r *Mailboxes) Save(accountID uint32, mailbox *models.Mailbox) error {
	ownerID := accountID
	rowID, err := r.nextRowID()
	if err != nil {
		return err
	}

	return saveValue(r.DB, r.tableName(), rowID, ownerID, mailbox)
}

func (r *Mailboxes) GetByID(rowID uint32) (models.Mailbox, error) {
	mb := models.Mailbox{
		ID: rowID,
	}
	blankEntries := kvs.ConvertToBlankEntries(r.tableName(), 0, rowID, mb)
	for _, e := range blankEntries {
		if err := kvs.Get(r.DB, &e); err != nil {
			return mb, err
		}

		if err := kvs.LoadEntry(&mb, e); err != nil {
			return mb, err
		}
	}
	return mb, nil
}

func (r *Mailboxes) GetAll(accountID uint32) ([]models.Mailbox, error) {
	mailboxes := make([]models.Mailbox, 1)

	blankEntries := kvs.ConvertToBlankEntries(r.tableName(), 0, 0, mailboxes[0])
	for _, ent := range blankEntries {
		// iterate over all stored values for this entry
		prefix := ent.PrefixKey()
		r.DB.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			var rows uint32 = 0
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				if rows >= uint32(len(mailboxes)) {
					mailboxes = append(mailboxes, models.Mailbox{
						ID: rows,
					})
				}
				item := it.Item()
				ent.RowID = rows
				if err := item.Value(func(val []byte) error {
					ent.Data = val
					return nil
				}); err != nil {
					return err
				}
				if err := kvs.LoadEntry(&mailboxes[rows], ent); err != nil {
					return err
				}
				rows++
			}

			return nil
		})
	}

	return mailboxes, nil
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
