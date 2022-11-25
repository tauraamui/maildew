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
	quitting bool
}

// InitProject initialize the mailui model for your program
func InitMail() tea.Model {
	m := Model{mode: auth}
	return m
}

// Init run any intial IO on program start
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handle IO and commands
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

// View return the text UI to be output to the terminal
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	return "Auth"
}
