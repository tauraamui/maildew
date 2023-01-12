package tui

import (
	"math/rand"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tauraamui/maildew/internal/storage/models"
	"github.com/tauraamui/maildew/internal/storage/repo"
)

var marginStyle = lipgloss.NewStyle().Margin(1, 2)

type emailsmodel struct {
	er         repo.Emails
	windowSize tea.WindowSizeMsg
	list       list.Model
}

func populateRepoWithFake(er *repo.Emails) {
	for i := 0; i < 100; i++ {
		er.Save(&models.Email{Subject: randomString(5, 20)})
	}
}

func newEmailListModel(er repo.Emails) emailsmodel {
	populateRepoWithFake(&er)
	items := newEmailsList(er)
	m := emailsmodel{er: er, list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "test@account.com"
	return m
}

func (m emailsmodel) Init() tea.Cmd {
	return nil
}

func newEmailsList(er repo.Emails) []list.Item {
	emails, err := er.GetAll()
	if err != nil {
		panic(err)
	}

	return emailsToItems(emails)
}

func (m emailsmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		h, v := marginStyle.GetFrameSize()
		m.list.SetSize(m.windowSize.Width-h, m.windowSize.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			m.list, cmd = m.list.Update(msg)
		}
	}

	return m, cmd
}

func (m emailsmodel) View() string {
	return marginStyle.Render(m.list.View())
}

func emailsToItems(emails []models.Email) []list.Item {
	items := make([]list.Item, len(emails))
	for i, email := range emails {
		items[i] = list.Item(email)
	}
	return items
}

func randomString(min, max int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	n := rand.Intn(max-min) + min
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
