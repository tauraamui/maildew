package tui

import tea "github.com/charmbracelet/bubbletea"

func registerUserCmd(u, p string) func() tea.Msg {
	return func() tea.Msg {
		return registerUserMsg{u, p}
	}
}

type registerUserMsg struct {
	Username, Password string
}
