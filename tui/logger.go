package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type LoggerComponent struct {
	width int
	logs  []string
}

func NewLoggerComponent(width int) *LoggerComponent {
	return &LoggerComponent{
		width: width,
		logs:  []string{},
	}
}

func (m *LoggerComponent) Init() tea.Cmd {
	return nil
}

func (m *LoggerComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *LoggerComponent) View() string {
	// Log view principal style.
	s := lg.NewStyle().
		BorderStyle(lg.NormalBorder()).
		BorderForeground(lg.Color("241")).
		Width(m.width)

	header := lg.NewStyle().
		Bold(true).
		PaddingLeft(1).
		Render("Log Messages")

	separator := lg.NewStyle().
		Foreground(lg.Color("240")).
		Render(strings.Repeat("─", m.width))

	rowStyle := lg.NewStyle().
		PaddingLeft(1)

	var messages []string
	if len(m.logs) == 0 {
		messages = append(messages, rowStyle.Render("No records..."))
	} else {
		for _, log := range m.logs {
			messages = append(messages, rowStyle.Render(log))
		}
	}

	return s.Render(lg.JoinVertical(
		lg.Top,
		header,
		separator,
		strings.Join(messages, "\n"),
	))
}

func (m *LoggerComponent) Push(message string) {
	m.logs = append(m.logs, message)
}
