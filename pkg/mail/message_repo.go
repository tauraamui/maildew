package mail

import (
	"io"

	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/maildew/internal/kvs"
)

const (
	messagesTableName = "messages"
)

type MessageRepo interface {
	DumpTo(w io.Writer) error
	Save(owner kvs.UUID, msg Message) error
	FetchByOwner(owner kvs.UUID) ([]Message, error)
	Close() error
}

func NewMessageRepo(db kvs.DB) MessageRepo {
	return messageRepo{DB: db}
}

type messageRepo struct {
	DB  kvs.DB
	seq *badger.Sequence
}

func (r messageRepo) DumpTo(w io.Writer) error {
	return r.DB.DumpTo(w)
}

func (r messageRepo) Save(owner kvs.UUID, msg Message) error {
	rowID, err := r.nextRowID()
	if err != nil {
		return err
	}

	return saveValueWithUUID(r.DB, r.tableName(), owner, rowID, msg)
}

func (r messageRepo) FetchByOwner(owner kvs.UUID) ([]Message, error) {
	return fetchByOwner[Message](r.DB, r.tableName(), owner)
}

func (r messageRepo) tableName() string {
	return messagesTableName
}

func (r messageRepo) nextRowID() (uint32, error) {
	if r.seq == nil {
		seq, err := r.DB.GetSeq([]byte(messagesTableName), 1)
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

func (r messageRepo) Close() error {
	if r.seq == nil {
		return nil
	}
	r.seq.Release()
	return nil
}
