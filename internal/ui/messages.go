package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// tickMsg is used to trigger periodic updates.
type tickMsg struct{}

// progressTickCmd generates tick messages at a specified interval to update the UI.
func progressTickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}
