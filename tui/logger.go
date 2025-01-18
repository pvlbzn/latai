package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type LoggerComponent struct {
	width int
	logs  []log

	showLast int
}

type log struct {
	id      int
	time    string
	message string
}

func newLog(id int, time, message string) log {
	return log{
		id:      id,
		time:    time,
		message: message,
	}
}

func NewLoggerComponent(width int) *LoggerComponent {
	return &LoggerComponent{
		width:    width,
		logs:     []log{},
		showLast: 5,
	}
}

func (m *LoggerComponent) WithShowLast(n int) *LoggerComponent {
	m.showLast = n
	return m
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
		// Render in a stack style.
		for i := len(m.logs); i > 0; i-- {
			if len(m.logs)-i > m.showLast {
				break
			}

			log := m.logs[i-1]
			row := fmt.Sprintf("%d\t%s\t%s", log.id, log.time, log.message)
			messages = append(messages, rowStyle.Render(row))
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
	id := len(m.logs)
	m.logs = append(m.logs, newLog(id, "now", message))
}
