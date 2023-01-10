package tui

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tauraamui/maildew/internal/storage/models"
	"github.com/tauraamui/maildew/internal/storage/repo"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	successStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	errorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))

	focusedButton  = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton  = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	dialogBoxStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true, true, true, true).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0)
)

type createaccountmodel struct {
	ar         repo.Accounts
	focusIndex int
	windowSize tea.WindowSizeMsg
	viewport   viewport.Model
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
	success    string
	err        error
}

func newCreateAccountModel(ar repo.Accounts) createaccountmodel {
	m := createaccountmodel{
		ar:     ar,
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Nickname"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Email"
			t.CharLimit = 64
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m createaccountmodel) Init() tea.Cmd {
	return textinput.Blink
}

func (m createaccountmodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)

	cmds := []tea.Cmd{vpCmd}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > textinput.CursorHide {
				m.cursorMode = textinput.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].SetCursorMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.err = nil
				cmds = append(cmds, createAccountCmd(m.inputs[0].Value(), m.inputs[1].Value(), m.inputs[2].Value()))
				cmds = append(cmds, clearFieldsResetFormCmd())
				return m, tea.Batch(cmds...)
			}

			// Cycle indexes
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

			return m, updateFocusedInputsCmd(m.focusIndex)
		}
	case tea.WindowSizeMsg:
		m.windowSize = msg
	case updateFocusedInputsMsg:
		fi := msg.index
		inputCmds := make([]tea.Cmd, len(m.inputs))
		for i := 0; i < len(m.inputs); i++ {
			if i == fi {
				// Set focused state
				inputCmds[i] = m.inputs[i].Focus()
				m.inputs[i].PromptStyle = focusedStyle
				m.inputs[i].TextStyle = focusedStyle
				continue
			}
			// Remove focused state
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = noStyle
			m.inputs[i].TextStyle = noStyle
		}
		cmds = append(cmds, inputCmds...)

		return m, tea.Batch(cmds...)
	case clearFieldsResetFormMsg:
		for i := 0; i < len(m.inputs); i++ {
			m.inputs[i].SetValue("")
		}
		m.focusIndex = 0

		return m, updateFocusedInputsCmd(m.focusIndex)
	case createAccountMsg:
		acc := models.Account{
			Email:    msg.email,
			Nick:     msg.nick,
			Password: msg.password,
		}
		if err := m.ar.Save(&acc); err != nil {
			m.err = err
		} else {
			m.success = fmt.Sprintf("account created, user ID: %d", acc.ID)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *createaccountmodel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m createaccountmodel) View() string {
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
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))
	if m.err != nil {
		b.WriteRune('\n')
		b.WriteRune('\n')
		b.WriteString(errorStyle.Render(fmt.Sprintf("%s...", m.err.Error())))
	}

	if len(m.success) > 0 {
		b.WriteRune('\n')
		b.WriteRune('\n')
		b.WriteString(successStyle.Render(fmt.Sprintf("%s...", m.success)))
	}

	ui := lipgloss.JoinVertical(lipgloss.Center, b.String())
	dialog := lipgloss.Place(m.windowSize.Width, m.windowSize.Height,
		lipgloss.Center, lipgloss.Center, dialogBoxStyle.Render(ui),
	)

	return dialog
}
