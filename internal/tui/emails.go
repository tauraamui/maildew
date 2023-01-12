package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauraamui/maildew/internal/storage/models"
	"github.com/tauraamui/maildew/internal/storage/repo"
)

type emailsmodel struct {
	er         repo.Emails
	windowSize tea.WindowSizeMsg
	list       list.Model
}

func populateRepoWithFake(er *repo.Emails) {
	for i := 0; i < 100; i++ {
		er.Save(&models.Email{Subject: fmt.Sprintf("Fake email %d", i)})
	}
}

func newEmailListModel(er repo.Emails) emailsmodel {
	return emailsmodel{er: er}
}

func (m emailsmodel) Init() tea.Cmd {
	populateRepoWithFake(&m.er)
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
