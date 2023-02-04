package storage

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Entry struct {
	TableName  string
	ColumnName string
	OwnerID    uint32
	RowID      uint32
	Data       []byte
}

func (e Entry) PrefixKey() []byte {
	return []byte(fmt.Sprintf("%s.%s.%d", e.TableName, e.ColumnName, e.OwnerID))
}

func (e Entry) Key() []byte {
	return []byte(fmt.Sprintf("%s.%s.%d.%d", e.TableName, e.ColumnName, e.OwnerID, e.RowID))
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

func ConvertToBlankEntries(tableName string, ownerID, rowID uint32, x interface{}) []Entry {
	v := reflect.ValueOf(x)
	return convertToEntries(tableName, ownerID, rowID, v, false)
}

func ConvertToEntries(tableName string, ownerID, rowID uint32, x interface{}) []Entry {
	v := reflect.ValueOf(x)
	return convertToEntries(tableName, ownerID, rowID, v, true)
}

func LoadEntry(s interface{}, entry Entry) error {
	// Convert the interface value to a reflect.Value so we can access its fields
	val := reflect.ValueOf(s).Elem()

	// Check if the entry's ColumnName field matches the name of any of the struct's fields
	field := val.FieldByName(cases.Title(language.English).String(entry.ColumnName))
	if !field.IsValid() {
		// check if entry column name matches full upper casing instead
		field = val.FieldByName(cases.Upper(language.English).String(entry.ColumnName))
		if !field.IsValid() {
			// The struct does not have a field with the same name as the entry's ColumnName, so return an error
			return fmt.Errorf("struct does not have a field with name %q", entry.ColumnName)
		}
	}

	// Convert the entry's Data field to the type of the target field
	if err := convertFromBytes(entry.Data, field.Addr().Interface()); err != nil {
		return fmt.Errorf("failed to convert entry data to field type: %v", err)
	}

	return nil
}

func LoadEntries(s interface{}, entries []Entry) error {
	for _, entry := range entries {
		if err := LoadEntry(s, entry); err != nil {
			return err
		}
	}

	return nil
}

// TODO:(tauraamui): come up with novel and well designed method of preserving full upper casing
func convertToEntries(tableName string, ownerID, rowID uint32, v reflect.Value, includeData bool) []Entry {
	entries := []Entry{}

	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		vv := reflect.Indirect(v)
		f := vv.Type().Field(i)

		fOpts := resolveFieldOptions(f)
		if fOpts.Ignore {
			continue
		}

		e := Entry{
			TableName:  tableName,
			ColumnName: strings.ToLower(f.Name),
			OwnerID:    ownerID,
			RowID:      rowID,
		}

		if includeData {
			bd, err := convertToBytes(v.Field(i).Interface())
			if err != nil {
				return entries
			}
			e.Data = bd
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
