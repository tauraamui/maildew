package mail

import (
	"errors"
	"io"
	"io/ioutil"
	"log"

	"github.com/emersion/go-smtp"
)

type MockSMTPServer struct{}

func (ms *MockSMTPServer) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

type Session struct{}

// Authenticate the user using SASL PLAIN.
func (s *Session) AuthPlain(username, password string) error {
	if username != "username" || password != "password" {
		return errors.New("invalid username or password")
	}
	return nil
}

// Set return path for currently processed message.
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	log.Println("mail from:", from)
	return nil
}

// Add recipient for currently processed message.
func (s *Session) Rcpt(to string) error {
	log.Println("rcpt to:", to)
	return nil
}

// Set currently processed message contents and send it.
//
// r must be consumed before Data returns.
func (s *Session) Data(r io.Reader) error {
	if b, err := ioutil.ReadAll(r); err != nil {
		return err
	} else {
		log.Println("data:", string(b))
	}
	return nil
}

// Discard currently processed message.
func (s *Session) Reset() {
}

// Free all resources associated with session.
func (s *Session) Logout() error {
	return nil
}
