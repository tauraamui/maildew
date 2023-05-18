package mock

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
)

type Backend struct {
	users map[string]*user
}

func (be *Backend) RegisterUser(username, password string) {
	if be.users == nil {
		be.users = map[string]*user{}
	}
	be.users[username] = &user{username: username, password: password}
}

func (be *Backend) Login(_ *imap.ConnInfo, username, password string) (backend.User, error) {
	user, ok := be.users[username]
	if ok && user.password == password {
		return user, nil
	}
	return nil, errors.New("bad username or password")
}

type LocalBackend interface {
	backend.Backend
	RegisterUser(username, password string)
	CreateMailbox(username, mbname string) error
	StoreMessage(username, mbname, body string)
}

// NOTE:(tauraamui) having our own implementation of a mock IMAP server
// backend will allow us more control for things like number of
// mailboxes and the number of messages per mailbox
func New() LocalBackend {
	usr := &user{username: "username", password: "password"}

	return &xbackend{
		users: map[string]*user{usr.username: usr},
	}
}

type xbackend struct {
	users map[string]*user
}

func (bk *xbackend) Login(_ *imap.ConnInfo, username, password string) (backend.User, error) {
	user, ok := bk.users[username]
	if ok && user.password == password {
		return user, nil
	}
	return nil, errors.New("bad username or password")
}

func (bk *xbackend) RegisterUser(username, password string) {
	if bk.users == nil {
		bk.users = map[string]*user{}
	}
	bk.users[username] = &user{username: username, password: password, mailboxes: map[string]*mailbox{}}
}

func (bk *xbackend) CreateMailbox(username, mbname string) error {
	mbname = strings.ToUpper(mbname)
	usr, ok := bk.users[username]
	if !ok {
		return fmt.Errorf("unable to create mailbox for non-existant user %s", username)
	}

	usr.mailboxes[mbname] = &mailbox{
		name:     mbname,
		user:     usr,
		messages: []*message{},
	}
	return nil
}

func (bk *xbackend) StoreMessage(username, mbname, body string) {
	msgs := bk.users[username].mailboxes[mbname].messages
	msgs = append(msgs, &message{
		Uid:   uint32(len(msgs)),
		Date:  time.Now(),
		Flags: []string{"\\Seen"},
		Size:  uint32(len(body)),
		Body:  []byte(body),
	})
}
