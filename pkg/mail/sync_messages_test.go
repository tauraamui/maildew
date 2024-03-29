package mail

import (
	"errors"
	"sort"
	"testing"

	"github.com/emersion/go-imap"
	"github.com/matryer/is"
)

func TestForEachMessageWithFetchingSuccessful(t *testing.T) {
	is := is.New(t)

	mconn := &mockRemoteConnection{
		mailboxes: makeRemoteConnectionData(map[uint32]string{
			3353: "Cats & Dogs",
			5393: "Re: neighbour noise complaint",
			3283: "Library - Book Overdue!",
		}),
	}

	fetchedSubjects := []string{}
	is.NoErr(forEachMessage(mconn, "INBOX", func(name string) error {
		fetchedSubjects = append(fetchedSubjects, name)
		return nil
	}))

	is = is.NewRelaxed(t)
	is.True(contains(fetchedSubjects, "Cats & Dogs"))
	is.True(contains(fetchedSubjects, "Re: neighbour noise complaint"))
	is.True(contains(fetchedSubjects, "Library - Book Overdue!"))
}

func TestForEachMessageWithFetchingEncountersAnImmediateErrorOnStart(t *testing.T) {
	is := is.New(t)

	mconn := &mockRemoteConnection{
		mailboxes: makeRemoteConnectionData(map[uint32]string{
			3353: "Cats & Dogs",
			5393: "Re: neighbour noise complaint",
			3283: "Library - Book Overdue!",
		}),
		fetch: func(mc mockRemoteConnection, seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
			defer close(ch)
			return errors.New("failed to initialise fetching process")
		},
	}

	fetchedSubjects := []string{}
	err := forEachMessage(mconn, "INBOX", func(name string) error {
		fetchedSubjects = append(fetchedSubjects, name)
		return nil
	})
	is.True(err != nil)
	is.Equal(err.Error(), "failed to initialise fetching process")

	is = is.NewRelaxed(t)
	is.Equal(len(fetchedSubjects), 0)
}

func TestForEachMessageWithFetchingEncountersAnErrorDuringProcess(t *testing.T) {
	is := is.New(t)

	mconn := &mockRemoteConnection{
		mailboxes: makeRemoteConnectionData(map[uint32]string{
			3353: "Cats & Dogs",
			5393: "Re: neighbour noise complaint",
			3283: "Library - Book Overdue!",
		}),
		fetch: func(mc mockRemoteConnection, seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
			defer close(ch)

			msgs := mc.mailboxes[mc.selected]
			ch <- msgs[0]
			ch <- msgs[1]
			return errors.New("failed during fetching process")
		},
	}

	fetchedSubjects := []string{}
	err := forEachMessage(mconn, "INBOX", func(name string) error {
		fetchedSubjects = append(fetchedSubjects, name)
		return nil
	})
	is.True(err != nil)
	is.Equal(err.Error(), "failed during fetching process")

	is = is.NewRelaxed(t)
	is.Equal(len(fetchedSubjects), 2)

	if doesNotContain(fetchedSubjects, "Library - Book Overdue!") {
		is.True(contains(fetchedSubjects, "Cats & Dogs"))
		is.True(contains(fetchedSubjects, "Re: neighbour noise complaint"))
	} else if doesNotContain(fetchedSubjects, "Cats & Dogs") {
		is.True(contains(fetchedSubjects, "Library - Book Overdue!"))
		is.True(contains(fetchedSubjects, "Re: neighbour noise complaint"))
	} else if doesNotContain(fetchedSubjects, "Re: neighbour noise complaint") {
		is.True(contains(fetchedSubjects, "Cats & Dogs"))
		is.True(contains(fetchedSubjects, "Library - Book Overdue!"))
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func doesNotContain(s []string, str string) bool {
	return !contains(s, str)
}

type mockRemoteConnection struct {
	fetch             fetchFunc
	selected          string
	mailboxes         map[string][]*imap.Message
	returnErrAfterNum int
	err               error
}

type fetchFunc func(mc mockRemoteConnection, seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error

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

func defaultMockFetchFunc(mc mockRemoteConnection, seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
	defer close(ch)
	for _, msg := range mc.mailboxes[mc.selected] {
		ch <- msg
	}
	return nil
}

func (mc mockRemoteConnection) Fetch(seqset *imap.SeqSet, items []imap.FetchItem, ch chan *imap.Message) error {
	if mc.fetch == nil {
		mc.fetch = defaultMockFetchFunc
	}
	return mc.fetch(mc, seqset, items, ch)
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
