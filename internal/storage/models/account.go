package models

import "strconv"

type Account struct {
	ID                    uint32 `mdb:"ignore"`
	Nick, Email, Password string
}

func (e *Account) SetID(id uint32)    { e.ID = id }
func (e *Account) Ref() interface{}   { return e }
func (e Account) Title() string       { return e.Email }
func (e Account) Description() string { return strconv.Itoa(int(e.ID)) }
func (e Account) FilterValue() string { return e.Email }
