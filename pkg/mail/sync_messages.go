package mail

import (
	"github.com/emersion/go-imap"
	"github.com/tauraamui/maildew/pkg/logging"
)

func SyncMessages(
	log logging.I,
	conn RemoteConnection,
	msgr MessageRepo,
	mb Mailbox,
) error {
	return syncMailboxMessages(conn, msgr, mb)
}

func syncMailboxMessages(conn RemoteConnection, msgr MessageRepo, mb Mailbox) error {
	if err := forEachMessage(conn, mb.Name, func(name string) error {
		if err := msgr.Save(mb.UUID, Message{}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func forEachMessage(conn RemoteConnection, mailboxName string, callback func(name string) error) error {
	mb, err := conn.Select(mailboxName, true)
	if err != nil {
		return err
	}

	msgsc := make(chan *imap.Message, 1)
	errc := make(chan error, 1)
	// NOTE:(tauraamui) the implementation of fetch already closes our receiver channel on completion
	go func() {
		errc <- conn.Fetch(buildSequence(mb.Messages), []imap.FetchItem{imap.FetchEnvelope}, msgsc)
	}()

	for {
		select {
		case err := <-errc:
			if err != nil {
				return err
			}
		case msg, more := <-msgsc:
			if err := callback(msg.Envelope.Subject); err != nil {
				return err
			}

			if !more {
				return nil
			}
		}
	}
}

func buildSequence(msgs uint32) *imap.SeqSet {
	from := uint32(1)
	to := msgs
	seqset := imap.SeqSet{}
	seqset.AddRange(from, to)
	return &seqset
}
