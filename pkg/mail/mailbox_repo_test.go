package mail_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/kvs"
	"github.com/tauraamui/maildew/pkg/mail"
)

func TestSaveMailbox(t *testing.T) {
	is := is.New(t)

	r, db, err := resolveMailboxRepoWithDB()
	is.NoErr(err)
	is.True(r != nil)
	defer r.Close()

	uuidID := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	ownerUUID := uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479")

	mailbox := mail.Mailbox{
		UUID: uuidID,
		Name: "Fake mailbox",
	}

	is.NoErr(r.Save(ownerUUID, mailbox))
	is.True(db != nil)
	is.NoErr(compareContentsWithExpected(db, map[string][]byte{
		"mailboxes": {0, 0, 0, 0, 0, 0, 0, 1},
		"mailboxes.uuid.f47ac10b-58cc-0372-8567-0e02b2c3d479.0": helperConvertToBytes(t, fmt.Sprintf("\"%s\"", uuidID.String())),
		"mailboxes.name.f47ac10b-58cc-0372-8567-0e02b2c3d479.0": helperConvertToBytes(t, "Fake mailbox"),
	}))

}

func TestFetchingMultipleMailboxes(t *testing.T) {
	is := is.New(t)

	r, err := resolveMailboxRepo()
	is.NoErr(err)
	is.True(r != nil)
	defer r.Close()

	ownerUUID := uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	toSave := [10]mail.Mailbox{}
	for i := 0; i < 10; i++ {
		toSave[i] = mail.Mailbox{
			UUID: uuid.New(),
			Name: fmt.Sprintf("INBOX%d", i),
		}
		r.Save(ownerUUID, toSave[i])
	}

	fetched, err := r.FetchByOwner(ownerUUID)
	is.NoErr(err)
	is.Equal(len(fetched), 10)

	is.Equal(fetched[0].Name, "INBOX0")
	is.Equal(fetched[1].Name, "INBOX1")
	is.Equal(fetched[2].Name, "INBOX2")
	is.Equal(fetched[3].Name, "INBOX3")
	is.Equal(fetched[4].Name, "INBOX4")
	is.Equal(fetched[5].Name, "INBOX5")
	is.Equal(fetched[6].Name, "INBOX6")
	is.Equal(fetched[7].Name, "INBOX7")
	is.Equal(fetched[8].Name, "INBOX8")
	is.Equal(fetched[9].Name, "INBOX9")
}

func resolveMailboxRepo() (mail.MailboxRepo, error) {
	mb, _, err := resolveMailboxRepoWithDB()
	return mb, err
}

func resolveMailboxRepoWithDB() (mail.MailboxRepo, *kvs.DB, error) {
	db, err := kvs.NewMemDB()
	if err != nil {
		return nil, &db, err
	}

	mr := mail.NewMailboxRepo(db)
	return mr, &db, nil
}

func compareContentsWithExpected(db *kvs.DB, exp map[string][]byte) error {
	return db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()

			ev, ok := exp[string(k)]
			if !ok {
				return fmt.Errorf("unexpected stored key: %s", k)
			}

			if err := item.Value(func(v []byte) error {
				if !bytes.Equal(ev, v) {
					return fmt.Errorf("expected does not match stored: %v != %v", ev, v)
				}
				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

func helperConvertToBytes(t *testing.T, i interface{}) []byte {
	t.Helper()
	b, err := convertToBytes(i)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func convertToBytes(i interface{}) ([]byte, error) {
	// Check the type of the interface.
	switch v := i.(type) {
	case []byte:
		// Return the input as a []byte if it is already a []byte.
		return v, nil
	case string:
		// Convert the string to a []byte and return it.
		return []byte(v), nil
	default:
		// Use json.Marshal to convert the interface to a []byte.
		return json.Marshal(v)
	}
}

type fakeUUID struct{ f string }

func (f fakeUUID) String() string {
	return f.f
}
