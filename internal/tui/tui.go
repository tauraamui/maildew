package tui

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tauraamui/maildew/internal/configdef"
	"github.com/tauraamui/maildew/internal/storage/repo"
)

// StartTea the entry point for the UI. Initializes the model.
func StartTea(cfg configdef.Values, ur repo.Accounts) {
	if cfg.Debug {
		if f, err := tea.LogToFile("debug.log", "debug"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			defer func() {
				err = f.Close()
				if err != nil {
					log.Fatal(err)
				}
			}()
		}
	}

	p := tea.NewProgram(InitMail(ur), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
