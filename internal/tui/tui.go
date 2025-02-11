package tui

import (
	"errors"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/pvlbzn/latai/internal/provider"
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
	var providers []provider.Provider
	l := NewLoggerComponent(70)

	// Initialize providers.
	openai, err := initializeProvider(
		l,
		provider.ModelProviderOpenAI,
		func() (provider.Provider, error) {
			return provider.NewOpenAI("")
		})
	if err == nil {
		providers = append(providers, openai)
	}

	bedrock, err := initializeProvider(
		l,
		provider.ModelProviderBedrock,
		func() (provider.Provider, error) {
			return provider.NewBedrock("", "")
		})
	if err == nil {
		providers = append(providers, bedrock)
	}

	groq, err := initializeProvider(
		l,
		provider.ModelProviderGroq,
		func() (provider.Provider, error) {
			return provider.NewGroq("")
		})
	if err == nil {
		providers = append(providers, groq)
	}

	t := NewTableComponent(providers, l)
	i := NewInfoComponent(70)

	return &TUIModel{
		tableComponent:  t,
		infoComponent:   i,
		loggerComponent: l,
	}, nil
}

func initializeProvider(l *LoggerComponent, name provider.ModelProvider, newProvider func() (provider.Provider, error)) (provider.Provider, error) {
	errProvider := errors.New("provider initialization failed")

	p, err := newProvider()
	if err != nil {
		if name == provider.ModelProviderBedrock {
			// Bedrock uses a different initialization logic, make message informative
			// for the use case.
			l.Push(fmt.Sprintf(
				"Bedrock not loaded, verify your `AWS_PROFILE` and `AWS_REGION`."))
		} else {
			l.Push(fmt.Sprintf(
				"%s not loaded. API key not found, `%s_API_KEY` envar is required.",
				name, strings.ToUpper(string(name))))
		}

		return nil, errProvider
	}

	if ok := p.VerifyAccess(); !ok {
		l.Push(fmt.Sprintf(
			"%s provider is not loaded. API key is invalid, verify your `%s_API_KEY`.",
			p.Name(), strings.ToUpper(string(p.Name()))))
		return nil, errProvider
	}

	l.Push(fmt.Sprintf("%s provider is loaded.", p.Name()))
	return p, nil
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
			m.loggerComponent.Push(fmt.Sprintf("Sorting, order asc: %v", m.tableComponent.sortAsc))
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
		m.loggerComponent.Push(fmt.Sprintf("Error measuring %s model: %s", msg.name, msg.err))
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
			m.loggerComponent.Push("Error while selecting row: " + err.Error())
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
