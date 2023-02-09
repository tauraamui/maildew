package mail

// TODO:(tauraamui) implement client, decide on the interface etc.,

import (
	"errors"

	"github.com/emersion/go-imap"
	imapclient "github.com/emersion/go-imap/client"
	"github.com/tauraamui/maildew/internal/storage"
	"github.com/tauraamui/maildew/internal/storage/models"
	"github.com/tauraamui/xerror/errgroup"
)

var ErrClientNotConnected error = errors.New("client is not connected")

type MessageUID uint32

type MessageHeader struct {
	MessageUID
}

type Message struct {
	Subject string
}

type Client interface {
	Connect(address string, account models.Account) error
	FetchMailbox(string, bool) (Mailbox, error)
	FetchAllMailboxes() ([]Mailbox, error)
	messageFetcher
	messageDeleter
	Close() error
}

type messageFetcher interface {
	fetchAllMessages(Mailbox) ([]Message, error)
	fetchAllMessageUIDs(Mailbox) ([]MessageUID, error)
}

type messageDeleter interface {
	deleteMessageByUID() // TODO:(tauraamui) -> implement
}

func NewClient(db storage.DB) Client {
	return &client{db: db}
}

func (c *client) Connect(address string, account models.Account) error {
	cc, err := imapclient.Dial(address)
	if err != nil {
		return err
	}

	if err := cc.Login(account.Email, account.Password); err != nil {
		return err
	}

	c.client = cc
	c.account = account

	return nil
}

func (c client) checkConnected() error {
	if c.client == nil {
		return ErrClientNotConnected
	}

	// TODO:(tauraamui) -> figure out how to best check logged out channel is closed

	return nil
}

type client struct {
	db      storage.DB
	client  *imapclient.Client
	account models.Account
}

func (c client) FetchMailbox(name string, ro bool) (Mailbox, error) {
	if err := c.checkConnected(); err != nil {
		return nil, err
	}

	m, err := c.client.Select(name, ro)
	if err != nil {
		return nil, err
	}

	return newMailbox(c.db, m.Name, c.account, c), nil
}

func (c client) FetchAllMailboxes() ([]Mailbox, error) {
	if err := c.checkConnected(); err != nil {
		return nil, err
	}

	mailboxesChan := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	mailboxes := []Mailbox{}
	go func() {
		done <- c.client.List("", "*", mailboxesChan)
	}()

	for m := range mailboxesChan {
		mailboxes = append(mailboxes, newMailbox(c.db, m.Name, c.account, c))
	}

	return mailboxes, <-done
}

func (c client) fetchAllMessages(mailbox Mailbox) ([]Message, error) {
	if err := c.checkConnected(); err != nil {
		return nil, err
	}

	mb, err := c.client.Select(mailbox.Name(), true)
	if err != nil {
		return nil, err
	}

	from := uint32(1)
	to := mb.Messages
	seqset := imap.SeqSet{}
	seqset.AddRange(from, to)

	msgsch := make(chan *imap.Message, 1)
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

func (c client) fetchAllMessageUIDs(mailbox Mailbox) ([]MessageUID, error) {
	if err := c.checkConnected(); err != nil {
		return nil, err
	}

	mb, err := c.client.Select(mailbox.Name(), true)
	if err != nil {
		return nil, err
	}

	seqset := imap.SeqSet{}
	seqset.AddRange(uint32(1), mb.Messages)

	msgsch := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	go func() {
		done <- c.client.Fetch(&seqset, []imap.FetchItem{imap.FetchUid}, msgsch)
	}()

	headers := []MessageUID{}
	for msg := range msgsch {
		headers = append(headers, MessageUID(msg.Uid))
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return headers, nil
}

func (c client) Close() error {
	if c.client != nil {
		errs := errgroup.I{}
		errs.Append(c.client.Logout())
		errs.Append(c.client.Close())
		return errs.ToErrOrNil()
	}

	return ErrClientNotConnected
}
