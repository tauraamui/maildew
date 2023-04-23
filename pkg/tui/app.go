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
	return model{
		log:      log,
		imapAddr: addr,
		repos:    r,
		active:   initialRegisterAccountModel(log, addr, r),
	}
}

func (m model) Init() tea.Cmd {
	return m.active.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
	}
	return m.active.Update(msg)
}

func (m model) View() string {
	return m.active.View()
}
