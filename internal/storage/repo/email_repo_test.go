package repo_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/kvs"
	"github.com/tauraamui/maildew/internal/storage/models"
	"github.com/tauraamui/maildew/internal/storage/repo"
)

func resolveEmailRepo() (repo.Emails, error) {
	db, err := kvs.NewMemDB()
	if err != nil {
		return repo.Emails{}, err
	}

	return repo.Emails{DB: db}, nil
}

func TestSaveEmail(t *testing.T) {
	t.Skip("pending migration to UUID")
	is := is.New(t)

	r, err := resolveEmailRepo()
	is.NoErr(err)
	defer r.Close()

	email := models.Email{
		Subject: "Fake email",
	}

	is.NoErr(r.Save(0, &email))

	is.NoErr(compareContentsWithExpected(r.DB, map[string][]byte{
		"emails":             {0, 0, 0, 0, 0, 0, 0, 100},
		"emails.subject.0.0": []byte("Fake email"),
	}))
}
