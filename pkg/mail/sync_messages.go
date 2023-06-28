package mail

import "github.com/tauraamui/maildew/pkg/logging"

func SyncMessages(
	log logging.I,
	conn RemoteConnection,
	msgr MessageRepo,
	mb Mailbox,
) error {
	return syncMailboxMessages(conn, msgr, mb)
}

func syncMailboxMessages(conn RemoteConnection, msgr MessageRepo, mb Mailbox) error {
	return nil
}
