package mock

import (
	"errors"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
)

type Backend struct {
	users map[string]*user
}

func (be *Backend) Login(_ *imap.ConnInfo, username, password string) (backend.User, error) {
	user, ok := be.users[username]
	if ok && user.password == password {
		return user, nil
	}
	return nil, errors.New("bad username or password")
}

func New() *Backend {
	usr := &user{username: "username", password: "password"}

	return &Backend{
		users: map[string]*user{usr.username: usr},
	}
}
