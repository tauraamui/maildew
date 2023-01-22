package mail

// TODO: implement client, decide on the interface etc.,

import (
	"github.com/emersion/go-imap"
	imapclient "github.com/emersion/go-imap/client"
)

type Mailbox struct {
	Name string
}

type Client interface {
	Mailboxes() ([]Mailbox, error)
	Close() error
}

func Connect(address, email, password string) (Client, error) {
	c, err := imapclient.Dial(address)
	if err := c.Login(email, password); err != nil {
		return nil, err
	}
	return client{
		client: c,
	}, err
}

type client struct {
	client *imapclient.Client
}

func (c client) Mailboxes() ([]Mailbox, error) {
	mailboxesChan := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	mailboxes := []Mailbox{}
	go func() {
		done <- c.client.List("", "", mailboxesChan)
	}()

	for m := range mailboxesChan {
		mailboxes = append(mailboxes, Mailbox{m.Name})
	}

	return mailboxes, <-done
}

func (c client) Close() error {
	return c.client.Close()
}
