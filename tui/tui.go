package tui

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/pvlbzn/genlat/evaluator"
	"github.com/pvlbzn/genlat/prompt"
	"github.com/pvlbzn/genlat/provider"
)

// TUIModel is a root of Genlat TUI application. It holds data and state
// for the whole application.
type TUIModel struct {
	// Dimensions.
	width int

	// Table data.
	table table.Model
	rows  []table.Row

	// Initialized services.
	provider provider.Provider
	models   []*provider.Model

	// Sorting order.
	sortAsc bool

	// Log events.
	loggerComponent *LoggerComponent
}

func NewTUIModel() TUIModel {
	// Initialize providers.
	p, err := provider.NewOpenAI("")
	if err != nil {
		panic(err)
	}

	// Fetch LLM models to list.
	models, err := p.GetLLMModels("")
	if err != nil {
		panic(err)
	}

	return makeTableModel(p, models)
}

func makeTableModel(provider provider.Provider, models []*provider.Model) TUIModel {
	// TODO: make width dynamic to whatever terminal sized
	width := 67
	height := len(models) + 1
	columns := []table.Column{
		{Title: "ID", Width: 2},
		{Title: "Name", Width: 32},
		{Title: "Provider", Width: 8},
		{Title: "Vendor", Width: 8},
		{Title: "Latency", Width: 7},
	}
	var rows []table.Row
	for i, m := range models {
		rows = append(rows, table.Row{strconv.Itoa(i), m.Name, m.Provider, m.Vendor, " "})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(height))

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lg.NormalBorder()).
		BorderForeground(lg.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lg.Color("#fff")).
		Background(lg.Color("#2b5ccc")).
		Bold(true)

	t.SetStyles(s)

	logger := NewLoggerComponent(width)
	logger.Push(fmt.Sprintf("loaded %d models", len(models)))

	return TUIModel{
		table:           t,
		rows:            rows,
		width:           width,
		provider:        provider,
		models:          models,
		sortAsc:         true,
		loggerComponent: logger,
	}
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
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}

		case "q", "ctrl+c": // Quit.
			return m, tea.Quit

		case "s": // Sort by latency.
			m.loggerComponent.Push(fmt.Sprintf("sorting, order asc: %v", m.sortAsc))
			return m, sortRowsCmd(m)

		case "J": // Scroll all the way up.
			m.table.SetCursor(0)
			return m, nil

		case "K": // Scroll all the way down.
			m.table.SetCursor(len(m.rows) - 1)
			return m, nil

		case "enter": // Run latency measurement for a selected model.
			// Get the selected row index
			selectedRowID, err := strconv.Atoi(m.table.SelectedRow()[0])
			if err != nil {
				panic(err)
			}

			// Log event.
			m.loggerComponent.Push(fmt.Sprintf(
				"measuring %s latency", m.rows[selectedRowID][1]))

			// Update fields.
			m.rows[selectedRowID][4] = "..."
			m.table.SetRows(m.rows)

			// Start the concurrent task and return a command
			return m, fetchModelLatencyCmd(m, selectedRowID)

		case "A": // Run latency measurement for all models.
			m.loggerComponent.Push(fmt.Sprintf("running %d parallel benchmarks", len(m.models)))

			for _, r := range m.rows {
				r[4] = "..."
			}
			m.table.SetRows(m.rows)

			return m, fetchAllModelLatencyCmd(m)
		}

	case latencyUpdatedMsg:
		// Log event.
		m.loggerComponent.Push(fmt.Sprintf(
			"%s latency %s ms", msg.name, msg.latency))

		// Update latency and whole table.
		m.rows[msg.id][4] = msg.latency
		m.table.SetRows(m.rows)

		return m, nil
	}

	var cmd tea.Cmd

	// Pass other messages to the table
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *TUIModel) View() string {
	tableView := m.makeTableView()

	// Stack vertically views.
	return lg.JoinVertical(
		lg.Top,
		tableView,
		m.loggerComponent.View(),
	)
}

// makeTableView returns a view of table of models and help string.
func (m *TUIModel) makeTableView() string {
	return lg.NewStyle().
		BorderStyle(lg.NormalBorder()).
		BorderForeground(lg.Color("241")).
		Render(m.table.View() + "\n" + m.makeHelpView())
}

// makeHelpView returns a view of a single help string.
func (m *TUIModel) makeHelpView() string {

	return lg.NewStyle().
		Foreground(lg.Color("241")).
		PaddingTop(1).
		PaddingLeft(1).
		Render(fmt.Sprintf("enter: run | A: run all | J/K: up/down | s: sort | q: quit"))
}

type latencyUpdatedMsg struct {
	id      int
	name    string
	latency string
}

func fetchModelLatencyCmd(m *TUIModel, modelRowID int) tea.Cmd {
	return func() tea.Msg {
		// Process the selected row (e.g., calculate latency or fetch new data)
		model := m.models[modelRowID]
		prompts, err := prompt.GetPrompts()
		if err != nil {
			panic(err)
		}

		m.loggerComponent.Push(fmt.Sprintf("sampling with %d default prompts", len(prompts)))

		eval := evaluator.NewEvaluator(m.provider, model, prompts...)
		res, err := eval.Evaluate()
		if err != nil {
			panic(err)
		}

		// Return an updateRowMsg to update the table row
		return latencyUpdatedMsg{
			id:      modelRowID,
			name:    res.ModelName,
			latency: fmt.Sprintf("%d", res.Latency.Milliseconds()),
		}
	}
}

// Create a batch of commands to fetch all latency, one model at a time.
func fetchAllModelLatencyCmd(m *TUIModel) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.models))
	for i := range m.models {
		cmds[i] = fetchModelLatencyCmd(m, i)
	}

	return tea.Batch(cmds...)
}

type sortRowsMsg struct{}

func sortRowsCmd(m *TUIModel) tea.Cmd {
	return func() tea.Msg {
		sort.SliceStable(m.rows, func(i, j int) bool {
			latencyI, errI := strconv.Atoi(m.rows[i][4])
			latencyJ, errJ := strconv.Atoi(m.rows[j][4])

			if errI != nil {
				latencyI = 0
			}
			if errJ != nil {
				latencyJ = 0
			}

			if m.sortAsc {
				return latencyI < latencyJ
			}

			return latencyI > latencyJ
		})

		// Toggle sorting order.
		m.sortAsc = !m.sortAsc

		m.table.SetRows(m.rows)

		return sortRowsMsg{}
	}
}
