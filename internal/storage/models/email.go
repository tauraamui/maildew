package models

type Email struct {
	ID                    uint64 `mdb:"ignore"`
	Subject string
}
