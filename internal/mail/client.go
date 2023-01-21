package mail

// TODO: implement client, decide on the interface etc.,

type Client interface {
	Mailboxes() interface{}
}

func Connect(email, password string) Client {
	return client{}
}

type client struct{}

func (c client) Mailboxes() interface{} {
	return nil
}
