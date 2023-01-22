package mail_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/emersion/go-imap/backend/memory"
	"github.com/emersion/go-imap/server"
	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/mail"
)

// NOTE: there are comments next to most of the assertions, especially
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

func TestClientListMailboxes(t *testing.T) {
	is := is.New(t)

	client, err, cleanup := setupClientConnection()
	defer cleanup()
	is.NoErr(err) // error connecting to imap server
	is.True(client != nil)

	mailboxes, err := client.Mailboxes()
	is.NoErr(err)
	is.Equal([]mail.Mailbox{{Name: "/"}}, mailboxes)
}

func setupClientConnection() (mail.Client, error, func() error) {
	l, err := setupListener()
	if err != nil {
		return nil, err, nil
	}

	err, shutdown := startLocalServer(l)
	if err != nil {
		return nil, err, func() error {
			// TODO: implement the usage of an error group to collect all errors
			//       from cleaning up connections etc., rather than ignoring them all
			l.Close()
			return nil
		}
	}

	addr := l.Addr().String()

	client, err := mail.Connect(addr, "username", "password")
	if err != nil {
		return nil, err, func() error {
			// TODO: implement the usage of an error group to collect all errors
			//       from cleaning up connections etc., rather than ignoring them all
			l.Close()
			shutdown()
			return nil
		}
	}

	return client, nil, func() error {
		// TODO: implement the usage of an error group to collect all errors
		//       from cleaning up connections etc., rather than ignoring them all
		l.Close()
		shutdown()
		return nil
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
	s := server.New(memory.New())
	s.AllowInsecureAuth = true

	go s.Serve(l)

	return nil, s.Close
}
