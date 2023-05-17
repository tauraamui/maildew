package main

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
	"github.com/tauraamui/maildew/internal/kvs"
	"github.com/tauraamui/maildew/pkg/logging"
	"github.com/tauraamui/maildew/pkg/mail"
)

type remoteAccount int
type remoteBox struct {
	ownerAccount int
}

type LocalAccountClone struct {
	RemoteRef          []byte // the mail servers remote account UUID equivilient
	LocalRef           kvs.UUID
	Username, Password string
}

type LocalBoxClone struct {
	RemoteRef []byte
	LocalRef  kvs.UUID
	Name      string
}

type LocalMessageClone struct {
	ID        uint32 `mdb:"ignore"`
	RemoteRef []byte
	LocalRef  kvs.UUID
	Subject   string
}

func main() {
	f, err := os.OpenFile("maildew.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	log := logging.New(logging.Options{Level: logging.DEBUG, Writer: f})
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.Info().Msg("Experiment for trying to understand box storage in relation to mail and ownership.")

	db, err := kvs.NewMemDB()
	if err != nil {
		log.Fatal().Msgf("unable to load in memory DB: %v\n", err)
	}
	defer db.Close()

	mail.RegisterAccount(log, "", mail.NewAccountRepo(db), mail.NewMailboxRepo(db), mail.NewMessageRepo(db), &mail.Account{
		Username: "test@place.com",
		Password: "fakepassword",
	})

	accRepo := localAccountRepo{DB: db}
	boxRepo := localBoxRepo{DB: db}
	msgRepo := localMessageRepo{DB: db}

	acc := LocalAccountClone{
		RemoteRef: []byte("g594tjrio"),
		LocalRef:  uuid.New(),
		Username:  "localacccopy",
		Password:  "notrelevant",
	}

	if err := accRepo.Save(acc); err != nil {
		log.Fatal().Msgf("unable to create local account in DB: %v\n", err)
	}

	inbox := LocalBoxClone{
		RemoteRef: []byte("whkotyor"),
		LocalRef:  uuid.New(),
		Name:      "INBOX",
	}

	if err := boxRepo.Save(acc.LocalRef, inbox); err != nil {
		log.Fatal().Msgf("unable to store local box in DB: %v\n", err)
	}

	junk := LocalBoxClone{
		RemoteRef: []byte("rgoiergo"),
		LocalRef:  uuid.New(),
		Name:      "JUNK",
	}

	if err := boxRepo.Save(acc.LocalRef, junk); err != nil {
		log.Fatal().Msgf("unable to store local box in DB: %v\n", err)
	}

	testMsgs := []LocalMessageClone{
		{
			RemoteRef: []byte("ykirgire"),
			LocalRef:  uuid.New(),
			Subject:   "RE: Testing testing 123!",
		}, {
			RemoteRef: []byte("hrthrtr"),
			LocalRef:  uuid.New(),
			Subject:   "Probably some spam.",
		}, {

			RemoteRef: []byte("wiergoerg"),
			LocalRef:  uuid.New(),
			Subject:   "Hot New Deals!",
		},
	}

	for _, testMsg := range testMsgs {
		if err := msgRepo.Save(inbox.LocalRef, testMsg); err != nil {
			log.Fatal().Msgf("unable to store local message in DB: %v\n", err)
		}
	}

	if err := msgRepo.Save(junk.LocalRef, LocalMessageClone{
		RemoteRef: []byte("rjgerigf"),
		LocalRef:  uuid.New(),
		Subject:   "Definitely some spam",
	}); err != nil {
		log.Fatal().Msgf("unable to store local message in DB: %v\n", err)
	}

	/*
		if err := db.DumpToStdout(); err != nil {
			log.Fatalf("unable to output in memory DB to stdout: %v\n", err)
		}
	*/

	inboxmsgs, err := msgRepo.GetMessages(inbox.LocalRef)
	if err != nil {
		log.Fatal().Msgf("unable to acquire messages owned by mailbox %s from DB: %v\n", inbox.LocalRef, err)
	}

	for _, msg := range inboxmsgs {
		fmt.Printf("INBOX MSG: %+v\n", msg)
	}

	junkmsgs, err := msgRepo.GetMessages(junk.LocalRef)
	if err != nil {
		log.Fatal().Msgf("unable to acquire messages owned by mailbox %s from DB: %v\n", junk.LocalRef, err)
	}

	for _, msg := range junkmsgs {
		fmt.Printf("JUNK MSG: %+v\n", msg)
	}
}

// -------------------------------------------------------------------

type localBoxRepo struct {
	DB  kvs.DB
	seq *badger.Sequence
}

func (r localBoxRepo) Save(owner kvs.UUID, box LocalBoxClone) error {
	rowID, err := r.nextRowID()
	if err != nil {
		return err
	}

	return saveValue(r.DB, r.tableName(), owner, rowID, box)
}

func (r localBoxRepo) tableName() string {
	return "localboxes"
}

func (r *localBoxRepo) nextRowID() (uint32, error) {
	if r.seq == nil {
		seq, err := r.DB.GetSeq([]byte(r.tableName()), 1)
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

// -------------------------------------------------------------------

type localMessageRepo struct {
	DB  kvs.DB
	seq *badger.Sequence
}

func (r localMessageRepo) Save(ownerID kvs.UUID, msg LocalMessageClone) error {
	rowID, err := r.nextRowID()
	if err != nil {
		return err
	}

	return saveValue(r.DB, r.tableName(), ownerID, rowID, msg)
}

func (r localMessageRepo) GetMessages(ownerID kvs.UUID) ([]LocalMessageClone, error) {
	messages := []LocalMessageClone{}

	blankEntries := kvs.ConvertToBlankEntriesWithUUID(r.tableName(), ownerID, 0, LocalMessageClone{})
	for _, ent := range blankEntries {
		// iterate over all stored values for this entry
		prefix := ent.PrefixKey()
		if err := r.DB.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			var rows uint32 = 0
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				// be very clear with our append conditions
				if len(messages) == 0 || rows >= uint32(len(messages)) {
					messages = append(messages, LocalMessageClone{
						ID: rows,
					})
				}
				item := it.Item()
				ent.RowID = rows
				if err := item.Value(func(val []byte) error {
					ent.Data = val
					return nil
				}); err != nil {
					return err
				}
				if err := kvs.LoadEntry(&messages[rows], ent); err != nil {
					return err
				}
				rows++
			}

			return nil
		}); err != nil {
			return nil, err
		}
	}

	return messages, nil
}

func (r localMessageRepo) tableName() string {
	return "localmessages"
}

func (r *localMessageRepo) nextRowID() (uint32, error) {
	if r.seq == nil {
		seq, err := r.DB.GetSeq([]byte(r.tableName()), 1)
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

// -------------------------------------------------------------------

type localAccountRepo struct {
	DB  kvs.DB
	seq *badger.Sequence
}

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
		seq, err := r.DB.GetSeq([]byte(r.tableName()), 1)
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
