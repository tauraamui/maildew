package core

import (
	"os"

	"github.com/gtank/cryptopasta"
)

func ResolveRootKey() {
	key := cryptopasta.NewEncryptionKey()
	os.WriteFile(".mailkey", key[:], os.ModePerm)
}
