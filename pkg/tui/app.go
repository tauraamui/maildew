package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tauraamui/maildew/pkg/logging"
	"github.com/tauraamui/maildew/pkg/mail"
)

var (
	focusedStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle            = lipgloss.NewStyle()
	focusedButton      = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton      = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	dialogContentStyle = lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).
				MarginRight(32)
	dialogBoxStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true, true, true, true).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(0, 1, 0)
)

type model struct {
	log        logging.I
	imapAddr   string
	repos      Repositories
	inputs     []textinput.Model
	focusIndex int
	windowSize tea.WindowSizeMsg
}

type Repositories struct {
	AccountRepo mail.AccountRepo
	MailboxRepo mail.MailboxRepo
	MessageRepo mail.MessageRepo
}

func Run(l logging.I, addr string, r Repositories) error {
	if _, err := tea.NewProgram(initialModel(l, addr, r), tea.WithAltScreen()).Run(); err != nil {
		return err
	}
	return nil
}

func initialModel(log logging.I, addr string, r Repositories) model {
	m := model{
		log:      log,
		imapAddr: addr,
		repos:    r,
		inputs:   make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()

		switch i {
		case 0:
			t.Placeholder = "Username"
			t.CharLimit = 48
			t.Focus()
		case 1:
			t.Placeholder = "Password"
			t.CharLimit = 48
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func registerUserCmd(u, p string) func() tea.Msg {
	return func() tea.Msg {
		return registerUserMsg{u, p}
	}
}

type registerUserMsg struct {
	Username, Password string
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case registerUserMsg:
		acc := mail.Account{Username: msg.Username, Password: msg.Password}
		mail.RegisterAccount(m.log, m.imapAddr, m.repos.AccountRepo, m.repos.MailboxRepo, m.repos.MessageRepo, acc)
	case tea.WindowSizeMsg:
		m.windowSize = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, registerUserCmd(m.inputs[0].Value(), m.inputs[1].Value())
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) resolveLongestInputValueWidth() int {
	i0w := lipgloss.Width(m.inputs[0].Value())
	i1w := lipgloss.Width(m.inputs[1].Value())
	if i0w == i1w {
		return i0w
	}

	if i1w > i0w {
		return i1w
	}

	return i0w
}

func (m model) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s", *button)

	textWidth := m.resolveLongestInputValueWidth()
	if textWidth >= 8 {
		dialogContentStyle.MarginRight(32 - (textWidth - 7))
	}
	content := dialogContentStyle.Render(b.String())
	dialog := lipgloss.Place(m.windowSize.Width, m.windowSize.Height,
		lipgloss.Center, lipgloss.Center, dialogBoxStyle.Render(content),
	)

	return dialog
}
