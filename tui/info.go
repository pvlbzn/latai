package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"math"
	"strings"
	"time"
)

// InfoComponent is an informational component which displays
// data about the latency measurement.
type InfoComponent struct {
	width int
	info  map[int]modelInfo
	model selectedModel
}

type modelInfo struct {
	rowID   int
	avg     string
	samples []time.Duration
}

func (s *modelInfo) getMaxLatency() time.Duration {
	if len(s.samples) == 0 {
		return 0
	}
	max := s.samples[0]
	for _, latency := range s.samples {
		if latency > max {
			max = latency
		}
	}
	return max
}

func (s *modelInfo) getMinLatency() time.Duration {
	if len(s.samples) == 0 {
		return 0
	}
	min := s.samples[0]
	for _, latency := range s.samples {
		if latency < min {
			min = latency
		}
	}
	return min
}

func (s *modelInfo) getJitter() time.Duration {
	if len(s.samples) == 0 {
		return 0
	}

	var sum time.Duration
	for _, latency := range s.samples {
		sum += latency
	}
	mean := sum / time.Duration(len(s.samples))

	var varianceSum float64
	for _, latency := range s.samples {
		diff := float64(latency - mean)
		varianceSum += diff * diff
	}
	variance := varianceSum / float64(len(s.samples))

	return time.Duration(math.Sqrt(variance))
}

// NewInfoComponent creates an instance of a new panel which displays
// data about models and latency details. Info panel shows details about
// models based on their row ID. Models are stored internally, and they are
// mapped onto their row IDs in the way it matches with table row IDs.
func NewInfoComponent(width int) *InfoComponent {
	return &InfoComponent{
		width: width,
		info:  make(map[int]modelInfo),
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
			Render(fmt.Sprintf(
				"Info: %s | %s | %s | %s",
				s.model.modelFamily, s.model.modelName, s.model.providerName, s.model.vendorName))
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

	var content string
	if info, ok := s.info[s.model.id]; !ok {
		content = rowStyle.
			Foreground(lg.Color("240")).
			Render("Press enter to run a measurement.")
	} else {
		data := fmt.Sprintf(
			"Runs: %d\tAvg: %s\tMin: %d\tMax: %d\tJitter: %d",
			len(info.samples), info.avg, info.getMinLatency().Milliseconds(), info.getMaxLatency().Milliseconds(), info.getJitter().Milliseconds())
		content = rowStyle.
			Foreground(lg.Color("231")).
			Render(fmt.Sprintf(data))
	}

	return container.Render(lg.JoinVertical(
		lg.Top,
		header,
		separator,
		content))
}

func (s *InfoComponent) AddInfo(rowID int, avg string, samples []time.Duration) {

	s.info[rowID] = modelInfo{
		rowID:   rowID,
		avg:     avg,
		samples: samples,
	}
}
