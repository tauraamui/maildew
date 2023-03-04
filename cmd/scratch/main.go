package main

import (
	"log"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/tauraamui/maildew/internal/kvs"
)

type remoteAccount int
type remoteBox struct {
	ownerAccount int
}

type LocalAccountClone struct {
	RemoteRef          []byte // the mail servers remote account UUID equivilient
	LocalRef           uuid.UUID
	Username, Password string
}

type LocalBoxClone struct {
	Owner     uuid.UUID
	RemoteRef []byte
	LocalRef  kvs.UUID
	Name      string
}

func main() {
	log.Println("Experiment for trying to understand box storage in relation to mail and ownership.")

	db, err := kvs.NewMemDB()
	if err != nil {
		log.Fatalf("unable to load in memory DB: %v\n", err)
	}
	defer db.Close()

	accRepo := localAccountRepo{DB: db}
	if err := accRepo.Save(LocalAccountClone{
		RemoteRef: []byte("g594tjrio"),
		LocalRef:  uuid.New(),
		Username:  "localacccopy",
		Password:  "notrelevant",
	}); err != nil {
		log.Fatalf("unable to create local account in DB: %v\n", err)
	}

	if err := db.DumpToStdout(); err != nil {
		log.Fatalf("unable to output in memory DB to stdout: %v\n", err)
	}
}

// -------------------------------------------------------------------

type localAccountRepo struct {
	DB  kvs.DB
	seq *badger.Sequence
}

const noOwner uint32 = 0

func (r localAccountRepo) Save(acc LocalAccountClone) error {
	rowID, err := r.nextRowID()
	if err != nil {
		return err
	}

	return saveValue(r.DB, r.tableName(), kvs.RootOwner{}, rowID, acc)
}

func (r localAccountRepo) tableName() string {
	return "localaccounts"
}

func (r *localAccountRepo) nextRowID() (uint32, error) {
	if r.seq == nil {
		seq, err := r.DB.GetSeq([]byte(r.tableName()), 100)
		if err != nil {
			return 0, err
		}
		r.seq = seq
	}

	s, err := r.seq.Next()
	if err != nil {
		return 0, err
	}

	return uint32(s), nil
}

func saveValue(db kvs.DB, tableName string, ownerID kvs.UUID, rowID uint32, v interface{}) error {
	entries := kvs.ConvertToEntriesWithUUID(tableName, ownerID, rowID, v)
	for _, e := range entries {
		if err := kvs.Store(db, e); err != nil {
			return err
		}
	}

	return nil
}
