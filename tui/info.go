package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/pvlbzn/latai/provider"
	"strings"
)

// InfoComponent is an informational component which displays
// data about the latency measurement.
type InfoComponent struct {
	width int
	info  []*modelInfo
	model selectedModel
}

type modelInfo struct {
	row          *TableRow
	measurements []*provider.Metric
}

// NewInfoComponent creates an instance of a new panel which displays
// data about models and latency details. Info panel shows details about
// models based on their row ID. Models are stored internally, and they are
// mapped onto their row IDs in the way it matches with table row IDs.
func NewInfoComponent(width int, rows []*TableRow) *InfoComponent {
	info := make([]*modelInfo, len(rows))
	for _, row := range rows {
		info = append(info, &modelInfo{row: row, measurements: make([]*provider.Metric, 0)})
	}

	return &InfoComponent{
		width: width,
		info:  info,
	}
}

func (s *InfoComponent) Init() tea.Cmd {
	return nil
}

func (s *InfoComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case modelSelectedMsg:
		s.model = selectedModel{
			id:           msg.id,
			providerName: msg.providerName,
			vendorName:   msg.vendorName,
			modelFamily:  msg.modelFamily,
			modelName:    msg.modelName,
		}
	}

	return s, tea.Batch()
}

func (s *InfoComponent) View() string {
	container := lg.NewStyle().
		BorderStyle(lg.NormalBorder()).
		BorderForeground(lg.Color("241")).
		Width(s.width)

	var header string
	if len(s.model.modelName) != 0 {
		header = lg.NewStyle().
			Bold(true).
			PaddingLeft(1).
			Render(fmt.Sprintf("Info: %s | %s | %s | %s", s.model.modelFamily, s.model.modelName, s.model.providerName, s.model.vendorName))
	} else {
		header = lg.NewStyle().
			Bold(true).
			PaddingLeft(1).
			Render("Info")
	}

	separator := lg.NewStyle().
		Foreground(lg.Color("240")).
		Render(strings.Repeat("─", s.width))

	rowStyle := lg.NewStyle().
		PaddingLeft(1)

	var info string
	if len(s.model.modelName) != 0 {
		info = rowStyle.
			Foreground(lg.Color("240")).
			Render(s.model.modelName)
	} else {
		info = rowStyle.
			Foreground(lg.Color("240")).
			Render("Select a model...")
	}

	return container.Render(lg.JoinVertical(
		lg.Top,
		header,
		separator,
		info))
}
