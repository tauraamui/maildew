package storage

import (
	"testing"

	"github.com/matryer/is"
)

func TestPersitIndex(t *testing.T) {
	is := is.New(t)

	db, err := NewMemDB()

	is.NoErr(err)

	db.persistIndex()
}
