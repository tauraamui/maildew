package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type dialogModel interface {
	SetPosition(x, y lipgloss.Position)
	Update(msg tea.Msg) tea.Cmd
	View() string
}

type errMsgModel struct {
	x, y   lipgloss.Position
	parent tea.Model
	err    error
}

func (m *errMsgModel) SetPosition(x, y lipgloss.Position) {
	m.x, m.y = x, y
}

func (m *errMsgModel) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return closeDialogCmd()
		}
	}
	return nil
}

func (m *errMsgModel) View() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n\n", errorStyle.Render("Error"))
	errStr := m.err.Error()
	if len(errStr) == 0 {
		errStr = "Something went wrong"
	} else {
		b.WriteString(strings.ToUpper(string(errStr[0])))
		b.WriteString(errStr[1:])
	}

	b.WriteString("\n\n")
	b.WriteString(focusedOKButton)

	return dialogBoxStyle.Copy().BorderForeground(lipgloss.Color("#874BFD")).Render(b.String())
}
