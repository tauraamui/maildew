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
	blankEntries := storage.ConvertToBlankEntries(accountsTableName, rowID, acc)
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
	// func IterateKeys(db *badger.DB, prefix []byte) error {
	// 	return db.View(func(txn *badger.Txn) error {
	// 		it := txn.NewIterator(badger.DefaultIteratorOptions)
	// 		defer it.Close()
	// 		prefixIter := it.Seek(prefix)

	// 		for prefixIter.ValidForPrefix(prefix) {
	// 			item := prefixIter.Item()
	// 			key := item.Key()
	// 			fmt.Println(key)
	// 			prefixIter.Next()
	// 		}

	// 		return nil
	// 	})
	// }

	return nil, nil
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
