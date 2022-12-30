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

func TestConvertToEntries(t *testing.T) {
	is := is.New(t)

	i := struct {
		Foo string
		Bar int
	}{
		Foo: "Foo",
		Bar: 4,
	}

	e := storage.ConvertToEntries("test", 0, i)
	is.Equal(len(e), 2)

	is = is.NewRelaxed(t)

	is.Equal(storage.Entry{
		TableName:  "test",
		ColumnName: "foo",
		Data:       []byte{70, 111, 111},
	}, e[0])

	is.Equal(storage.Entry{
		TableName:  "test",
		ColumnName: "bar",
		Data:       []byte{0, 0, 0, 0, 0, 0, 0, 4},
	}, e[1])
}

func TestSequences(t *testing.T) {
	is := is.New(t)

	db, err := storage.NewMemDB()
	is.NoErr(err)
	defer db.Close()

	fruitEntry := storage.Entry{
		TableName:  "fruits",
		ColumnName: "color",
	}

	chocolateEntry := storage.Entry{
		TableName:  "chocolate",
		ColumnName: "flavour",
	}

	fruitSeq, err := db.GetSeq(fruitEntry.PrefixKey(), 100)
	is.NoErr(err) // error occurred on getting db sequence
	defer fruitSeq.Release()

	chocolateSeq, err := db.GetSeq(chocolateEntry.PrefixKey(), 100)
	is.NoErr(err) // error occurred on getting db sequence
	defer fruitSeq.Release()

	id, err := fruitSeq.Next()
	is.NoErr(err) // error occurred when aquiring next iter value
	is.Equal(id, uint64(0))

	id, err = chocolateSeq.Next()
	is.NoErr(err) // error occurred when aquiring next iter value
	is.Equal(id, uint64(0))

	id, err = fruitSeq.Next()
	is.NoErr(err) // error occurred when aquiring next iter value
	is.Equal(id, uint64(1))

	id, err = chocolateSeq.Next()
	is.NoErr(err) // error occurred when aquiring next iter value
	is.Equal(id, uint64(1))

	id, err = fruitSeq.Next()
	is.NoErr(err) // error occurred when aquiring next iter value
	is.Equal(id, uint64(2))

	id, err = chocolateSeq.Next()
	is.NoErr(err) // error occurred when aquiring next iter value
	is.Equal(id, uint64(2))
}
