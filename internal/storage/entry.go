package storage

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

type Entry struct {
	TableName  string
	ColumnName string
	RowID      uint64
	Data       []byte
}

func (e Entry) PrefixKey() []byte {
	return []byte(fmt.Sprintf("%s.%s", e.TableName, e.ColumnName))
}

func (e Entry) Key() []byte {
	return []byte(fmt.Sprintf("%s.%s.%d", e.TableName, e.ColumnName, e.RowID))
}

func Store(db DB, e Entry) error {
	return db.conn.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(e.Key()), e.Data)
	})
}

func Get(db DB, e *Entry) error {
	return db.conn.View(func(txn *badger.Txn) error {
		item, err := txn.Get(e.Key())
		if err != nil {
			return err
		}

		if err := item.Value(func(val []byte) error {
			e.Data = val
			return nil
		}); err != nil {
			return err
		}

		return nil
	})
}
