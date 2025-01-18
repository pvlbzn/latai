package tui

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
	"github.com/pvlbzn/genlat/provider"
)

type TableModel struct {
	// Dimensions.
	width int

	// Table data.
	table table.Model
	rows  []table.Row

	// Initialized services.
	openaiProvider *provider.OpenAI
	models         []*provider.Model

	// Sorting order.
	sortAsc bool

	// Log events.
	loggerComponent *LoggerComponent
}

func NewTableModel() TableModel {
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

	width := 72
	height := len(models) + 1
	columns := []table.Column{
		{Title: "ID", Width: 2},
		{Title: "Name", Width: 32},
		{Title: "Provider", Width: 8},
		{Title: "Vendor", Width: 8},
		{Title: "Latency (ms)", Width: 12},
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

	return TableModel{
		table:           t,
		rows:            rows,
		width:           width,
		openaiProvider:  p,
		models:          models,
		sortAsc:         true,
		loggerComponent: NewLoggerComponent(width),
	}
}

func (m *TableModel) Init() tea.Cmd {
	return nil
}

// Update returns a new model and a command. Commands are functions
// which designed to perform side effects.
func (m *TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "s":
			m.loggerComponent.Push("Sorting!")
			return m, sortRowsCmd(m)
		case "enter":
			// Get the selected row index
			selectedRowID, err := strconv.Atoi(m.table.SelectedRow()[0])
			if err != nil {
				panic(err)
			}

			m.rows[selectedRowID][4] = "..."
			m.table.SetRows(m.rows)

			// Start the concurrent task and return a command
			return m, fetchModelLatencyCmd(m, selectedRowID)
		}
	case updateRowMsg:
		// Update the latency column in the selected row
		m.rows[msg.id][4] = msg.data[4] // Assuming latency is in the 5th column

		// Update the table with the modified rows
		m.table.SetRows(m.rows)
		return m, nil
	}

	// Pass other messages to the table
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *TableModel) View() string {
	tableView := m.makeTableView()

	// Stack vertically views.
	return lg.JoinVertical(
		lg.Top,
		tableView,
		m.loggerComponent.View(),
	)
}

// makeTableView returns a view of table of models and help string.
func (m *TableModel) makeTableView() string {
	return lg.NewStyle().
		BorderStyle(lg.NormalBorder()).
		BorderForeground(lg.Color("241")).
		Render(m.table.View() + "\n" + m.makeHelpView())
}

// makeHelpView returns a view of a single help string.
func (m *TableModel) makeHelpView() string {
	return lg.NewStyle().
		Foreground(lg.Color("241")).
		Render(fmt.Sprintf(" s: sort by latency | q: quit"))
}

func Run() {
	m := NewTableModel()

	if _, err := tea.NewProgram(&m).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
