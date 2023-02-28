package repo

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/maildew/internal/kvs"
	"github.com/tauraamui/maildew/internal/storage/models"
)

const (
	emailsTableName = "emails"
)

type Emails struct {
	DB  kvs.DB
	seq *badger.Sequence
}

func (r *Emails) Save(accountID uint32, email *models.Email) error {
	ownerID := accountID
	rowID, err := r.nextRowID()
	if err != nil {
		return err
	}

	return saveValue(r.DB, r.tableName(), rowID, ownerID, email)
}

func (r *Emails) GetByID(rowID uint32) (models.Email, error) {
	acc := models.Email{
		ID: rowID,
	}
	blankEntries := kvs.ConvertToBlankEntries(r.tableName(), 0, rowID, acc)
	for _, e := range blankEntries {
		if err := kvs.Get(r.DB, &e); err != nil {
			return acc, err
		}

		if err := kvs.LoadEntry(&acc, e); err != nil {
			return acc, err
		}
	}

	return acc, nil
}

// TODO:(tauraamui) really all of the generic "Getters" and "Setters" methods
// of these repos are identical thanks to using the kvs backend
// so should move each of these into using Go generics rather than
// copying them for each type.
func (r *Emails) GetAll(accountID uint32) ([]models.Email, error) {
	emails := make([]models.Email, 1)

	blankEntries := kvs.ConvertToBlankEntries(r.tableName(), 0, 0, emails[0])
	for _, ent := range blankEntries {
		// iterate over all stored values for this entry
		prefix := ent.PrefixKey()
		r.DB.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			var rows uint32 = 0
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				if rows >= uint32(len(emails)) {
					emails = append(emails, models.Email{
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
				if err := kvs.LoadEntry(&emails[rows], ent); err != nil {
					return err
				}
				rows++
			}

			return nil
		})
	}

	return emails, nil
}

func (r *Emails) tableName() string {
	return emailsTableName
}

func (r *Emails) nextRowID() (uint32, error) {
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

func (r *Emails) Close() {
	if r.seq == nil {
		return
	}
	r.seq.Release()
}
