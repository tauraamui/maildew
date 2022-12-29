package storage_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/storage"
)

func TestEntryStoreValuesInTable(t *testing.T) {
	is := is.New(t)

	e := storage.Entry{
		TableName:  "users",
		ColumnName: "email",
		Data:       []byte{0x33},
	}

	db, err := storage.NewMemDB()
	is.NoErr(err)
	defer db.Close()

	seq, err := db.GetSeq(e.PrefixKey(), 100)
	is.NoErr(err) // error occurred on getting db sequence
	defer seq.Release()

	id, err := seq.Next()
	is.NoErr(err) // error occurred when aquiring next iter value

	e.RowID = id

	is.NoErr(storage.Store(db, e)) // error occurred when calling store

	newEntry := storage.Entry{
		TableName:  e.TableName,
		ColumnName: e.ColumnName,
		RowID:      e.RowID,
		Data:       nil,
	}
	is.NoErr(storage.Get(db, &newEntry))

	is.Equal(newEntry.Data, []byte{0x33})
}
