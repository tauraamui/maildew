package mail

type Remote interface {
	GetMailboxes() error
}

type Syncable interface {
	Sync() error
}
