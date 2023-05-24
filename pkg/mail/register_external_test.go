package mail

import (
	"fmt"
	"net"
	"testing"

	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-imap/server"
	"github.com/tauraamui/maildew/internal/kvs"
	"github.com/tauraamui/maildew/internal/mail/mock"
	"github.com/tauraamui/maildew/internal/storage/models"
)

func TestPrintOutMailboxes(t *testing.T) {

	l, err := setupListener()

	backend := mock.New()
	backend.RegisterUser("username", "password")
	for i := 0; i < 20; i++ {
		backend.CreateMailbox("username", fmt.Sprintf("INBOX%d", i+1))
	}

	err, shutdown := startLocalServerWithBackend(l, backend)
	defer shutdown()

	if err != nil {
		t.Fatalf("unable to start local IMAP server: %v", err)
	}

	acc := Account{Username: "username", Password: "password"}
	cc, err := acquireClientConn(l.Addr().String(), acc, false)
	if err != nil {
		t.Fatal(err)
	}
	defer cc.Close()

	db, err := kvs.NewMemDB()
	if err != nil {
		t.Fatalf("unable to open in memory KVS: %v", err)
	}

	mbRepo := NewMailboxRepo(db)

	persistMailboxes(cc, mbRepo, acc)
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

	if users != nil {
		for _, u := range users {
			mockBackend.RegisterUser(u.Email, u.Password)
		}
	}
	s := server.New(mockBackend)
	s.AllowInsecureAuth = true

	go s.Serve(l)

	return nil, s.Close
}
