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

	msgc := make(chan *imap.Message, 1)
	errc := make(chan error)
	defer close(errc)
	go func() {
		errc <- conn.Fetch(buildSequence(mb.Messages), []imap.FetchItem{imap.FetchEnvelope}, msgc)
	}()

	// if an error is encountered, msgc should be closed automatically
	for msg := range msgc {
		if msg == nil || msg.Envelope == nil {
			continue
		}

		if err := callback(msg.Envelope.Subject); err != nil {
			return err
		}
	}

	return <-errc
}

func buildSequence(msgs uint32) *imap.SeqSet {
	from := uint32(1)
	to := msgs
	seqset := imap.SeqSet{}
	seqset.AddRange(from, to)
	return &seqset
}
