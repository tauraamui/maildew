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

func fetchByOwner[E Account | Mailbox | Message](db kvs.DB, tableName string, owner kvs.UUID) ([]E, error) {
	entries := make([]E, 1)

	blankEntries := kvs.ConvertToBlankEntriesWithUUID(tableName, owner, 0, entries[0])
	for _, ent := range blankEntries {
		// iterate over all stored values for this entry
		prefix := ent.PrefixKey()
		db.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			var rows uint32 = 0
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				if rows >= uint32(len(entries)) {
					entries = append(entries)
				}
				item := it.Item()
				ent.RowID = rows
				if err := item.Value(func(val []byte) error {
					ent.Data = val
					return nil
				}); err != nil {
					return err
				}

				if err := kvs.LoadEntry(&entries[rows], ent); err != nil {
					return err
				}
				rows++
			}
			return nil
		})
	}
	return entries, nil
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
