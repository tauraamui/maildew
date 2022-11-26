package account

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Nick, Email, Password string
}

type Repository interface {
	CreateAccount(nick, email, pass string) error
	GetAccounts() ([]Account, error)
}

type AccountRepository struct {
	DB *gorm.DB
}

func (r AccountRepository) CreateAccount(nick, email, pass string) error {
	account := Account{Nick: nick, Email: email, Password: pass}
	result := r.DB.Create(&account)
	return result.Error
}

func (r AccountRepository) GetAccounts() ([]Account, error) {
	var accounts []Account
	result := r.DB.Find(&accounts)
	return accounts, result.Error
}
