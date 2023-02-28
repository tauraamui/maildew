package main

import (
	"errors"

	"github.com/tauraamui/maildew/internal/config"
	"github.com/tauraamui/maildew/internal/configdef"
	"github.com/tauraamui/maildew/internal/kvs"
	"github.com/tauraamui/maildew/internal/storage/repo"
	"github.com/tauraamui/maildew/internal/tui"
)

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
