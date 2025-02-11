package cmd

import (
	"fmt"
	"github.com/pvlbzn/latai/internal/tui"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() {
	m, err := tui.NewTUIModel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseAllMotion()).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
