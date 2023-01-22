package mock

import (
	"errors"

	"github.com/emersion/go-imap/backend"
)

type user struct {
	username  string
	password  string
	mailboxes map[string]*mailbox
}

func (u *user) Username() string {
	return u.username
}

func (u *user) ListMailboxes(subscribed bool) (mailboxes []backend.Mailbox, err error) {
	for _, mailbox := range u.mailboxes {
		if subscribed && !mailbox.subscribed {
			continue
		}

		mailboxes = append(mailboxes, mailbox)
	}
	return
}

func (u *user) GetMailbox(name string) (mailbox backend.Mailbox, err error) {
	mailbox, ok := u.mailboxes[name]
	if !ok {
		err = errors.New("no such mailbox")
	}
	return
}

func (u *user) CreateMailbox(name string) error {
	if _, ok := u.mailboxes[name]; ok {
		return errors.New("Mailbox already exists")
	}

	u.mailboxes[name] = &mailbox{name: name, user: u}
	return nil
}

func (u *user) DeleteMailbox(name string) error {
	if name == "INBOX" {
		return errors.New("Cannot delete INBOX")
	}
	if _, ok := u.mailboxes[name]; !ok {
		return errors.New("No such mailbox")
	}

	delete(u.mailboxes, name)
	return nil
}

func (u *user) RenameMailbox(existingName, newName string) error {
	mbox, ok := u.mailboxes[existingName]
	if !ok {
		return errors.New("No such mailbox")
	}

	u.mailboxes[newName] = &mailbox{
		name:     newName,
		messages: mbox.messages,
		user:     u,
	}

	mbox.messages = nil

	if existingName != "INBOX" {
		delete(u.mailboxes, existingName)
	}

	return nil
}

func (u *user) Logout() error {
	return nil
}
