package mail

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/maildew/internal/kvs"
)

func fetchByOwner[E Account | Mailbox | Message](db kvs.DB, tableName string, owner kvs.UUID) ([]E, error) {
	entries := make([]E, 1)

	blankEntries := kvs.ConvertToBlankEntriesWithUUID(tableName, owner, 0, entries[0])
	for _, ent := range blankEntries {
		// iterate over all stored values for this entry
		prefix := ent.PrefixKey()
		db.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			var rows uint32 = 0
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				if rows >= uint32(len(entries)) {
					entries = append(entries, *new(E))
				}
				item := it.Item()
				ent.RowID = rows
				if err := item.Value(func(val []byte) error {
					ent.Data = val
					return nil
				}); err != nil {
					return err
				}

				if err := kvs.LoadEntry(&entries[rows], ent); err != nil {
					return err
				}
				rows++
			}
			return nil
		})
	}
	return entries, nil
}
