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
	log := logging.New(logging.Options{Level: logging.DEBUG})
	log.Debug().Msg("test debug message")
	log.Info().Msg("MAILDEW REGISTRATION")

	l, err := setupListener()
	if err != nil {
		log.Fatal().Msgf("unable to start localhost TCP listener: %v", err)
	}

	username, password := os.Getenv("MD_USERNAME"), os.Getenv("MD_PASSWORD")
	backend := mock.New()
	//backend.RegisterUser(username, password)
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

	acc := mail.Account{
		Username: username,
		Password: password,
	}

	if err := mail.RegisterAccount(log, l.Addr().String(), accRepo, mbRepo, msgRepo, acc); err != nil {
		log.Fatal().Msgf("failed to register new account: %v", err)
	}

	db.DumpToStdout()

	l.Close()
	shutdown()

	if err := tui.Run(); err != nil {
		log.Fatal().Msgf("failed to load TUI: %v", err)
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
