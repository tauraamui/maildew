package models

import "strconv"

type Email struct {
	ID      uint32 `mdb:"ignore"`
	Subject string
}

func (e *Email) SetID(id uint32)    { e.ID = id }
func (e *Email) Ref() interface{}   { return e }
func (e Email) Title() string       { return e.Subject }
func (e Email) Description() string { return strconv.Itoa(int(e.ID)) }
func (e Email) FilterValue() string { return e.Subject }
