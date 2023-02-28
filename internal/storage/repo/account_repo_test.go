package repo_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/kvs"
	"github.com/tauraamui/maildew/internal/storage/models"
	"github.com/tauraamui/maildew/internal/storage/repo"
)

func resolveRepo() (repo.Accounts, error) {
	db, err := kvs.NewMemDB()
	if err != nil {
		return repo.Accounts{}, err
	}

	return repo.Accounts{DB: db}, nil
}

func TestSaveUser(t *testing.T) {
	is := is.New(t)

	r, err := resolveRepo()
	is.NoErr(err)
	defer r.Close()

	user := models.Account{
		Email:    "test@place.com",
		Nick:     "Test User",
		Password: "fefweiofeifwwef",
	}

	is.NoErr(r.Save(&user))

	is.NoErr(compareContentsWithExpected(r.DB, map[string][]byte{
		"accounts":              {0, 0, 0, 0, 0, 0, 0, 1},
		"accounts.email.0.0":    []byte("test@place.com"),
		"accounts.nick.0.0":     []byte("Test User"),
		"accounts.password.0.0": []byte("fefweiofeifwwef"),
	}))
}

func TestGetUser(t *testing.T) {
	is := is.New(t)

	r, err := resolveRepo()
	is.NoErr(err)
	defer r.Close()

	is.NoErr(insertContents(r.DB, map[string][]byte{
		"accounts":              {0, 0, 0, 0, 0, 0, 0, 1},
		"accounts.email.0.0":    []byte("test@place.com"),
		"accounts.nick.0.0":     []byte("Test User"),
		"accounts.password.0.0": []byte("fefweiofeifwwef"),
	}))

	acc, err := r.GetByID(0)
	is.NoErr(err)

	is.Equal(acc.Email, "test@place.com")
	is.Equal(acc.Nick, "Test User")
	is.Equal(acc.Password, "fefweiofeifwwef")
}

func TestGetAllUsers(t *testing.T) {
	is := is.New(t)

	r, err := resolveRepo()
	is.NoErr(err)
	defer r.Close()

	is.NoErr(insertContents(r.DB, map[string][]byte{
		"accounts":              {0, 0, 0, 0, 0, 0, 0, 1},
		"accounts.email.0.0":    []byte("first@place.com"),
		"accounts.nick.0.0":     []byte("First User"),
		"accounts.password.0.0": []byte("wwqdwdqdqwdqd"),
		"accounts.email.0.1":    []byte("second@place.com"),
		"accounts.nick.0.1":     []byte("Second User"),
		"accounts.password.0.1": []byte("gigioregioigr"),
	}))

	accs, err := r.GetAll()
	is.NoErr(err)
	is.Equal(len(accs), 2)

	first, second := accs[0], accs[1]
	is.Equal(first.Nick, "First User")
	is.Equal(second.Nick, "Second User")
	is.Equal(first.Email, "first@place.com")
	is.Equal(second.Email, "second@place.com")
	is.Equal(first.Password, "wwqdwdqdqwdqd")
	is.Equal(second.Password, "gigioregioigr")
}

func insertContents(db kvs.DB, cnts map[string][]byte) error {
	return db.Update(func(txn *badger.Txn) error {
		for k, v := range cnts {
			if err := txn.Set([]byte(k), v); err != nil {
				return err
			}
		}
		return nil
	})
}

func compareContentsWithExpected(db kvs.DB, exp map[string][]byte) error {
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
