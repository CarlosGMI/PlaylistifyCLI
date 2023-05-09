package tui

import (
	"fmt"

	"github.com/CarlosGMI/Playlistify/utils"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TableModel struct {
	table       table.Model
	tableType   string
	updatable   bool
	previewText string
}
type SelectedItemMsg struct {
	Item string
}

var TableTypes = []string{utils.PlaylistsTable, utils.SongsTable}
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
		{Title: "ARTISTS", Width: 50},
	},
}

func CreateTable(tableType string, rows []table.Row, updatable bool, previewText string) TableModel {
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
		Bold(true)

	if !updatable {
		styles.Selected.Bold(false)
		styles.Selected.UnsetForeground()
	} else {
		styles.Selected.Foreground(lipgloss.Color("#FFFFFF"))
		styles.Selected.Background(lipgloss.Color(utils.ColorSpotifyGreenOpaque))
	}

	newTable.SetStyles(styles)

	return TableModel{newTable, tableType, updatable, previewText}
}

func (model TableModel) Init() tea.Cmd {
	return nil
}

func (model TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if !model.updatable {
		return model, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return model, tea.Quit
		case "enter":
			searchModel := CreateSearchModel(true, model.table.SelectedRow()[0], "")
			msg := SelectedItemMsg{model.table.SelectedRow()[0]}

			return searchModel.Update(msg)
		}
	}

	model.table, cmd = model.table.Update(msg)
	cmds = append(cmds, cmd)

	return model, tea.Batch(cmds...)
}

func (model TableModel) View() string {
	if model.updatable && model.tableType == utils.PlaylistsTable {
		return fmt.Sprintf("\n%s\n\n%s\n\n", model.previewText, tableBaseStyle.Render(model.table.View()))
	}

	return fmt.Sprintf("\n%s\n\n", tableBaseStyle.Render(model.table.View()))
}
