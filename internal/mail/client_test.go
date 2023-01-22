package mail_test

import (
	"testing"

	"github.com/emersion/go-imap/backend/memory"
	"github.com/emersion/go-imap/server"
	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/mail"
)

func TestClientListingMailboxes(t *testing.T) {
	is := is.New(t)

	errc := make(chan error, 1)
	shutdown := startLocalServer(errc)
	defer shutdown()

	// FIX: connection refused not sure why
	client, err := mail.Connect("username", "password")
	is.NoErr(err)

	mailboxes, err := client.Mailboxes()
	is.NoErr(err)
	is.Equal(mailboxes, nil)
}

func startLocalServer(errc chan error) func() error {
	s := server.New(memory.New())

	s.Addr = ":1143"
	s.AllowInsecureAuth = true

	go func(errc chan error) {
		errc <- s.ListenAndServe()
	}(errc)

	return s.Close
}
