// ui.go
package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// StartUI initializes and starts the Bubble Tea program for the UI.
func StartUI(progressChan chan float64) {
	p := tea.NewProgram(InitialModel(progressChan), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "UI encountered an error: %v\n", err)
		os.Exit(1)
	}
}
