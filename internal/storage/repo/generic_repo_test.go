package repo_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/kvs"
	"github.com/tauraamui/maildew/internal/storage/models"
	"github.com/tauraamui/maildew/internal/storage/repo"
)

func resolveGenericRepo() (repo.GenericRepo, error) {
	db, err := kvs.NewMemDB()
	if err != nil {
		return repo.GenericRepo{}, err
	}

	return repo.GenericRepo{TableName: "mailboxes", DB: db}, nil
}

func TestSaveGeneric(t *testing.T) {
	is := is.New(t)

	r, err := resolveGenericRepo()
	is.NoErr(err)
	defer r.Close()

	mailbox := models.Mailbox{
		UID:  83,
		Name: "Fake",
	}

	is.NoErr(r.Save(0, &mailbox))
	is.NoErr(compareContentsWithExpected(r.DB, map[string][]byte{
		"mailboxes":          {0, 0, 0, 0, 0, 0, 0, 100},
		"mailboxes.uid.0.0":  helperConvertToBytes(t, 83),
		"mailboxes.name.0.0": helperConvertToBytes(t, "Fake"),
	}))
}
