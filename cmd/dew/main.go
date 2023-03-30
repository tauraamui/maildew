package main

import (
	"os"

	"github.com/tauraamui/maildew/internal/kvs"
	"github.com/tauraamui/maildew/pkg/logging"
	"github.com/tauraamui/maildew/pkg/mail"
)

func main() {
	log := logging.New(logging.Options{Level: logging.DEBUG})
	log.Debug().Msg("test debug message")
	log.Info().Msg("MAILDEW REGISTRATION")

	db, err := kvs.NewMemDB()
	if err != nil {
		log.Fatal().Msgf("unable to open in memory KVS: %v", err)
	}

	accRepo := mail.NewAccountRepo(db)
	mbRepo := mail.NewMailboxRepo(db)
	msgRepo := mail.NewMessageRepo(db)

	acc := mail.Account{
		Username: os.Getenv("MD_USERNAME"),
		Password: os.Getenv("MD_PASSWORD"),
	}

	if err := mail.RegisterAccount(log, accRepo, mbRepo, msgRepo, acc); err != nil {
		log.Fatal().Msgf("failed to register new account: %v", err)
	}
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
