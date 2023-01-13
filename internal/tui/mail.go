package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauraamui/maildew/internal/storage/repo"
	// account "github.com/tauraamui/maildew/internal/storage"
)

type mode int

const (
	auth mode = iota
	rootMode
	createAccountMode
	emailsListMode
)

// Model the entryui model definition
type Model struct {
	mode          mode
	root          tea.Model
	createAccount tea.Model
	windowSize    tea.WindowSizeMsg
	quitting      bool
}

// InitProject initialize the mailui model for your program
func InitMail(ar repo.Accounts, er repo.Emails) tea.Model {
	m := Model{
		// root: newRootModel(ar, er),
		createAccount: newCreateAccountModel(ar),
	}
	return &m
}

// Init run any intial IO on program start
func (m *Model) Init() tea.Cmd {
	m.mode = createAccountMode
	// m.mode = rootMode

	return nil
}

// Update handle IO and commands
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case switchModeMsg:
		m.mode = msg.mode
		return m, nil
	case tea.WindowSizeMsg:
		m.windowSize = msg
	}

	var cmd tea.Cmd
	switch m.mode {
	case createAccountMode:
		m.createAccount, cmd = m.createAccount.Update(msg)
		return m, cmd
	case rootMode:
		m.root, cmd = m.root.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View return the text UI to be output to the terminal
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	switch m.mode {
	case createAccountMode:
		return m.createAccount.View()
	case rootMode:
		return m.root.View()
	}

	// should probably just panic here
	return "Nothing"
}
