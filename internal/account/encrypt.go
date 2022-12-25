package account

import (
	"os"

	"github.com/gtank/cryptopasta"
)

func SetupLocal() {
	key := cryptopasta.NewEncryptionKey()
	os.WriteFile(".mailkey", key[:], os.ModePerm)
}
