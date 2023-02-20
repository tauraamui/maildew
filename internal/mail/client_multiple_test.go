package mail_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/mail"
	"github.com/tauraamui/maildew/internal/mail/mock"
	"github.com/tauraamui/maildew/internal/storage"
	"github.com/tauraamui/maildew/internal/storage/models"
)

func TestClientConnectToMultipleAccounts(t *testing.T) {
	is := is.New(t)

	l, err := setupListener()
	is.NoErr(err) // error setting up the net listener
	defer l.Close()

	backend := mock.Backend{}
	backend.RegisterUser("fake1@place.com", "fakepass")
	backend.RegisterUser("fake2@place.com", "secondfakepass")
	backend.RegisterUser("fake3@place.com", "thirdfakepass")

	err, shutdown := startLocalServerWithBackend(l, &backend)
	is.NoErr(err)

	defer func() {
		is.NoErr(shutdown())
	}()

	addr := l.Addr().String()

	client := mail.NewClient(storage.DB{})
	is.True(client != nil)

	err = client.Connect(addr, models.Account{
		Email: "fake1@place.com", Password: "fakepass",
	})
	is.NoErr(err) // error connecting to imap server

	err = client.Connect(addr, models.Account{
		Email: "fake2@place.com", Password: "secondfakepass",
	})
	is.NoErr(err) // error connecting to imap server

	err = client.Connect(addr, models.Account{
		Email: "fake3@place.com", Password: "thirdfakepass",
	})
	is.NoErr(err) // error connecting to imap server
}
