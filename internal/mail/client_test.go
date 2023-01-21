package mail_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/mail"
)

func TestClientListingMailboxes(t *testing.T) {
	is := is.New(t)

	client := mail.Connect("username", "password")

	is.Equal(client.Mailboxes(), nil)
}
