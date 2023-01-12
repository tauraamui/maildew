package models

import "strconv"

type Account struct {
	ID                    uint64 `mdb:"ignore"`
	Nick, Email, Password string
}

func (e Account) Title() string       { return e.Email }
func (e Account) Description() string { return strconv.Itoa(int(e.ID)) }
func (e Account) FilterValue() string { return e.Email }
