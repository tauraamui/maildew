package mail

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/kvs"
)

func TestFetchByOwner(t *testing.T) {
	is := is.New(t)

	db, err := kvs.NewMemDB()
	is.NoErr(err)

	accRepo := NewAccountRepo(db)
	mbRepo := NewMailboxRepo(db)

	mbcount := 10

	accWithMboxes := Account{UUID: uuid.New(), Username: "username1", Password: "password"}
	accRepo.Save(accWithMboxes)

	for i := 0; i < mbcount; i++ {
		mbRepo.Save(accWithMboxes.UUID, Mailbox{
			UUID: uuid.New(),
			Name: fmt.Sprintf("INBOX%d", i),
		})
	}

	accWithoutMboxes := Account{UUID: uuid.New(), Username: "username2", Password: "password"}
	accRepo.Save(accWithoutMboxes)

	fetchedMboxes, err := mbRepo.FetchByOwner(accWithMboxes.UUID)
	is.NoErr(err)
	is.Equal(len(fetchedMboxes), mbcount)

	fetchedMboxes, err = mbRepo.FetchByOwner(accWithoutMboxes.UUID)
	is.NoErr(err)
	is.Equal(len(fetchedMboxes), 0)
}
