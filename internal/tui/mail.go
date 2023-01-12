package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauraamui/maildew/internal/storage/repo"
	// account "github.com/tauraamui/maildew/internal/storage"
)

type mode int

const (
	auth mode = iota
	createAccountMode
	emailsListMode
)

type (
	createAccountMsg struct {
		nick, email, password string
	}
	updateFocusedInputsMsg struct {
		index int
	}
	resetFormMsg        struct{}
	authenticateUserMsg struct{}
)

// Model the entryui model definition
type Model struct {
	ar            repo.Accounts
	mode          mode
	createAccount tea.Model
	emailList     tea.Model
	windowSize    tea.WindowSizeMsg
	auth          tea.Model
	quitting      bool
}

// InitProject initialize the mailui model for your program
func InitMail(ar repo.Accounts, er repo.Emails) tea.Model {
	m := Model{
		ar:            ar,
		createAccount: newCreateAccountModel(ar),
		emailList:     newEmailListModel(er),
	}
	return &m
}

// Init run any intial IO on program start
func (m *Model) Init() tea.Cmd {
	m.mode = createAccountMode
	// accs, _ := m.ar.GetAccounts()
	// if len(accs) == 0 {
	// 	m.mode = createAccountMode
	// }
	return nil
}

// Update handle IO and commands
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
	}

	var cmd tea.Cmd
	switch m.mode {
	case createAccountMode:
		m.createAccount, cmd = m.createAccount.Update(msg)
		return m, cmd
	case emailsListMode:
		// m.list, cmd = m.list.Update(msg)
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
	case emailsListMode:
		// return m.list.View()
		return ""
	}

	return "Nothing"
}
