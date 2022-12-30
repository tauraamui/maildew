package models

type Account struct {
	ID                    uint64 `mdb:"ignore"`
	Nick, Email, Password string
}
