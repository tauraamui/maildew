package mail_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-imap/server"
	"github.com/matryer/is"
	"github.com/tauraamui/maildew/internal/kvs"
	"github.com/tauraamui/maildew/internal/mail"
	"github.com/tauraamui/maildew/internal/mail/mock"
	"github.com/tauraamui/maildew/internal/storage/models"
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

	client := mail.NewClient(kvs.DB{})
	err = client.Connect(addr, mail.Account{
		Username: "username", Password: "password",
	})
	is.NoErr(err) // error connecting to imap server
	is.True(client != nil)
}

func TestClientMethodsProtectedIfNoConnection(t *testing.T) {
	is := is.New(t)

	client := mail.NewClient(kvs.DB{})
	is.True(client != nil)

	mb, err := client.FetchMailbox(mail.Account{}, "INBOX", true)
	is.Equal(mb, nil)
	is.True(err != nil)
	is.Equal(err, mail.ErrClientNotConnected)

	mbs, err := client.FetchAllMailboxes(mail.Account{})
	is.Equal(mbs, nil)
	is.True(err != nil)
	is.Equal(err, mail.ErrClientNotConnected)

	err = client.Close()
	is.True(err != nil)
	is.Equal(err, mail.ErrClientNotConnected)
}

func TestClientFetchMailboxes(t *testing.T) {
	is := is.New(t)

	client, err, cleanup := setupClientConnection()
	defer cleanup()
	is.NoErr(err)          // error connecting to imap server
	is.True(client != nil) // ensure client is not nil

	mailboxes, err := client.FetchAllMailboxes(mail.Account{Username: "username"})
	is.NoErr(err) // error fetching mailboxes

	is.True(len(mailboxes) > 0)
	is.Equal(mailboxes[0].Name(), "INBOX")
}

func TestClientFetchAllInboxMessages(t *testing.T) {
	t.Skip()
	is := is.New(t)

	client, err, cleanup := setupClientConnection()
	defer cleanup()
	is.NoErr(err)          // error connecting to imap server
	is.True(client != nil) // ensure client is not nil

	mb, err := client.FetchMailbox(mail.Account{Username: "username"}, "INBOX", true)
	is.NoErr(err) // error fetching mailbox of INBOX name
	msgs, err := mb.FetchAllMessages()
	is.NoErr(err)                                                               // error fetching inbox messages
	is.Equal(msgs, []mail.Message{{Subject: "A little message, just for you"}}) // list of messages does not match expected
}

func TestClientFetchAllInboxMessageUIDs(t *testing.T) {
	t.Skip()
	is := is.New(t)

	client, err, cleanup := setupClientConnection()
	defer cleanup()
	is.NoErr(err)
	is.True(client != nil)

	mb, err := client.FetchMailbox(mail.Account{Username: "username"}, "INBOX", true)
	is.NoErr(err) // error fetching mailbox of INBOX name
	uids, err := mb.FetchAllMessageUIDs()
	is.NoErr(err)
	is.Equal(uids, []mail.MessageUID{6})
}

func setupClientConnection() (mail.Client, error, func() error) {
	l, err := setupListener()
	if err != nil {
		return nil, err, nil
	}

	err, shutdown := startLocalServer(l, []models.Account{
		{Email: "username", Password: "password"},
	}...)
	if err != nil {
		return nil, err, func() error {
			// NOTE:(tauraamui) since starting the server was the cause of the error
			//       it's unnecessary to call the given shutdown callback
			return l.Close()
		}
	}

	addr := l.Addr().String()

	client := mail.NewClient(kvs.DB{})
	if err := client.Connect(addr, mail.Account{Username: "username", Password: "password"}); err != nil {
		return nil, err, func() error {
			errs := errgroup.I{}

			errs.Append(l.Close())
			errs.Append(shutdown())
			return errs.ToErrOrNil()
		}
	}

	return client, nil, func() error {
		errs := errgroup.I{}

		errs.Append(client.Close())
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

func startLocalServerWithBackend(l net.Listener, backend backend.Backend) (error, func() error) {
	s := server.New(backend)
	s.AllowInsecureAuth = true

	go s.Serve(l)

	return nil, s.Close
}

func startLocalServer(l net.Listener, users ...models.Account) (error, func() error) {
	mockBackend := mock.New()

	body := "From: contact@example.org\r\n" +
		"To: contact@example.org\r\n" +
		"Subject: A little message, just for you\r\n" +
		"Date: Wed, 11 May 2016 14:31:59 +0000\r\n" +
		"Message-ID: <0000000@localhost/>\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"Hi there :)"
	if users != nil {
		for _, u := range users {
			mockBackend.RegisterUser(u.Email, u.Password)
			mockBackend.CreateMailbox(u.Email, "INBOX")
			mockBackend.StoreMessage(u.Email, "INBOX", body)
		}
	}
	s := server.New(mockBackend)
	s.AllowInsecureAuth = true

	go s.Serve(l)

	return nil, s.Close
}
