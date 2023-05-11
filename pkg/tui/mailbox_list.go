package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauraamui/maildew/pkg/logging"
	"github.com/tauraamui/maildew/pkg/mail"
)

type mailboxListModel struct {
	log        logging.I
	windowSize tea.WindowSizeMsg
	mbrepo     mail.MailboxRepo
	list       [10]string
}

func initialMailboxListModel(log logging.I, mbrepo mail.MailboxRepo) mailboxListModel {
	return mailboxListModel{
		log:    log,
		mbrepo: mbrepo,
	}
}

func (m mailboxListModel) Init() tea.Cmd {
	return nil
}

func (m mailboxListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return nil, nil
}

func (m mailboxListModel) View() string {
	return "mailbox list"
}
