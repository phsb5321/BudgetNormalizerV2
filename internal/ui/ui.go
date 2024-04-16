// internal/ui/ui.go
package ui

import (
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type tickMsg time.Time
type progressMsg float64

type model struct {
	progress progress.Model
	percent  float64
	done     chan bool
	info     chan string
	suffix   string
}

func initialModel(done chan bool, info chan string) model {
	return model{
		progress: progress.New(progress.WithDefaultGradient()),
		done:     done,
		info:     info,
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil
	case progressMsg:
		m.percent = float64(msg)
		cmd := m.progress.SetPercent(m.percent)
		return m, cmd
	case tickMsg:
		if m.percent >= 1.0 {
			return m, tea.Quit
		}
		select {
		case <-m.done:
			m.percent = 1.0
		case info := <-m.info:
			m.suffix = info
		default:
		}
		return m, tickCmd()
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	default:
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.ViewAs(m.percent) + " " + m.suffix + "\n\n" +
		pad + helpStyle("Processing data... (Press 'q' or 'ctrl+c' to quit)")
}

const tickInterval = 50

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func StartUI(done chan bool, info chan string) {
	m := initialModel(done, info)
	p := tea.NewProgram(m)
	if err := func() error {
		_, err := p.Run()
		return err
	}(); err != nil {
		log.Println("Error running UI:", err)
	}
}
