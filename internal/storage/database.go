package storage

import (
	"github.com/dgraph-io/badger/v3"
)

type DB struct {
	conn *badger.DB
}

func NewDB(db *badger.DB) (DB, error) {
	return newDB(false)
}

func NewMemDB() (DB, error) {
	return newDB(true)
}

func newDB(inMemory bool) (DB, error) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(inMemory))
	if err != nil {
		return DB{}, err
	}

	return DB{conn: db}, nil
}

func (db DB) GetSeq(key []byte, bandwidth uint64) (*badger.Sequence, error) {
	return db.conn.GetSequence(key, bandwidth)
}

func (db DB) Close() error {
	return db.conn.Close()
}
