package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauraamui/maildew/pkg/logging"
)

type errMsgModel struct {
	log        logging.I
	parent     tea.Model
	windowSize tea.WindowSizeMsg
	err        error
}

func (m errMsgModel) Init() tea.Cmd {
	return nil
}

func (m errMsgModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
	case tea.KeyMsg:
		m.log.Debug().Msg("key update from error message model")
		switch msg.String() {
		case "enter":
			return m.parent, nil
		}
	}
	return m, nil
}

func (m errMsgModel) View() string {
	var b strings.Builder
	b.WriteString(m.err.Error())
	b.WriteRune('\n')
	b.WriteString(focusedOKButton)

	return wrapInDialog(b.String(), m.windowSize, dialogBoxStyle)
}
