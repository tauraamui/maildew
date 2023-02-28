package mail

// TODO:(tauraamui) implement client, decide on the interface etc.,

import (
	"errors"

	"github.com/emersion/go-imap"
	imapclient "github.com/emersion/go-imap/client"
	"github.com/tauraamui/maildew/internal/kvs"
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
	Connect(address string, account Account) error
	FetchMailbox(Account, string, bool) (Mailbox, error)
	FetchAllMailboxes(Account) ([]Mailbox, error)
	messageFetcher
	Close() error
}

type messageFetcher interface {
	fetchAllMessages(Account, Mailbox) ([]Message, error)
	fetchAllMessageUIDs(Account, Mailbox) ([]MessageUID, error)
}

func NewClient(db kvs.DB) Client {
	return &client{db: db}
}

func (c *client) Connect(ipaddress string, account Account) error {
	cc, err := imapclient.Dial(ipaddress)
	if err != nil {
		return err
	}

	if err := cc.Login(account.Username, account.Password); err != nil {
		return err
	}

	if c.loggedInAccounts == nil {
		c.loggedInAccounts = make(map[string]*imapclient.Client)
	}
	c.loggedInAccounts[account.Username] = cc

	return nil
}

func (c client) getConnection(username string) (*imapclient.Client, error) {
	client, ok := c.loggedInAccounts[username]
	if !ok {
		return nil, ErrClientNotConnected
	}

	return client, nil
}

type client struct {
	db               kvs.DB
	loggedInAccounts map[string]*imapclient.Client
}

func (c client) FetchMailbox(acc Account, name string, ro bool) (Mailbox, error) {
	conn, err := c.getConnection(acc.Username)
	if err != nil {
		return nil, err
	}

	m, err := conn.Select(name, ro)
	if err != nil {
		return nil, err
	}

	return newMailbox(c.db, m.Name, acc, c), nil
}

func (c client) FetchAllMailboxes(acc Account) ([]Mailbox, error) {
	conn, err := c.getConnection(acc.Username)
	if err != nil {
		return nil, err
	}

	mailboxesChan := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	mailboxes := []Mailbox{}
	go func() {
		done <- conn.List("", "*", mailboxesChan)
	}()

	for m := range mailboxesChan {
		mailboxes = append(mailboxes, newMailbox(c.db, m.Name, acc, c))
	}

	return mailboxes, <-done
}

func (c client) fetchAllMessages(acc Account, mailbox Mailbox) ([]Message, error) {
	conn, err := c.getConnection(acc.Username)
	if err != nil {
		return nil, err
	}

	mb, err := conn.Select(mailbox.Name(), true)
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
		done <- conn.Fetch(&seqset, []imap.FetchItem{imap.FetchEnvelope}, msgsch)
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

func (c client) fetchAllMessageUIDs(acc Account, mailbox Mailbox) ([]MessageUID, error) {
	conn, err := c.getConnection(acc.Username)
	if err != nil {
		return nil, err
	}

	mb, err := conn.Select(mailbox.Name(), true)
	if err != nil {
		return nil, err
	}

	seqset := imap.SeqSet{}
	seqset.AddRange(uint32(1), mb.Messages)

	msgsch := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	go func() {
		done <- conn.Fetch(&seqset, []imap.FetchItem{imap.FetchUid}, msgsch)
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
	if c.loggedInAccounts != nil {
		errs := errgroup.I{}
		for _, client := range c.loggedInAccounts {
			errs.Append(client.Logout())
			errs.Append(client.Close())
		}
		return errs.ToErrOrNil()
	}

	return ErrClientNotConnected
}
