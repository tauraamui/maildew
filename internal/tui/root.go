package tui

import (
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

type focusStatus int

const (
	accountsFocus focusStatus = iota
	emailsFocus
)

type rootmodel struct {
	status       focusStatus
	windowSize   tea.WindowSizeMsg
	accountsList tea.Model
	emailsList   tea.Model
}

func newRootModel(ar repo.Accounts, er repo.Emails) rootmodel {
	return rootmodel{
		status:       accountsFocus,
		accountsList: newAccountsListModel(ar),
		emailsList:   newEmailListModel(er),
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
		m.accountsList, cmd = m.accountsList.Update(msg)
		cmds = append(cmds, cmd)

		m.emailsList, cmd = m.emailsList.Update(msg)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			switch m.status {
			case accountsFocus:
				m.status = emailsFocus
			case emailsFocus:
				m.status = accountsFocus
			}
		}
	}

	switch m.status {
	case accountsFocus:
		m.accountsList, cmd = m.accountsList.Update(msg)
		cmds = append(cmds, cmd)
	case emailsFocus:
		m.emailsList, cmd = m.emailsList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

//  FIX:(tauraamui) whenever the longer list is scrolled through when it gets to the last page
//		 it shows entry numbers for entries which are no longer rendered or present
//		 this rendering issue resolves itself when the focus is switched back and forth.

func (m rootmodel) View() string {
	var accountsView string
	var emailsView string

	switch m.status {
	case accountsFocus:
		accountsView = focusedModelStyle.Render(m.accountsList.View())
		emailsView = modelStyle.Render(m.emailsList.View())
	case emailsFocus:
		accountsView = modelStyle.Render(m.accountsList.View())
		emailsView = focusedModelStyle.Render(m.emailsList.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, accountsView, emailsView)
}
