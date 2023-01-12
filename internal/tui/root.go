package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauraamui/maildew/internal/storage/repo"
)

type rootmodel struct {
	windowSize tea.WindowSizeMsg
	focusIndex int
	accounts   tea.Model
	emails     tea.Model
}

func newRootModel(ar repo.Accounts, er repo.Emails) rootmodel {
	return rootmodel{
		accounts: newAccountsListModel(ar),
		emails:   newEmailListModel(er),
	}
}

func (m rootmodel) Init() tea.Cmd {
	return nil
}

func (m rootmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			if m.focusIndex == 0 {
				m.accounts, cmd = m.accounts.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
			m.emails, cmd = m.emails.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m rootmodel) View() string {
	return m.accounts.View() + m.emails.View()
}
