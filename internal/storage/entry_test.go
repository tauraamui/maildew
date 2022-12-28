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
		RowID:      0,
		Data:       []byte{},
	}

	db, err := storage.NewMemDB()
	is.NoErr(err)
	defer db.Close()

	seq, err := db.GetSeq(e.PrefixKey(), 100)
	is.NoErr(err) // error occurred on getting db sequence

	id, err := seq.Next()
	is.NoErr(err) // error occurred when aquiring next iter value

	e.RowID = id

	is.NoErr(storage.Store(db, e)) // error occurred when calling store
}
