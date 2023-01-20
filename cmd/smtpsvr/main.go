package main

import (
	"log"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/tauraamui/maildew/internal/mail"
)

func main() {
	svr := mail.MockServer{}

	s := smtp.NewServer(&svr)

	s.Addr = ":1025"
	s.Domain = "localhost"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("starting smtp server at:", s.Domain+s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
