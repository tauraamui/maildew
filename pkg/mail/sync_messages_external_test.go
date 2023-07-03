package mail

import (
	"sort"
	"testing"

	"github.com/emersion/go-imap"
	"github.com/matryer/is"
)

func TestForEachMessage(t *testing.T) {
	is := is.New(t)

	mconn := &mockRemoteConnection{
		mailboxes: makeRemoteConnectionData(map[uint32]string{
			3353: "Cats & Dogs",
			5393: "Re: neighbour noise complaint",
		}),
	}

	fetchedSubjects := []string{}
	is.NoErr(forEachMessage(mconn, "INBOX", func(name string) error {
		fetchedSubjects = append(fetchedSubjects, name)
		return nil
	}))

	is.Equal(fetchedSubjects, []string{
		"Cats & Dogs",
		"Re: neighbour noise complaint",
	})
}

type mockRemoteConnection struct {
	selected          string
	mailboxes         map[string][]*imap.Message
	returnErrAfterNum int
	err               error
}

func sortedKeys(m map[string][]*imap.Message) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (mc mockRemoteConnection) List(ref, name string, ch chan *imap.MailboxInfo) error {
	defer close(ch)

	for i, name := range sortedKeys(mc.mailboxes) {
		if mc.returnErrAfterNum > 0 && i >= mc.returnErrAfterNum {
			return mc.err
		}
		ch <- &imap.MailboxInfo{Name: name}
	}
	return mc.err
}

func (mc *mockRemoteConnection) Select(name string, readOnly bool) (*imap.MailboxStatus, error) {
	mc.selected = name
	return &imap.MailboxStatus{
		Messages: uint32(len(mc.mailboxes[mc.selected])),
	}, nil
}

func (mc mockRemoteConnection) Fetch(seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
	defer close(ch)
	for _, msg := range mc.mailboxes[mc.selected] {
		ch <- msg
	}
	return nil
}

func (mc mockRemoteConnection) Close() error { return nil }

func makeRemoteConnectionData(inboxMessages map[uint32]string) map[string][]*imap.Message {
	mailboxesAndMessages := map[string][]*imap.Message{
		"INBOX":    {},
		"WORK":     {},
		"SHOPPING": {},
		"SPAM":     {},
		"DRAFTS":   {},
		"OTHER":    {},
		"MISC":     {},
		"SCHOOL":   {},
		"LIBRARY":  {},
		"MISC2":    {},
		"MISC3":    {},
		"MISC4":    {},
	}

	for uuid, name := range inboxMessages {
		mailboxesAndMessages["INBOX"] = append(mailboxesAndMessages["INBOX"], &imap.Message{
			Uid:      uuid,
			Envelope: &imap.Envelope{Subject: name},
		})
	}

	return mailboxesAndMessages
}
