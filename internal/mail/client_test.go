package mail_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/emersion/go-imap/server"
	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/mail"
	"github.com/tauraamui/maildew/internal/mail/mock"
	"github.com/tauraamui/xerror/errgroup"
)

// NOTE:(tauraamui) there are comments next to most of the assertions, especially
//		 the no error ones because the matryer/is library shows the comments
//       in the test output when those assertions fail

func TestClientConnectToLocalMockServer(t *testing.T) {
	is := is.New(t)

	l, err := setupListener()
	is.NoErr(err) // error setting up the net listener
	defer l.Close()

	err, shutdown := startLocalServer(l)
	is.NoErr(err)

	defer func() {
		is.NoErr(shutdown())
	}()

	addr := l.Addr().String()

	client, err := mail.Connect(addr, "username", "password")
	is.NoErr(err) // error connecting to imap server
	is.True(client != nil)
}

func TestClientFetchMailboxes(t *testing.T) {
	is := is.New(t)

	client, err, cleanup := setupClientConnection()
	defer cleanup()
	is.NoErr(err)          // error connecting to imap server
	is.True(client != nil) // ensure client is not nil

	mailboxes, err := client.FetchAllMailboxes()
	is.NoErr(err) // error fetching mailboxes
	is.Equal([]mail.Mailbox{{Name: "INBOX"}}, mailboxes)
}

func TestClientFetchAllInboxMessages(t *testing.T) {
	is := is.New(t)

	client, err, cleanup := setupClientConnection()
	defer cleanup()
	is.NoErr(err)          // error connecting to imap server
	is.True(client != nil) // ensure client is not nil

	msgs, err := client.FetchAllMessages(mail.Mailbox{Name: "INBOX"})
	is.NoErr(err)                                                               // error fetching inbox messages
	is.Equal(msgs, []mail.Message{{Subject: "A little message, just for you"}}) // list of messages does not match expected
}

func TestClientFetchAllInboxMessageUIDs(t *testing.T) {
	is := is.New(t)

	client, err, cleanup := setupClientConnection()
	defer cleanup()
	is.NoErr(err)
	is.True(client != nil)

	uids, err := client.FetchAllMessageUIDs(mail.Mailbox{Name: "INBOX"})
	is.NoErr(err)
	is.Equal(uids, []mail.MessageUID{6})
}

func setupClientConnection() (mail.Client, error, func() error) {
	l, err := setupListener()
	if err != nil {
		return nil, err, nil
	}

	err, shutdown := startLocalServer(l)
	if err != nil {
		return nil, err, func() error {
			// NOTE:(tauraamui) since starting the server was the cause of the error
			//       it's unnecessary to call the given shutdown callback
			return l.Close()
		}
	}

	addr := l.Addr().String()

	client, err := mail.Connect(addr, "username", "password")
	if err != nil {
		return nil, err, func() error {
			errs := errgroup.I{}

			errs.Append(l.Close())
			errs.Append(shutdown())
			return errs.ToErrOrNil()
		}
	}

	return client, nil, func() error {
		errs := errgroup.I{}

		errs.Append(l.Close())
		errs.Append(shutdown())
		return errs.ToErrOrNil()
	}
}

func setupListener() (net.Listener, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("cannot listen: %w", err)
	}
	return l, nil
}

func startLocalServer(l net.Listener) (error, func() error) {
	s := server.New(mock.New())
	s.AllowInsecureAuth = true

	go s.Serve(l)

	return nil, s.Close
}
