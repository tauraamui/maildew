package storage

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"

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

func ConvertToEntries(tableName string, rowID uint64, x interface{}) []Entry {
	v := reflect.ValueOf(x)

	entries := []Entry{}

	for i := 0; i < v.NumField(); i++ {
		vv := reflect.Indirect(v)
		e := Entry{
			TableName:  tableName,
			ColumnName: strings.ToLower(vv.Type().Field(i).Name),
			RowID:      rowID,
			Data:       ConvertToBytes(v.Field(i).Interface()),
		}

		entries = append(entries, e)
	}

	return entries
}

func ConvertToBytes(x interface{}) []byte {
	switch v := x.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	case int:
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(v))
		return b
	case bool:
		if v {
			return []byte{0x1}
		}
		return []byte{0x0}
	}
	return nil
}
