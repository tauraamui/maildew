package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type dialogModel interface {
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

type errMsgModel struct {
	parent tea.Model
	err    error
}

func (m *errMsgModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m.parent, nil
		}
	}
	return m.parent, nil
}

func (m *errMsgModel) View() string {
	var b strings.Builder
	b.WriteString(m.err.Error())
	b.WriteRune('\n')
	b.WriteString(focusedOKButton)

	return dialogBoxStyle.Render(b.String())
	// return wrapInDialog(b.String(), m.windowSize, dialogBoxStyle)
}
