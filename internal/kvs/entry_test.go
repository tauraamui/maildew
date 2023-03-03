package kvs_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/kvs"
)

func TestEntryStoreValuesInTable(t *testing.T) {
	is := is.New(t)

	e := kvs.Entry{
		TableName:  "users",
		ColumnName: "email",
		OwnerID:    11,
		Data:       []byte{0x33},
	}

	db, err := kvs.NewMemDB()
	is.NoErr(err)
	defer db.Close()

	seq, err := db.GetSeq(e.PrefixKey(), 100)
	is.NoErr(err) // error occurred on getting db sequence
	defer seq.Release()

	id, err := seq.Next()
	is.NoErr(err) // error occurred when aquiring next iter value

	e.RowID = uint32(id)

	is.NoErr(kvs.Store(db, e)) // error occurred when calling store

	newEntry := kvs.Entry{
		TableName:  e.TableName,
		ColumnName: e.ColumnName,
		OwnerID:    11,
		RowID:      e.RowID,
		Data:       nil,
	}
	is.NoErr(kvs.Get(db, &newEntry))

	is.Equal(newEntry.Data, []byte{0x33})
}

func TestGettingEntryOutOfTableErrorIncorrectKey(t *testing.T) {
	is := is.New(t)

	e := kvs.Entry{
		TableName:  "user",
		ColumnName: "emailz",
		OwnerID:    33,
		Data:       []byte{0x33},
	}

	db, err := kvs.NewMemDB()
	is.NoErr(err)
	defer db.Close()

	seq, err := db.GetSeq(e.PrefixKey(), 100)
	is.NoErr(err) // error occurred on getting db sequence
	defer seq.Release()

	id, err := seq.Next()
	is.NoErr(err) // error occurred when aquiring next iter value

	e.RowID = uint32(id)

	is.NoErr(kvs.Store(db, e)) // error occurred when calling store

	newEntry := kvs.Entry{
		TableName:  e.TableName,
		ColumnName: e.ColumnName,
		OwnerID:    11,
		RowID:      e.RowID,
		Data:       nil,
	}
	is.Equal(kvs.Get(db, &newEntry).Error(), "key not found: user.emailz.11.0")
	is.Equal(newEntry.Data, nil)
}

func TestConvertToEntries(t *testing.T) {
	is := is.New(t)

	source := struct {
		Foo string
		Bar int
	}{
		Foo: "Foo",
		Bar: 4,
	}

	e := kvs.ConvertToEntries("test", 0, 0, source)
	is.Equal(len(e), 2)

	is = is.NewRelaxed(t)

	is.Equal(kvs.Entry{
		TableName:  "test",
		ColumnName: "foo",
		Data:       []byte{70, 111, 111},
	}, e[0])

	is.Equal(kvs.Entry{
		TableName:  "test",
		ColumnName: "bar",
		Data:       []byte{52},
	}, e[1])
}

func TestUpdateStruct(t *testing.T) {
	// Define a struct type to use for the test
	type TestStruct struct {
		Field1 string
		Field2 int
		Field3 bool
	}

	// Create a slice of Entry values to use as input
	entries := []kvs.Entry{
		{ColumnName: "field1", Data: []byte("hello")},
		{ColumnName: "field2", Data: []byte("123")},
		{ColumnName: "field3", Data: []byte("true")},
	}

	s := TestStruct{}

	is := is.New(t)

	is.NoErr(kvs.LoadEntries(&s, entries)) // LoadEntries returned an error
	// Check that the values of the TestStruct fields were updated correctly
	expected := TestStruct{Field1: "hello", Field2: 123, Field3: true}
	is.Equal(s, expected) // Use the Equal method of the is package to compare the values
}

func TestSequences(t *testing.T) {
	is := is.New(t)

	db, err := kvs.NewMemDB()
	is.NoErr(err)
	defer db.Close()

	fruitEntry := kvs.Entry{
		TableName:  "fruits",
		ColumnName: "color",
	}

	chocolateEntry := kvs.Entry{
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
