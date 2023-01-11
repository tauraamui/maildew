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

	entries := storage.ConvertToEntries(accountsTableName, rowID, *user)
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
	blankEntries := storage.ConvertToBlankEntries(r.tableName(), rowID, acc)
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
	// acquire all entries for the type "Account"
	// this basically means for every single struct field,
	// return a list of entry types
	blankEntries := storage.ConvertToBlankEntries(r.tableName(), 0, models.Account{})

	// For each entry type (each struct field on "Account") search using the
	// prefix key for every stored value for that entry type. Each time we
	// find a value that means that there is a full "Account" that we need to
	// append to the accounts list.
	for _, ent := range blankEntries {
		prefix := ent.PrefixKey()
		r.DB.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				accounts = append(accounts, models.Account{})
				// item := it.Item()
				// key := item.Key()
				// err := item.Value()
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
		seq, err := r.DB.GetSeq([]byte(accountsTableName), 100)
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
