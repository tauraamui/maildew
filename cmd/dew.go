package main

import (
	"errors"

	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/maildew/internal/config"
	"github.com/tauraamui/maildew/internal/configdef"
	account "github.com/tauraamui/maildew/internal/storage"
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

	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		panic(err)
	}

	ar := account.AccountRepository{DB: db}
	tui.StartTea(cfg, ar)
}
