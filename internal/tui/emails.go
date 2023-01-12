package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauraamui/maildew/internal/storage/repo"
)

type emailsmodel struct {
	er         repo.Emails
	windowSize tea.WindowSizeMsg
}

func newListModel() emailsmodel {
	return emailsmodel{}
}

func (m emailsmodel) Init() tea.Cmd {
	return nil
}

func (m emailsmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
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
