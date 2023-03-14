package mail

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/maildew/internal/kvs"
)

const (
	mailboxesTableName = "mailboxes"
)

type MailboxRepo interface {
	Save(owner kvs.UUID, mailbox Mailbox) error
	Close()
}

func NewMailboxRepo(db kvs.DB) MailboxRepo {
	return mailboxRepo{DB: db}
}

type mailboxRepo struct {
	DB  kvs.DB
	seq *badger.Sequence
}

func (r mailboxRepo) Save(owner kvs.UUID, mailbox Mailbox) error {
	rowID, err := r.nextRowID()
	if err != nil {
		return err
	}

	return saveValueWithUUID(r.DB, r.tableName(), owner, rowID, mailbox)
}

func saveValueWithUUID(db kvs.DB, tableName string, ownerID kvs.UUID, rowID uint32, v interface{}) error {
	entries := kvs.ConvertToEntriesWithUUID(tableName, ownerID, rowID, v)
	for _, e := range entries {
		if err := kvs.Store(db, e); err != nil {
			return err
		}
	}

	return nil
}

func (r mailboxRepo) tableName() string {
	return accountsTableName
}

func (r mailboxRepo) nextRowID() (uint32, error) {
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

func (r mailboxRepo) Close() {
	if r.seq == nil {
		return
	}
	r.seq.Release()
}
