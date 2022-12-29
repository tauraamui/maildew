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

func (r *Accounts) Save(user models.Account) error {
	emailEntry := storage.Entry{
		TableName:  accountsTableName,
		ColumnName: "email",
		Data:       []byte(user.Email),
	}

	if r.seq == nil {
		seq, err := r.DB.GetSeq(emailEntry.PrefixKey(), 100)
		if err != nil {
			return err
		}
		r.seq = seq
	}

	rowID, err := r.seq.Next()
	if err != nil {
		return err
	}

	emailEntry.RowID = rowID

	nickEntry := storage.Entry{
		TableName:  accountsTableName,
		ColumnName: "nick",
		Data:       []byte(user.Nick),
		RowID:      rowID,
	}

	passwordEntry := storage.Entry{
		TableName:  accountsTableName,
		ColumnName: "password",
		Data:       []byte(user.Password),
		RowID:      rowID,
	}

	if err := storage.Store(r.DB, emailEntry); err != nil {
		return err
	}

	if err := storage.Store(r.DB, nickEntry); err != nil {
		return err
	}

	if err := storage.Store(r.DB, passwordEntry); err != nil {
		return err
	}

	return nil
}

func (r Accounts) Close() {
	if r.seq == nil {
		return
	}
	r.seq.Release()
}
