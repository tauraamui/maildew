package repo

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/maildew/internal/kvs"
	"github.com/tauraamui/maildew/internal/storage/models"
)

const (
	accountsTableName = "accounts"
)

type Accounts struct {
	DB  kvs.DB
	seq *badger.Sequence
}

func (r *Accounts) Save(user *models.Account) error {
	rowID, err := r.nextRowID()
	if err != nil {
		return err
	}

	return saveValue(r.DB, r.tableName(), rowID, 0, user)
}

func (r *Accounts) GetByID(rowID uint32) (models.Account, error) {
	acc := models.Account{
		ID: uint32(rowID),
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

func (r *Accounts) GetAll() ([]models.Account, error) {
	accounts := []models.Account{}

	blankEntries := kvs.ConvertToBlankEntries(r.tableName(), 0, 0, models.Account{})
	for _, ent := range blankEntries {
		// iterate over all stored values for this entry
		prefix := ent.PrefixKey()
		r.DB.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			var rows uint32 = 0
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				// be very clear with our append conditions
				if len(accounts) == 0 || rows >= uint32(len(accounts)) {
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
				if err := kvs.LoadEntry(&accounts[rows], ent); err != nil {
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

func (r *Accounts) nextRowID() (uint32, error) {
	if r.seq == nil {
		seq, err := r.DB.GetSeq([]byte(accountsTableName), 1)
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

func (r Accounts) Close() {
	if r.seq == nil {
		return
	}
	r.seq.Release()
}
