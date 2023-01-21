package main

import (
	"log"

	"github.com/emersion/go-imap/backend/memory"
	"github.com/emersion/go-imap/server"
)

func main() {
	// create a memory backend
	be := memory.New()

	// create a new server
	s := server.New(be)
	s.Addr = ":1143"

	// since we will use this server for testing only we're allowing
	// auth over unencrypted connections
	s.AllowInsecureAuth = true

	log.Println("starting imap server at: localhost:" + s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
