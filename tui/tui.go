package tui

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/pvlbzn/latai/provider"
)

var (
	ErrIndexNotFound = errors.New("index not found")
)

// TUIModel is a root of Genlat TUI application. It holds data and state
// for the whole application.
type TUIModel struct {
	tableComponent  *TableComponent
	loggerComponent *LoggerComponent

	width  int
	height int
}

func NewTUIModel() (*TUIModel, error) {
	// Initialize providers.
	openai, err := provider.NewOpenAI("")
	if err != nil {
		return nil, err
	}

	bedrock, err := provider.NewBedrock(provider.DefaultAWSRegion, provider.DefaultAWSProfile)
	if err != nil {
		return nil, err
	}

	groq, err := provider.NewGroq("")
	if err != nil {
		return nil, err
	}

	providers := []provider.Provider{openai, bedrock, groq}

	l := NewLoggerComponent(67)
	t := NewTableComponent(providers, l)

	return &TUIModel{
		tableComponent:  t,
		loggerComponent: l,
	}, nil
}

func (m *TUIModel) Init() tea.Cmd {
	return nil
}

// Update returns a new model and a command. Commands are functions
// which designed to perform side effects.
func (m *TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Key press handlers.
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.tableComponent.ToggleFocus()
			return m, nil

		case "q", "ctrl+c": // Quit.
			return m, tea.Quit

		case "s": // Sort by latency.
			m.loggerComponent.Push(fmt.Sprintf("sorting, order asc: %v", m.tableComponent.sortAsc))
			return m, m.tableComponent.SortByLatency()

		case "J": // Scroll all the way up.
			m.tableComponent.ScrollTop()
			return m, nil

		case "K": // Scroll all the way down.
			m.tableComponent.ScrollBottom()
			return m, nil

		case "enter": // Run latency measurement for a selected model.
			return m, m.tableComponent.MeasureRowLatency()

		case "A": // Run latency measurement for all models.
			return m, m.tableComponent.MeasureAllRowLatency()
		}

	case latencyUpdatedMsg:
		m.loggerComponent.Push(fmt.Sprintf("%s latency %s ms", msg.name, msg.latency))
		m.tableComponent.UpdateLatency(msg.id, msg.latency)
		return m, nil

	case latencyErrMsg:
		m.loggerComponent.Push(fmt.Sprintf("error measuring %s model: %s", msg.name, msg.err))
		m.tableComponent.SetLatencyError(msg.id)
		return m, nil
	}

	// Pass other messages to the table component.
	var cmd tea.Cmd
	m.tableComponent.table, cmd = m.tableComponent.table.Update(msg)

	return m, cmd
}

func (m *TUIModel) View() string {
	return lg.JoinVertical(
		lg.Top,
		m.tableComponent.View(),
		m.loggerComponent.View(),
	)
}

type latencyUpdatedMsg struct {
	id      int
	name    string
	latency string
}

type latencyErrMsg struct {
	id   int
	name string
	err  string
}
