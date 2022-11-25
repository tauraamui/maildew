package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type emailsmodel struct {
}

func newListModel() emailsmodel {
	return emailsmodel{}
}

func (m emailsmodel) Init() tea.Cmd {
	return textinput.Blink
}

func (m emailsmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m emailsmodel) View() string {
	return "emails list"
}
