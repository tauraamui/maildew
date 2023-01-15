package repo

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/maildew/internal/storage"
	"github.com/tauraamui/maildew/internal/storage/models"
)

const (
	accountsTableName = "accounts"
)

type Accounts struct {
	DB  storage.DB
	seq *badger.Sequence
}

func (r *Accounts) Save(user *models.Account) error {
	rowID, err := r.nextRowID()
	if err != nil {
		return err
	}

	entries := storage.ConvertToEntries(accountsTableName, 0, rowID, *user)
	for _, e := range entries {
		if err := storage.Store(r.DB, e); err != nil {
			return err
		}
	}

	user.ID = rowID

	return nil
}

func (r *Accounts) GetByID(rowID uint64) (models.Account, error) {
	acc := models.Account{
		ID: rowID,
	}
	blankEntries := storage.ConvertToBlankEntries(r.tableName(), 0, rowID, acc)
	for _, e := range blankEntries {
		if err := storage.Get(r.DB, &e); err != nil {
			return acc, err
		}

		if err := storage.LoadEntry(&acc, e); err != nil {
			return acc, err
		}
	}

	return acc, nil
}

func (r *Accounts) GetAll() ([]models.Account, error) {
	accounts := []models.Account{}

	blankEntries := storage.ConvertToBlankEntries(r.tableName(), 0, 0, models.Account{})
	for _, ent := range blankEntries {
		// iterate over all stored values for this entry
		prefix := ent.PrefixKey()
		r.DB.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			var rows uint64 = 0
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				// be very clear with our append conditions
				if len(accounts) == 0 || rows >= uint64(len(accounts)) {
					accounts = append(accounts, models.Account{
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
				if err := storage.LoadEntry(&accounts[rows], ent); err != nil {
					return err
				}
				rows++
			}

			return nil
		})
	}

	return accounts, nil
}

func (r *Accounts) tableName() string {
	return accountsTableName
}

func (r *Accounts) nextRowID() (uint64, error) {
	if r.seq == nil {
		seq, err := r.DB.GetSeq([]byte(accountsTableName), 1)
		if err != nil {
			return 0, err
		}
		r.seq = seq
	}

	return r.seq.Next()
}

func (r Accounts) Close() {
	if r.seq == nil {
		return
	}
	r.seq.Release()
}
