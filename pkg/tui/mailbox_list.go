package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauraamui/maildew/pkg/logging"
	"github.com/tauraamui/maildew/pkg/mail"
)

type mailboxListModel struct {
	log        logging.I
	windowSize tea.WindowSizeMsg
	mbrepo     mail.MailboxRepo
	acc        mail.Account
	list       []string
}

func initialMailboxListModel(log logging.I, mbrepo mail.MailboxRepo, acc mail.Account) *mailboxListModel {
	return &mailboxListModel{
		log:    log,
		list:   []string{},
		mbrepo: mbrepo,
		acc:    acc,
	}
}

func (m *mailboxListModel) Init() tea.Cmd {
	mboxes, err := m.mbrepo.FetchByOwner(m.acc.UUID)
	m.log.Debug().Msgf("retreived %d mailboxes for account %s", len(mboxes), m.acc.UUID.String())
	if err != nil {
		m.log.Error().Msgf("unable to fetch mailboxes: %w", err)
		// not sure how to handle this yet
	}

	for _, mbx := range mboxes {
		m.list = append(m.list, mbx.Name)
	}

	return nil
}

func (m *mailboxListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *mailboxListModel) View() string {
	sb := strings.Builder{}
	for _, mb := range m.list {
		sb.WriteString(mb)
		sb.WriteRune('\n')
	}

	return sb.String()
}
