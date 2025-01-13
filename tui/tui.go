package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func InitialModel() Model {
	return Model{
		choices:  []string{"Try X", "Try Y", "Try Z"},
		selected: make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) View() string {
	s := "TODO List\n\n"

	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor
		}

		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "X" // selected
		}

		s += fmt.Sprintf("%s %s %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit\n"

	return s
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Key press?
	case tea.KeyMsg:
		// What key
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}
