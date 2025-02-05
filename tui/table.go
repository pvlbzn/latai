package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/pvlbzn/latai/evaluator"
	"github.com/pvlbzn/latai/prompt"
	"github.com/pvlbzn/latai/provider"
	"math"
	"sort"
	"strconv"
)

type TableComponent struct {
	// Table data.
	table table.Model
	rows  []table.Row

	// Initialized providers and their models.
	providers []*tuiProvider

	// Sorting order.
	sortAsc bool

	logger *LoggerComponent
}

type tuiProvider struct {
	provider provider.Provider
	models   []*provider.Model
}

func NewTableComponent(providers []provider.Provider, logger *LoggerComponent) *TableComponent {
	var tuiProviders []*tuiProvider
	for _, p := range providers {
		models := p.GetLLMModels("")
		tuiProviders = append(tuiProviders, &tuiProvider{
			provider: p,
			models:   models,
		})
	}

	t, r := makeTableModel(tuiProviders)

	return &TableComponent{
		table:     t,
		rows:      r,
		providers: tuiProviders,
		logger:    logger,
	}
}

func makeTableModel(tuiProviders []*tuiProvider) (table.Model, []table.Row) {
	height := 30
	columns := []table.Column{
		{Title: "ID", Width: 2},
		{Title: "Name", Width: 32},
		{Title: "Provider", Width: 8},
		{Title: "Vendor", Width: 8},
		{Title: "Latency", Width: 7},
	}

	// Get sequential list of all models.
	var models []*provider.Model
	for _, p := range tuiProviders {
		models = append(models, p.models...)
	}

	// Create rows
	var rows []table.Row
	for i, m := range models {
		rows = append(rows, table.Row{strconv.Itoa(i), m.Name, string(m.Provider), string(m.Vendor), " "})
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

	return t, rows
}

func (s *TableComponent) getModelByRowID(rowID int) (provider.Provider, *provider.Model, error) {
	type pair struct {
		provider provider.Provider
		model    *provider.Model
	}

	var index []*pair
	for _, p := range s.providers {
		for _, m := range p.models {
			index = append(index, &pair{p.provider, m})
		}
	}

	if rowID > len(index) {
		return nil, nil, ErrIndexNotFound
	}

	return index[rowID].provider, index[rowID].model, nil
}

func (s *TableComponent) countAllModels() int {
	var sum int
	for _, p := range s.providers {
		sum += len(p.models)
	}

	return sum
}

func (s *TableComponent) Init() tea.Cmd {
	return nil
}

func (s *TableComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case latencyUpdatedMsg:
		// Log event.
		s.logger.Push(fmt.Sprintf(
			"%s latency %s ms", msg.name, msg.latency))

		// Update latency and whole table.
		s.rows[msg.id][4] = msg.latency
		s.table.SetRows(s.rows)

		return s, nil

	case latencyErrMsg:
		s.logger.Push(fmt.Sprintf("error measuring %s model: %s", msg.name, msg.err))
		s.rows[msg.id][4] = "err"
		s.table.SetRows(s.rows)
	}

	var cmd tea.Cmd

	// Pass other messages to the table
	s.table, cmd = s.table.Update(msg)
	return s, cmd
}

func (s *TableComponent) View() string {
	return lg.JoinVertical(
		lg.Top,
		s.makeTableView(),
	)
}

// makeTableView returns a view of table of models and help string.
func (s *TableComponent) makeTableView() string {
	return lg.NewStyle().
		BorderStyle(lg.NormalBorder()).
		BorderForeground(lg.Color("241")).
		Render(s.table.View() + "\n" + s.makeHelpView())
}

// makeHelpView returns a view of a single help string.
func (s *TableComponent) makeHelpView() string {

	return lg.NewStyle().
		Foreground(lg.Color("241")).
		PaddingTop(1).
		PaddingLeft(1).
		Render(fmt.Sprintf("enter: run | A: run all | J/K: up/down | s: sort | q: quit"))
}

func (s *TableComponent) ToggleFocus() {
	if s.table.Focused() {
		s.table.Blur()
	} else {
		s.table.Focus()
	}
}

func (s *TableComponent) SortByLatency() tea.Cmd {
	return sortRowsCmd(s)
}

func (s *TableComponent) ScrollTop() {
	s.table.SetCursor(0)
}

func (s *TableComponent) ScrollBottom() {
	s.table.SetCursor(len(s.rows) - 1)
}

func (s *TableComponent) MeasureRowLatency() tea.Cmd {
	selectedRowID, err := strconv.Atoi(s.table.SelectedRow()[0])
	if err != nil {
		s.logger.Push("error selecting row ID: " + err.Error())
	}

	// Log event.
	s.logger.Push(fmt.Sprintf("measuring %s latency", s.rows[selectedRowID][1]))

	// Update fields.
	s.rows[selectedRowID][4] = "..."
	s.table.SetRows(s.rows)

	// Start the concurrent task and return a command
	return fetchModelLatencyCmd(s, selectedRowID)
}

func (s *TableComponent) MeasureAllRowLatency() tea.Cmd {
	s.logger.Push(fmt.Sprintf("running %d parallel benchmarks", s.countAllModels()))

	for _, r := range s.rows {
		r[4] = "..."
	}
	s.table.SetRows(s.rows)

	return fetchAllModelLatencyCmd(s)
}

func (s *TableComponent) UpdateLatency(id int, latency string) {
	s.rows[id][4] = latency
	s.table.SetRows(s.rows)
}

func (s *TableComponent) SetLatencyError(id int) {
	s.rows[id][4] = "err"
	s.table.SetRows(s.rows)
}

func fetchModelLatencyCmd(t *TableComponent, modelRowID int) tea.Cmd {
	return func() tea.Msg {
		// Process the selected row (e.g., calculate latency or fetch new data)
		p, m, err := t.getModelByRowID(modelRowID)
		if err != nil {
			return latencyErrMsg{modelRowID, m.Name, err.Error()}
		}

		prompts, err := prompt.GetPrompts()
		if err != nil {
			return latencyErrMsg{modelRowID, m.Name, err.Error()}
		}

		t.logger.Push(fmt.Sprintf("sampling with %d default prompts", len(prompts)))

		eval := evaluator.NewEvaluator(p, m, prompts...)
		res, err := eval.Evaluate()
		if err != nil {
			return latencyErrMsg{modelRowID, m.Name, err.Error()}
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
func fetchAllModelLatencyCmd(t *TableComponent) tea.Cmd {
	counter := 0
	cmds := make([]tea.Cmd, t.countAllModels())
	for _, p := range t.providers {
		for range p.models {
			cmds[counter] = fetchModelLatencyCmd(t, counter)
			counter++
		}
	}

	return tea.Batch(cmds...)
}

type sortRowsMsg struct{}

func sortRowsCmd(s *TableComponent) tea.Cmd {
	return func() tea.Msg {
		sort.SliceStable(s.rows, func(i, j int) bool {
			latencyI, errI := strconv.Atoi(s.rows[i][4])
			latencyJ, errJ := strconv.Atoi(s.rows[j][4])

			if s.sortAsc {
				if errI != nil {
					latencyI = math.MaxInt
				}
				if errJ != nil {
					latencyJ = math.MaxInt
				}

				return latencyI < latencyJ
			}

			if errI != nil {
				latencyI = 0
			}
			if errJ != nil {
				latencyJ = 0
			}

			return latencyI > latencyJ
		})

		// Toggle sorting order.
		s.sortAsc = !s.sortAsc

		s.table.SetRows(s.rows)

		return sortRowsMsg{}
	}
}
