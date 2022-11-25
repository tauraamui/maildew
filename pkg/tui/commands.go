package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func createAccountCmd(nick, email, pass string) tea.Cmd {
	return func() tea.Msg {
		log.Printf("NICK: %s, EMAIL: %s, PASS: %s\n", nick, email, pass)
		return createAccountMsg{}
	}
}
