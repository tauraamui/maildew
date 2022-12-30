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
