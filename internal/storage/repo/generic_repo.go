package repo

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/maildew/internal/storage"
)

type Value interface {
	SetID(id uint32)
	Ref() interface{}
}

type GenericRepo struct {
	TableName string
	DB        storage.DB
	seq       *badger.Sequence
}

func (r *GenericRepo) Save(ownerID uint32, v Value) error {
	rowID, err := r.nextRowID()
	if err != nil {
		return err
	}

	entries := storage.ConvertToEntries(r.TableName, ownerID, rowID, v)
	for _, e := range entries {
		if err := storage.Store(r.DB, e); err != nil {
			return err
		}
	}

	v.SetID(rowID)
	return nil
}

func (r *GenericRepo) nextRowID() (uint32, error) {
	if r.seq == nil {
		seq, err := r.DB.GetSeq([]byte(r.TableName), 100)
		if err != nil {
			return 0, err
		}
		r.seq = seq
	}

	s, err := r.seq.Next()
	if err != nil {
		return 0, err
	}

	return uint32(s), nil
}

func (r *GenericRepo) Close() {
	if r.seq == nil {
		return
	}
	r.seq.Release()
}
