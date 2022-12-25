package main

import (
	"log"

	"github.com/tauraamui/maildew/internal/account"
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
	account.SetupLocal()
	db := openSQLite()
	ar := account.AccountRepository{DB: db}
	tui.StartTea(ar)
}
