package tui

import (
	"errors"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/pvlbzn/latai/provider"
)

var (
	ErrIndexNotFound = errors.New("index not found")
)

// TUIModel is a root of Latai TUI application. It holds data and state
// for the whole application.
type TUIModel struct {
	tableComponent  *TableComponent
	infoComponent   *InfoComponent
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

	l := NewLoggerComponent(70)
	t := NewTableComponent(providers, l)
	i := NewInfoComponent(70)

	return &TUIModel{
		tableComponent:  t,
		infoComponent:   i,
		loggerComponent: l,
	}, nil
}

func (m *TUIModel) Init() tea.Cmd {
	return nil
}

// Update returns a new model and a command. Commands are functions
// which designed to perform side effects.
func (m *TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.tableComponent.ToggleFocus()
			return m, nil

		case "q", "ctrl+c":
			// Quit.
			return m, tea.Quit

		case "s":
			// Sort by latency.
			m.loggerComponent.Push(fmt.Sprintf("sorting, order asc: %v", m.tableComponent.sortAsc))
			return m, m.tableComponent.SortByLatency()

		case "J":
			// Scroll all the way up.
			m.tableComponent.ScrollTop()
			cmds = append(cmds, m.notifySelection())

		case "K":
			// Scroll all the way down.
			m.tableComponent.ScrollBottom()
			cmds = append(cmds, m.notifySelection())

		case "j", "down":
			// One row down.
			if m.tableComponent.MoveCursorDown() {
				cmds = append(cmds, m.notifySelection())
			}

		case "k", "up":
			// One row up.
			if m.tableComponent.MoveCursorUp() {
				cmds = append(cmds, m.notifySelection())
			}

		case "enter":
			// Run latency measurement for a selected model.
			return m, m.tableComponent.MeasureRowLatency()

		case "A":
			// Run latency measurement for all models.
			return m, m.tableComponent.MeasureAllRowLatency()
		}

	case latencyUpdatedMsg:
		m.loggerComponent.Push(fmt.Sprintf("%s latency %s ms", msg.name, msg.latency))
		m.tableComponent.UpdateLatency(msg.id, msg.latency)
		m.infoComponent.AddInfo(msg.id, msg.latency, msg.samples)
		return m, nil

	case latencyErrMsg:
		m.loggerComponent.Push(fmt.Sprintf("error measuring %s model: %s", msg.name, msg.err))
		m.tableComponent.SetLatencyError(msg.id)
		return m, nil

	case modelSelectedMsg:
		m.infoComponent.Update(msg)
		return m, nil
	}

	// Pass other messages.
	var tableCmd tea.Cmd
	m.tableComponent.table, tableCmd = m.tableComponent.table.Update(msg)
	cmds = append(cmds, tableCmd)

	return m, tea.Batch(cmds...)
}

type selectedModel struct {
	id           int
	providerName provider.ModelProvider
	vendorName   provider.ModelVendor
	modelFamily  provider.ModelFamily
	modelName    string
}

type modelSelectedMsg struct {
	selectedModel
}

// Message components on table row selection change.
func (m *TUIModel) notifySelection() tea.Cmd {
	// ✅ Delay selection processing slightly to ensure correct row is fetched
	return tea.Tick(time.Millisecond*10, func(time.Time) tea.Msg {
		id, model, err := m.tableComponent.GetSelectedRow()
		if err != nil {
			m.loggerComponent.Push("error while selecting row: " + err.Error())
		}

		return modelSelectedMsg{
			selectedModel{
				id:           id,
				providerName: model.Provider,
				modelName:    model.Name,
				vendorName:   model.Vendor,
				modelFamily:  model.Family,
			},
		}
	})
}

func (m *TUIModel) View() string {
	return lg.JoinVertical(
		lg.Top,
		m.tableComponent.View(),
		m.infoComponent.View(),
		m.loggerComponent.View(),
	)
}

type latencyUpdatedMsg struct {
	id      int
	name    string
	latency string
	samples []time.Duration
}

type latencyErrMsg struct {
	id   int
	name string
	err  string
}
