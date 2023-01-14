package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauraamui/maildew/internal/storage/repo"
	// account "github.com/tauraamui/maildew/internal/storage"
)

type (
	mode   int
	status int
)

const (
	rootStatus status = iota
	createAccountStatus
)

type Model struct {
	status        status
	root          tea.Model
	createAccount tea.Model
	windowSize    tea.WindowSizeMsg
	quitting      bool
}

func InitMail(ar repo.Accounts, er repo.Emails) tea.Model {
	m := Model{
		root:          newRootModel(ar, er),
		createAccount: newCreateAccountModel(ar),
	}
	return &m
}

// Init run any intial IO on program start
func (m *Model) Init() tea.Cmd {
	m.status = rootStatus

	return nil
}

// Update handle IO and commands
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
	}

	var cmd tea.Cmd
	cmds := []tea.Cmd{}
	switch m.status {
	case rootStatus:
		m.root, cmd = m.root.Update(msg)
		cmds = append(cmds, cmd)
	case createAccountStatus:
		m.createAccount, cmd = m.createAccount.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View return the text UI to be output to the terminal
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	switch m.status {
	case rootStatus:
		return m.root.View()
	case createAccountStatus:
		return m.createAccount.View()
	}

	// should probably just panic here
	return "A problem or bug is occurring, this text should never appear... Check the logs, or something."
}
