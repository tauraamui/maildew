package mail

import (
	"io"

	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/maildew/internal/kvs"
)

const (
	mailboxesTableName = "mailboxes"
)

type MailboxRepo interface {
	DumpTo(w io.Writer) error
	Save(owner kvs.UUID, mailbox Mailbox) error
	FetchByOwner(owner kvs.UUID) ([]Mailbox, error)
	Close()
}

func NewMailboxRepo(db kvs.DB) MailboxRepo {
	return mailboxRepo{DB: db}
}

type mailboxRepo struct {
	DB  kvs.DB
	seq *badger.Sequence
}

func (r mailboxRepo) DumpTo(w io.Writer) error {
	return r.DB.DumpTo(w)
}

func (r mailboxRepo) Save(owner kvs.UUID, mailbox Mailbox) error {
	rowID, err := r.nextRowID()
	if err != nil {
		return err
	}

	return saveValueWithUUID(r.DB, r.tableName(), owner, rowID, mailbox)
}

func (r mailboxRepo) FetchByOwner(owner kvs.UUID) ([]Mailbox, error) {
	return fetchByOwner[Mailbox](r.DB, r.tableName(), owner)
}

func saveValueWithUUID(db kvs.DB, tableName string, ownerID kvs.UUID, rowID uint32, v interface{}) error {
	if v == nil {
		return nil
	}
	entries := kvs.ConvertToEntriesWithUUID(tableName, ownerID, rowID, v)
	for _, e := range entries {
		if err := kvs.Store(db, e); err != nil {
			return err
		}
	}

	return nil
}

func (r mailboxRepo) tableName() string {
	return mailboxesTableName
}

func (r mailboxRepo) nextRowID() (uint32, error) {
	if r.seq == nil {
		seq, err := r.DB.GetSeq([]byte(r.tableName()), 1)
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

func (r mailboxRepo) Close() {
	if r.seq == nil {
		return
	}
	r.seq.Release()
}
