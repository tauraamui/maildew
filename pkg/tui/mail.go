package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	auth mode = iota
	list
)

// Model the entryui model definition
type Model struct {
	mode     mode
	auth     tea.Model
	quitting bool
}

// InitProject initialize the mailui model for your program
func InitMail() tea.Model {
	m := Model{
		mode: auth,
		auth: newAuthModel(),
	}
	return m
}

// Init run any intial IO on program start
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handle IO and commands
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	if m.mode == auth {
		return m.auth.Update(msg)
	}

	return nil, nil
}

// View return the text UI to be output to the terminal
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	switch m.mode {
	case auth:
		m.auth.View()
	}

	return "Nothing"
}
