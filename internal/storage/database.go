package storage

import (
	"encoding/binary"

	"github.com/dgraph-io/badger/v3"
)

type DB struct {
	conn     *badger.DB
	incIndex int
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

func (db DB) persistIndex() {
	db.conn.Update(func(txn *badger.Txn) error {
		ib := make([]byte, 8)
		binary.LittleEndian.PutUint64(ib, uint64(db.incIndex))
		txn.Set([]byte("x_pk"), ib)
		return nil
	})
}
