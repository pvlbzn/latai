package tui

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pvlbzn/genlat/evaluator"
	"github.com/pvlbzn/genlat/prompt"
	"github.com/pvlbzn/genlat/provider"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#FFFDF5"))

type TableModel struct {
	table table.Model
	rows  []table.Row

	openaiProvider *provider.OpenAI
	models         []*provider.Model

	sortAsc bool
}

type updateRowMsg struct {
	id   int
	data table.Row
}

func NewTableModel() TableModel {
	p, err := provider.NewOpenAI("")
	if err != nil {
		panic(err)
	}

	models, err := p.GetLLMModels("")
	if err != nil {
		panic(err)
	}

	columns := []table.Column{
		{Title: "ID", Width: 2},
		{Title: "Name", Width: 32},
		{Title: "Provider", Width: 8},
		{Title: "Vendor", Width: 8},
		{Title: "Latency (ms)", Width: 12},
	}

	var rows []table.Row
	for i, m := range models {
		rows = append(rows, table.Row{strconv.Itoa(i), m.Name, m.Provider, m.Vendor, "?"})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(32))

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true)

	t.SetStyles(s)

	return TableModel{
		table:          t,
		rows:           rows,
		openaiProvider: p,
		models:         models,
		sortAsc:        true,
	}
}

func (m TableModel) Init() tea.Cmd {
	return nil
}

func (m TableModel) View() string {
	helpText := " 's' sort by latency | `q` quit"

	return baseStyle.Render(
		m.table.View() + "\n" +
			helpText)
}

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			sort.SliceStable(m.rows, func(i, j int) bool {
				latencyI, errI := strconv.Atoi(m.rows[i][4])
				latencyJ, errJ := strconv.Atoi(m.rows[j][4])

				// Handle "..." or invalid values
				if errI != nil {
					latencyI = 0 // Assign a large number to represent "..."
				}
				if errJ != nil {
					latencyJ = 0 // Assign a large number to represent "..."
				}

				if m.sortAsc {
					return latencyI < latencyJ
				}

				return latencyI > latencyJ
			})

			// Toggle sorting order.
			m.sortAsc = !m.sortAsc

			m.table.SetRows(m.rows)
			return m, nil
		case "enter":
			// Get the selected row index
			selectedRowID, err := strconv.Atoi(m.table.SelectedRow()[0])
			if err != nil {
				panic(err)
			}

			m.rows[selectedRowID][4] = "..."
			m.table.SetRows(m.rows)

			// Start the concurrent task and return a command
			return m, fetchLatencyCmd(m, selectedRowID)
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

func fetchLatencyCmd(m TableModel, rowID int) tea.Cmd {
	return func() tea.Msg {
		// Process the selected row (e.g., calculate latency or fetch new data)
		model := m.models[rowID]
		prompts, err := prompt.GetPrompts()
		if err != nil {
			panic(err)
		}

		eval := evaluator.NewEvaluator(m.openaiProvider, model, prompts...)
		res, err := eval.Evaluate()
		if err != nil {
			panic(err)
		}

		// Return an updateRowMsg to update the table row
		return updateRowMsg{
			id: rowID,
			data: table.Row{
				m.rows[rowID][0], // ID
				m.rows[rowID][1], // Name
				m.rows[rowID][2], // Provider
				m.rows[rowID][3], // Vendor
				fmt.Sprintf("%d", res.Latency.Milliseconds()), // Updated Latency
			},
		}
	}
}

func Run() {
	m := NewTableModel()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
