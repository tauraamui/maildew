package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func createAccountCmd(nick, email, pass string) tea.Cmd {
	return func() tea.Msg {
		return createAccountMsg{nick, email, pass}
	}
}
