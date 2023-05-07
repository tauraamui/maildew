package tui

import tea "github.com/charmbracelet/bubbletea"

func closeDialogCmd() func() tea.Msg {
	return func() tea.Msg { return closeDialogMsg{} }
}

type closeDialogMsg struct{}
