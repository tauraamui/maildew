package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauraamui/maildew/pkg/tui/constants"
)

type mode int

const (
	auth mode = iota
	createAccount
	list
)

type (
	createAccountMsg    struct{}
	authenticateUserMsg struct{}
)

// Model the entryui model definition
type Model struct {
	mode          mode
	createAccount tea.Model
	auth          tea.Model
	list          tea.Model
	quitting      bool
}

// InitProject initialize the mailui model for your program
func InitMail() tea.Model {
	m := Model{
		mode:          createAccount,
		createAccount: newCreateAccountModel(),
		list:          newListModel(),
	}
	return m
}

// Init run any intial IO on program start
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handle IO and commands
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
	case authenticateUserMsg:
		m.mode = list
	case createAccountMsg:
		m.mode = list
	}

	var cmd tea.Cmd
	switch m.mode {
	case createAccount:
		m.createAccount, cmd = m.createAccount.Update(msg)
		return m, cmd
	case list:
		m.list, cmd = m.list.Update(msg)
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
	case createAccount:
		return m.createAccount.View()
	case list:
		return m.list.View()
	}

	return "Nothing"
}
