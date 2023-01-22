package mail

// TODO:(tauraamui) implement client, decide on the interface etc.,

import (
	"github.com/emersion/go-imap"
	imapclient "github.com/emersion/go-imap/client"
)

type Mailbox struct {
	Name string
}

type Message struct {
	Subject string
}

type Client interface {
	Mailboxes() ([]Mailbox, error)
	Messages(Mailbox) ([]Message, error)
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
		done <- c.client.List("", "*", mailboxesChan)
	}()

	for m := range mailboxesChan {
		mailboxes = append(mailboxes, Mailbox{m.Name})
	}

	return mailboxes, <-done
}

func (c client) Messages(mailbox Mailbox) ([]Message, error) {
	mb, err := c.client.Select(mailbox.Name, true)
	if err != nil {
		return nil, err
	}

	from := uint32(1)
	to := mb.Messages
	if mb.Messages > 3 {
		to = mb.Messages - 3
	}

	seqset := imap.SeqSet{}
	seqset.AddRange(from, to)

	msgsch := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.client.Fetch(&seqset, []imap.FetchItem{imap.FetchEnvelope}, msgsch)
	}()

	msgs := []Message{}
	for msg := range msgsch {
		msgs = append(msgs, Message{Subject: msg.Envelope.Subject})
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return msgs, nil
}

func (c client) Close() error {
	return c.client.Close()
}
