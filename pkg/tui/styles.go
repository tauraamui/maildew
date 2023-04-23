package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	errorStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#E3242B"))
	noStyle             = lipgloss.NewStyle()
	focusedSubmitButton = focusedStyle.Copy().Render("[ Submit ]")
	focusedOKButton     = focusedStyle.Copy().Render("[ OK ]")
	blurredSubmitButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	dialogContentStyle  = lipgloss.NewStyle().Border(lipgloss.HiddenBorder()).
				MarginRight(32)
	dialogBoxStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true, true, true, true).
			Padding(0, 1, 0)
)

func wrapInDialog(c string, ws tea.WindowSizeMsg, stl lipgloss.Style) string {
	return lipgloss.Place(ws.Width, ws.Height,
		lipgloss.Center, lipgloss.Center, stl.Render(c),
	)
}
