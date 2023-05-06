package tui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TableModel struct {
	table     table.Model
	updatable bool
}

type UpdateFunction func()

var TableTypes = []string{"PLAYLISTS", "SONGS"}
var tableBaseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
var columns = map[string][]table.Column{
	TableTypes[0]: {
		{Title: "PLAYLIST ID", Width: 12},
		{Title: "PLAYLIST NAME", Width: 50},
		{Title: "TOTAL TRACKS", Width: 20},
	},
	TableTypes[1]: {
		{Title: "#", Width: 12},
		{Title: "NAME", Width: 50},
		{Title: "ARTISTS", Width: 100},
	},
}

func CreateTable(tableType string, rows []table.Row, updatable bool) TableModel {
	newTable := table.New(
		table.WithColumns(columns[tableType]),
		table.WithRows(rows),
		table.WithFocused(updatable),
		table.WithHeight(len(rows)),
	)
	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	if !updatable {
		styles.Selected.Bold(false)
		styles.Selected.UnsetForeground()
	}

	newTable.SetStyles(styles)

	return TableModel{newTable, updatable}
}

func (model TableModel) Init() tea.Cmd {
	return nil
}

func (model TableModel) View() string {
	result := tableBaseStyle.Render(model.table.View()) + "\n"

	return result
}

func (model TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if !model.updatable {
		return model, tea.Quit
	}

	return model, cmd
}
