package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauraamui/maildew/internal/storage/models"
	"github.com/tauraamui/maildew/internal/storage/repo"
)

type accountsmodel struct {
	ar         repo.Accounts
	windowSize tea.WindowSizeMsg
	list       list.Model
}

func newAccountsListModel(ar repo.Accounts) accountsmodel {
	ar.Save(&models.Account{
		Nick:     "First Test Account",
		Email:    "firsttest@account.com",
		Password: "jiejfofowif",
	})
	ar.Save(&models.Account{
		Nick:     "Second Test Account",
		Email:    "secondtest@account.com",
		Password: "grejgioregoirjgi",
	})
	items := newAccountsList(ar)
	m := accountsmodel{ar: ar, list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "Accounts"
	return m
}

func (m accountsmodel) Init() tea.Cmd {
	return nil
}

func newAccountsList(ar repo.Accounts) []list.Item {
	accounts, err := ar.GetAll()
	if err != nil {
		panic(err)
	}

	return accountsToItems(accounts)
}

func (m accountsmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		h, v := marginStyle.GetFrameSize()
		m.list.SetSize(m.windowSize.Width-h, m.windowSize.Height-v)
	}

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m accountsmodel) View() string {
	return marginStyle.Render(m.list.View())
}

func accountsToItems(accounts []models.Account) []list.Item {
	items := make([]list.Item, len(accounts))
	for i, account := range accounts {
		items[i] = list.Item(account)
	}
	return items
}
