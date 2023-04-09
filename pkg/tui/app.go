package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauraamui/maildew/pkg/logging"
	"github.com/tauraamui/maildew/pkg/mail"
)

type model struct {
	log                  logging.I
	imapAddr             string
	repos                Repositories
	windowSize           tea.WindowSizeMsg
	registerAccountModel tea.Model
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
	return model{
		log:                  log,
		imapAddr:             addr,
		repos:                r,
		registerAccountModel: initialRegisterAccountModel(log),
	}
}

func (m model) Init() tea.Cmd {
	return m.registerAccountModel.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
	}
	return m.registerAccountModel.Update(msg)
}

func (m model) View() string {
	return m.registerAccountModel.View()
}
