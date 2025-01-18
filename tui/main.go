package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() {
	m := NewTUIModel()

	if _, err := tea.NewProgram(&m).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
