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
		mailboxes: makeRemoteConnectionData(),
	}

	is.NoErr(forEachMessage(mconn, "INBOX", func(name string) error {
		return nil
	}))

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
	for _, msg := range mc.mailboxes[mc.selected] {
		ch <- msg
	}
	close(ch)
	return nil
}

func (mc mockRemoteConnection) Close() error { return nil }

func makeRemoteConnectionData() map[string][]*imap.Message {
	return map[string][]*imap.Message{
		"INBOX": {
			{
				Uid: 321,
				Envelope: &imap.Envelope{
					Subject: "Test inbox message",
				},
			},
			{
				Uid: 5940,
				Envelope: &imap.Envelope{
					Subject: "Car insurance ad",
				},
			},
			{
				Uid: 623943,
				Envelope: &imap.Envelope{
					Subject: "Order is 15 days late",
				},
			},
			{
				Uid: 65096,
				Envelope: &imap.Envelope{
					Subject: "Feel happy!",
				},
			},
		},
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
}
