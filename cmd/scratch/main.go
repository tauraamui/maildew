package main

import (
	"log"

	"github.com/google/uuid"
)

type remoteAccount int
type remoteBox struct {
	ownerAccount int
}

type localAccountClone struct {
	remoteRef uint // the mail servers remote account UUID equivilient
	localRef  uuid.UUID
}
type localBoxClone struct {
}

func main() {
	log.Println("Experiment for trying to understand box storage in relation to mail and ownership.")
}
