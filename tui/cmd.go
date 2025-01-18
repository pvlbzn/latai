package tui

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pvlbzn/genlat/evaluator"
	"github.com/pvlbzn/genlat/prompt"
)

type latencyUpdatedMsg struct {
	id   int
	data table.Row
}

func fetchModelLatencyCmd(m *TableModel, modelRowID int) tea.Cmd {
	return func() tea.Msg {
		// Process the selected row (e.g., calculate latency or fetch new data)
		model := m.models[modelRowID]
		prompts, err := prompt.GetPrompts()
		if err != nil {
			panic(err)
		}

		eval := evaluator.NewEvaluator(m.provider, model, prompts...)
		res, err := eval.Evaluate()
		if err != nil {
			panic(err)
		}

		// Return an updateRowMsg to update the table row
		return latencyUpdatedMsg{
			id: modelRowID,
			data: table.Row{
				m.rows[modelRowID][0],                         // ID
				m.rows[modelRowID][1],                         // Name
				m.rows[modelRowID][2],                         // Provider
				m.rows[modelRowID][3],                         // Vendor
				fmt.Sprintf("%d", res.Latency.Milliseconds()), // Updated Latency
			},
		}
	}
}

type sortRowsMsg struct{}

func sortRowsCmd(m *TableModel) tea.Cmd {
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
