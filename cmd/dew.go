package main

import (
	"errors"
	"log"

	"github.com/tauraamui/dragondaemon/pkg/config"
	"github.com/tauraamui/dragondaemon/pkg/configdef"
	account "github.com/tauraamui/maildew/internal/storage"
	"github.com/tauraamui/maildew/internal/tui"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openSQLite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("mail.db"))
	if err != nil {
		log.Fatalf("unable to open DB: %v\n", err)
	}

	if err = db.AutoMigrate(account.Account{}); err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	// cfg, err := config.Load()
	// if err != nil {
	// 	panic(err)
	// }

	// config.ResolveRootKey()

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

	db := openSQLite()
	ar := account.AccountRepository{DB: db}
	tui.StartTea(ar)
}
