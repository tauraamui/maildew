package repo_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/storage"
	"github.com/tauraamui/maildew/internal/storage/models"
	"github.com/tauraamui/maildew/internal/storage/repo"
)

func resolveEmailRepo() (repo.Emails, error) {
	db, err := storage.NewMemDB()
	if err != nil {
		return repo.Emails{}, err
	}

	return repo.Emails{DB: db}, nil
}

func TestSaveEmail(t *testing.T) {
	is := is.New(t)

	r, err := resolveEmailRepo()
	is.NoErr(err)
	defer r.Close()

	email := models.Email{
		Subject: "Fake email",
	}

	is.NoErr(r.Save(&email))

	is.NoErr(compareContentsWithExpected(r.DB, map[string][]byte{
		"emails":           {0, 0, 0, 0, 0, 0, 0, 100},
		"emails.subject.0": []byte("Fake email"),
	}))
}
