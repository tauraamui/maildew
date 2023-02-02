package models

import "strconv"

type Mailbox struct {
	ID   uint32 `mdb:"ignore"`
	UID  uint32
	Name string
}

func (m Mailbox) Title() string       { return m.Name }
func (m Mailbox) Description() string { return strconv.Itoa(int(m.ID)) }
func (m Mailbox) FilterValue() string { return m.Name }
