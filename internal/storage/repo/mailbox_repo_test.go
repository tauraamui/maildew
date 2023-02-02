package repo_test

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/storage"
	"github.com/tauraamui/maildew/internal/storage/models"
	"github.com/tauraamui/maildew/internal/storage/repo"
)

func resolveMailboxRepo() (repo.Mailboxes, error) {
	db, err := storage.NewMemDB()
	if err != nil {
		return repo.Mailboxes{}, err
	}

	return repo.Mailboxes{DB: db}, nil
}

func TestSaveMailboxSuccess(t *testing.T) {
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
		"mailboxes.name.0.0": []byte("Fake mailbox"),
	}))
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
