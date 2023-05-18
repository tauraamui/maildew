package mail

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/maildew/internal/kvs"
)

func fetchByOwner[E any](db kvs.DB, tableName string, owner kvs.UUID) ([]E, error) {
	dest := []E{}

	typeRef := new(E)

	blankEntries := kvs.ConvertToBlankEntriesWithUUID(tableName, owner, 0, typeRef)
	for _, ent := range blankEntries {
		// iterate over all stored values for this entry
		prefix := ent.PrefixKey()
		db.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			var rows uint32 = 0
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				if len(dest) == 0 || rows >= uint32(len(dest)) {
					dest = append(dest, *new(E))
				}
				item := it.Item()
				ent.RowID = rows
				if err := item.Value(func(val []byte) error {
					ent.Data = val
					return nil
				}); err != nil {
					return err
				}

				if err := kvs.LoadEntry(&dest[rows], ent); err != nil {
					return err
				}
				rows++
			}
			return nil
		})
	}
	return dest, nil
}
