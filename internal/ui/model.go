// model.go
package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	progressBar  progress.Model
	progressChan <-chan float64
	quitting     bool
	err          error
	keys         *keyMap
}

// InitialModel sets up the initial state of the UI model.
func InitialModel(progressChan <-chan float64) model {
	m := model{
		progressBar:  progress.New(progress.WithDefaultGradient()),
		progressChan: progressChan,
		keys:         newKeyMap(),
	}
	return m
}

// Init initializes the model's command.
func (m model) Init() tea.Cmd {
	return m.progressTickCmd()
}

// progressTickCmd generates tick messages at a specified interval to update the UI.
func (m model) progressTickCmd() tea.Cmd {
	return tea.Every(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

// Update handles all incoming messages and updates the state of the model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progressBar.Width = msg.Width - 4 // Adjust width to fit screen
		return m, nil

	case tickMsg:
		select {
		case progress, ok := <-m.progressChan:
			if !ok {
				m.quitting = true
				return m, tea.Quit
			}
			m.progressBar.SetPercent(progress) // Correctly update progress
		default:
		}
		return m, m.progressTickCmd()

	case progress.FrameMsg:
		progressModel, cmd := m.progressBar.Update(msg)
		m.progressBar = progressModel.(progress.Model) // Correct type assertion
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	default:
		return m, nil
	}
}

// View renders the progress bar and the stop button.
func (m model) View() string {
	if m.quitting {
		return "Done!"
	}
	return m.progressBar.View()
}

type keyMap struct {
	stop key.Binding
}

func newKeyMap() *keyMap {
	return &keyMap{
		stop: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("ctrl+c/q", "stop"),
		),
	}
}
