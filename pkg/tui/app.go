package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauraamui/maildew/pkg/logging"
	"github.com/tauraamui/maildew/pkg/mail"
)

type model struct {
	log        logging.I
	imapAddr   string
	repos      Repositories
	windowSize tea.WindowSizeMsg
	active     tea.Model
}

type Repositories struct {
	AccountRepo mail.AccountRepo
	MailboxRepo mail.MailboxRepo
	MessageRepo mail.MessageRepo
}

func Run(l logging.I, addr string, r Repositories) error {
	if _, err := tea.NewProgram(initialModel(l, addr, r), tea.WithAltScreen()).Run(); err != nil {
		return err
	}
	return nil
}

func initialModel(log logging.I, addr string, r Repositories) model {
	m := model{
		log:      log,
		imapAddr: addr,
		repos:    r,
	}

	m.active = initialRegisterAccountModel(log, m, addr, r)
	return m
}

func (m model) Init() tea.Cmd {
	return m.active.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	hasActive := m.active != nil
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
	case tea.KeyMsg:
		if !hasActive {
			switch msg.String() {
			case "ctrl+c", "esc":
				return m, tea.Quit
			}
		}
	case openMailboxListMsg:
		m.active = msg.mailboxListModel
		return m, m.active.Init()
	}

	if hasActive {
		return m.active.Update(msg)
	}

	return m, nil
}

func (m model) View() string {
	if m.active == nil {
		return "maildew app has no active view at this time"
	}
	return m.active.View()
}
