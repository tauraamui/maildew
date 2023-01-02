package storage

import (
	"encoding/json"
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
		f := vv.Type().Field(i)

		fOpts := resolveFieldOptions(f)
		if fOpts.Ignore {
			continue
		}

		bd, err := convertToBytes(v.Field(i).Interface())
		if err != nil {
			return entries
		}

		e := Entry{
			TableName:  tableName,
			ColumnName: strings.ToLower(f.Name),
			RowID:      rowID,
			Data:       bd,
		}

		entries = append(entries, e)
	}

	return entries
}

func convertToBytes(i interface{}) ([]byte, error) {
	// Check the type of the interface.
	switch v := i.(type) {
	case []byte:
		// Return the input as a []byte if it is already a []byte.
		return v, nil
	case string:
		// Convert the string to a []byte and return it.
		return []byte(v), nil
	default:
		// Use json.Marshal to convert the interface to a []byte.
		return json.Marshal(v)
	}
}

func convertFromBytes(data []byte, i interface{}) error {
	// Check that the destination argument is a pointer.
	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}

	// Check the type of the interface.
	switch v := i.(type) {
	case *[]byte:
		// Set the value of the interface to the []byte if it is a pointer to a []byte.
		*v = data
		return nil
	case *string:
		// Convert the []byte to a string and set the value of the interface to the string.
		*v = string(data)
		return nil
	default:
		// Use json.Unmarshal to convert the []byte to the interface.
		return json.Unmarshal(data, v)
	}
}

type mdbFieldOptions struct {
	Ignore bool
}

func resolveFieldOptions(f reflect.StructField) mdbFieldOptions {
	mdbTagValue := f.Tag.Get("mdb")
	return mdbFieldOptions{
		Ignore: strings.Contains(mdbTagValue, "ignore"),
	}
}
