package main

import (
	"fmt"
	"net"
	"os"

	"github.com/emersion/go-imap/backend"
	"github.com/emersion/go-imap/server"
	"github.com/tauraamui/maildew/internal/kvs"
	"github.com/tauraamui/maildew/internal/mail/mock"
	"github.com/tauraamui/maildew/internal/storage/models"
	"github.com/tauraamui/maildew/pkg/logging"
	"github.com/tauraamui/maildew/pkg/mail"
	"github.com/tauraamui/maildew/pkg/tui"
)

func main() {
	f, err := os.OpenFile("maildew.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log := logging.New(logging.Options{Level: logging.DEBUG, Writer: f})
	log.Info().Msg("MAILDEW v0.0.0a")

	l, err := setupListener()
	if err != nil {
		log.Fatal().Msgf("unable to start localhost TCP listener: %v", err)
	}

	backend := mock.New()
	backend.RegisterUser("username", "password")
	for i := 0; i < 11; i++ {
		backend.CreateMailbox("username", fmt.Sprintf("INBOX%d", i+1))
	}

	body := "From: contact@example.org\r\n" +
		"To: contact@example.org\r\n" +
		"Subject: A little message, just for you\r\n" +
		"Date: Wed, 11 May 2016 14:31:59 +0000\r\n" +
		"Message-ID: <0000000@localhost/>\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"Hi there :)"
	backend.StoreMessage("username", "INBOX1", body)

	err, shutdown := startLocalServerWithBackend(l, backend)
	if err != nil {
		log.Fatal().Msgf("unable to start local IMAP server: %v", err)
	}

	db, err := kvs.NewMemDB()
	if err != nil {
		log.Fatal().Msgf("unable to open in memory KVS: %v", err)
	}

	accRepo := mail.NewAccountRepo(db)
	mbRepo := mail.NewMailboxRepo(db)
	msgRepo := mail.NewMessageRepo(db)

	if err := tui.Run(
		log,
		l.Addr().String(),
		tui.Repositories{
			AccountRepo: accRepo,
			MailboxRepo: mbRepo,
			MessageRepo: msgRepo,
		}); err != nil {
		log.Fatal().Msgf("failed to load TUI: %v", err)
	}

	db.DumpTo(f)

	l.Close()
	shutdown()
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

/*
func main() {
	if err := config.DefaultCreator().Create(); err != nil {
		if err != nil {
			if !errors.Is(err, configdef.ErrConfigAlreadyExists) {
				panic(err)
			}
		}
	}

	cfg, err := config.DefaultResolver().Resolve()
	if err != nil {
		panic(err)
	}

	db, err := kvs.NewMemDB()
	if err != nil {
		panic(err)
	}

	ar := repo.Accounts{DB: db}
	defer ar.Close()

	er := repo.Emails{DB: db}
	defer er.Close()

	tui.StartTea(cfg, ar, er)

	db.DumpToStdout()
}
*/
