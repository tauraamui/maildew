package main

import (
	"github.com/tauraamui/maildew/pkg/logging"
)

func main() {
	log := logging.New(logging.Options{Level: logging.DEBUG})
	log.Debug().Msg("test debug message")
	log.Info().Msg("MAILDEW REGISTRATION")
	//conn := mail.RegisterAccount()
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
