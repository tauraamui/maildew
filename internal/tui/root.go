package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tauraamui/maildew/internal/storage/repo"
)

var (
	modelStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
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
		m.accounts, cmd = m.accounts.Update(msg)
		cmds = append(cmds, cmd)

		m.emails, cmd = m.emails.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.focusIndex++
			if m.focusIndex > 1 {
				m.focusIndex = 0
			}
		}
	}

	if m.focusIndex == 0 {
		m.accounts, cmd = m.accounts.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	m.emails, cmd = m.emails.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m rootmodel) View() string {
	var b strings.Builder
	if m.focusIndex == 0 {
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(m.accounts.View()), modelStyle.Render(m.emails.View())))
		return b.String()
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(m.accounts.View()), focusedModelStyle.Render(m.emails.View())))
	return b.String()
}
