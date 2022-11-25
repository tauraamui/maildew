package tui

import tea "github.com/charmbracelet/bubbletea"

type Auth struct{}

func (a Auth) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (a Auth) View() string {
	return ""
}
