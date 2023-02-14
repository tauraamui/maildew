package mock

import (
	"errors"
	"fmt"
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
	fmt.Printf("mock users: %v\n", be.users[username])
}

func (be *Backend) Login(_ *imap.ConnInfo, username, password string) (backend.User, error) {
	user, ok := be.users[username]
	if ok && user.password == password {
		return user, nil
	}
	return nil, errors.New("bad username or password")
}

// NOTE:(tauraamui) having our own implementation of a mock IMAP server
// backend will allow us more control for things like number of
// mailboxes and the number of messages per mailbox
func New() *Backend {
	usr := &user{username: "username", password: "password"}

	body := "From: contact@example.org\r\n" +
		"To: contact@example.org\r\n" +
		"Subject: A little message, just for you\r\n" +
		"Date: Wed, 11 May 2016 14:31:59 +0000\r\n" +
		"Message-ID: <0000000@localhost/>\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"Hi there :)"

	usr.mailboxes = map[string]*mailbox{
		"INBOX": {
			name: "INBOX",
			user: usr,
			messages: []*message{
				{
					Uid:   6,
					Date:  time.Now(),
					Flags: []string{"\\Seen"},
					Size:  uint32(len(body)),
					Body:  []byte(body),
				},
			},
		},
	}

	return &Backend{
		users: map[string]*user{usr.username: usr},
	}
}
