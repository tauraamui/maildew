package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	auth mode = iota
	list
)

type (
	authenticateUserMsg struct{}
)

// Model the entryui model definition
type Model struct {
	mode     mode
	auth     tea.Model
	list     tea.Model
	quitting bool
}

// InitProject initialize the mailui model for your program
func InitMail() tea.Model {
	m := Model{
		mode: auth,
		auth: newAuthModel(),
		list: newListModel(),
	}
	return m
}

// Init run any intial IO on program start
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handle IO and commands
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case authenticateUserMsg:
		m.mode = list
	}

	var cmd tea.Cmd
	switch m.mode {
	case auth:
		m.auth, cmd = m.auth.Update(msg)
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
	case auth:
		return m.auth.View()
	case list:
		return m.list.View()
	}

	return "Nothing"
}
