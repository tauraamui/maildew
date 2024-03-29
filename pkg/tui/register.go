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

type registerAccountModel struct {
	log        logging.I
	parent     tea.Model
	imapAddr   string
	r          Repositories
	inputs     []textinput.Model
	focusIndex int
	windowSize tea.WindowSizeMsg
	errDialog  dialogModel
}

func initialRegisterAccountModel(log logging.I, parent tea.Model, imapAddr string, r Repositories) registerAccountModel {
	m := registerAccountModel{
		log:      log,
		parent:   parent,
		imapAddr: imapAddr,
		r:        r,
		inputs:   make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()

		switch i {
		case 0:
			t.Placeholder = "Username"
			t.CharLimit = 48
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
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

func (m registerAccountModel) Init() tea.Cmd {
	return textinput.Blink
}

type errorMessageMsg struct {
	err error
}

type returnToParentMsg struct {
	cc  mail.RemoteConnection
	acc mail.Account
}

func registerAccountCmd(l logging.I, imapAddr string, u, p string, r Repositories) func() tea.Msg {
	return func() tea.Msg {
		acc := mail.Account{Username: u, Password: p}
		cc, err := mail.RegisterAccount(l, imapAddr, r.AccountRepo, r.MailboxRepo, &acc, mail.ResolveClientConnector(imapAddr, acc))
		if err != nil {
			return errorMessageMsg{err}
		}
		return returnToParentMsg{cc, acc}
	}
}

type openMailboxListMsg struct {
	mailboxListModel tea.Model
}

func openMailboxListCmd(l logging.I, mbRepo mail.MailboxRepo, acc mail.Account) func() tea.Msg {
	return func() tea.Msg {
		return openMailboxListMsg{
			mailboxListModel: initialMailboxListModel(l, mbRepo, acc),
		}
	}
}

func (m registerAccountModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case returnToParentMsg:
		return m.parent, openMailboxListCmd(m.log, m.r.MailboxRepo, msg.acc)
	case errorMessageMsg:
		m.errDialog = &errMsgModel{
			parent: m,
			err:    msg.err,
		}
		return m, nil
	case closeDialogMsg:
		m.errDialog = nil
	case tea.WindowSizeMsg:
		m.windowSize = msg
	case tea.KeyMsg:
		if m.errDialog != nil {
			return m, m.errDialog.Update(msg)
		}
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, registerAccountCmd(m.log, m.imapAddr, m.inputs[0].Value(), m.inputs[1].Value(), m.r)
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

func (m registerAccountModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m registerAccountModel) resolveLongestInputValueWidth() int {
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

func (m registerAccountModel) View() string {
	var b strings.Builder

	b.WriteString("Register new account\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredSubmitButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedSubmitButton
	}
	fmt.Fprintf(&b, "\n\n%s", *button)

	textWidth := m.resolveLongestInputValueWidth()
	if textWidth >= 18 {
		dialogContentStyle.MarginRight(32 - (textWidth - 17))
	}

	var style = dialogBoxStyle
	if m.errDialog == nil {
		style = style.Copy().BorderForeground(lipgloss.Color("#874BFD"))
	}
	bg := wrapInDialog(dialogContentStyle.Render(b.String()), m.windowSize, style)
	if m.errDialog != nil {
		fg := m.errDialog.View()
		x := (m.windowSize.Width / 2) - (lipgloss.Width(fg) / 2)
		y := (m.windowSize.Height / 2) - (lipgloss.Height(fg) / 2)

		m.errDialog.SetPosition(lipgloss.Position(x), lipgloss.Position(y))
		fg = m.errDialog.View()
		return placeOverlay(x, y, fg, bg, false)
	}

	return bg
}
