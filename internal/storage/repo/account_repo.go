package repo

import (
	"errors"

	"github.com/tauraamui/maildew/internal/storage"
	"github.com/tauraamui/maildew/internal/storage/models"
)

type Accounts struct {
	DB storage.DB
}

func (r Accounts) Save(user models.Account) error {
	return errors.New("saving user accounts not supported")
}
