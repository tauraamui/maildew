package repo_test

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/kvs"
	"github.com/tauraamui/maildew/internal/storage/models"
	"github.com/tauraamui/maildew/internal/storage/repo"
)

func resolveMailboxRepo() (repo.Mailboxes, error) {
	db, err := kvs.NewMemDB()
	if err != nil {
		return repo.Mailboxes{}, err
	}

	return repo.Mailboxes{DB: db}, nil
}

func TestSaveMailbox(t *testing.T) {
	t.Skip("pending migration to UUID")
	is := is.New(t)

	r, err := resolveMailboxRepo()
	is.NoErr(err)
	defer r.Close()

	mailbox := models.Mailbox{
		UID:  112,
		Name: "Fake mailbox",
	}

	is.NoErr(r.Save(0, &mailbox))
	is.NoErr(compareContentsWithExpected(r.DB, map[string][]byte{
		"mailboxes":          {0, 0, 0, 0, 0, 0, 0, 100},
		"mailboxes.uid.0.0":  helperConvertToBytes(t, 112),
		"mailboxes.name.0.0": helperConvertToBytes(t, "Fake mailbox"),
	}))
}

func TestGetMailbox(t *testing.T) {
	t.Skip("pending migration to UUID")
	is := is.New(t)

	r, err := resolveMailboxRepo()
	is.NoErr(err)
	defer r.Close()

	is.NoErr(insertContents(r.DB, map[string][]byte{
		"mailboxes":          {0, 0, 0, 0, 0, 0, 0, 1},
		"mailboxes.uid.0.0":  helperConvertToBytes(t, 684),
		"mailboxes.name.0.0": helperConvertToBytes(t, "Fake mailbox"),
	}))

	mb, err := r.GetByID(0)
	is.NoErr(err)

	is.Equal(mb.UID, uint32(684))
	is.Equal(mb.Name, "Fake mailbox")
}

func TestGetAllMailboxes(t *testing.T) {
	t.Skip("pending migration to UUID")
	is := is.New(t)

	r, err := resolveMailboxRepo()
	is.NoErr(err)
	defer r.Close()

	is.NoErr(insertContents(r.DB, map[string][]byte{
		"mailboxes":          {0, 0, 0, 0, 0, 0, 0, 1},
		"mailboxes.uid.0.0":  helperConvertToBytes(t, 244),
		"mailboxes.name.0.0": helperConvertToBytes(t, "First fake mailbox"),
		"mailboxes.uid.0.1":  helperConvertToBytes(t, 968),
		"mailboxes.name.0.1": helperConvertToBytes(t, "Second fake mailbox"),
	}))

	mbs, err := r.GetAll(0)
	is.NoErr(err)
	is.Equal(len(mbs), 2)

	first, second := mbs[0], mbs[1]
	is.Equal(first.UID, uint32(244))
	is.Equal(first.Name, "First fake mailbox")
	is.Equal(second.UID, uint32(968))
	is.Equal(second.Name, "Second fake mailbox")
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
