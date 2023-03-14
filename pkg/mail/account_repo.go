package mail

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/maildew/internal/kvs"
)

const (
	accountsTableName = "accounts"
)

type AccountRepo interface {
	Save(user Account) error
	Close()
}

func NewAccountRepo(db kvs.DB) AccountRepo {
	return accountRepo{DB: db}
}

type accountRepo struct {
	DB  kvs.DB
	seq *badger.Sequence
}

func (r accountRepo) Save(user Account) error {
	rowID, err := r.nextRowID()
	if err != nil {
		return err
	}

	return saveValueWithUUID(r.DB, r.tableName(), kvs.RootOwner{}, rowID, user)
}

func (r accountRepo) tableName() string {
	return accountsTableName
}

func (r accountRepo) nextRowID() (uint32, error) {
	if r.seq == nil {
		seq, err := r.DB.GetSeq([]byte(accountsTableName), 1)
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

func (r accountRepo) Close() {
	if r.seq == nil {
		return
	}
	r.seq.Release()
}
