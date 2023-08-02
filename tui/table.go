package tui

import (
	"fmt"
	"strings"

	"github.com/CarlosGMI/Playlistify/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	textTable "github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/term"
)

type tableKeymap struct {
	newSearch           key.Binding
	newSearchInPlaylist key.Binding
	switchMode          key.Binding
}
type tableHelpOption struct {
	character   string
	description string
	condition   bool
}
type tableContext struct {
	selectedPlaylist string
}
type TableModel struct {
	table       table.Model
	tableType   string
	updatable   bool
	previewText string
	showHelp    bool
	mode        string
	viewport    viewport.Model
	context     tableContext
}
type SelectedItemMsg struct {
	Item string
}

const defaultTableHeight = 15

var tableKeys = tableKeymap{
	newSearch: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "New search"),
	),
	newSearchInPlaylist: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "New search in current playlist"),
	),
	switchMode: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Switch table mode"),
	),
}
var tableHelpOptions = []tableHelpOption{
	{"n", "new search", true},
	{"p", "new search in current playlist", true},
	{"s", "switch to text view", true},
	{"s", "switch to table view", false},
	{"q", "quit", true},
	{"esc", "quit", true},
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

func CreateTable(
	tableType string,
	rows []table.Row,
	textRows []textTable.Row,
	updatable bool,
	previewText string,
	context tableContext,
) TableModel {
	var tableHeight = defaultTableHeight
	isSearchTable := tableType == TableTypes[1]
	isSongsTable := tableType == TableTypes[0]
	showHelp := isSearchTable || (isSongsTable && !updatable)
	tableHelpOptions[0].condition = isSearchTable
	tableHelpOptions[1].condition = isSearchTable
	terminalWidth, _, err := term.GetSize(0)

	if err != nil {
		terminalWidth = 122
	}

	if len(rows) < tableHeight {
		tableHeight = len(rows)
	}

	textTableViewport := viewport.New(terminalWidth, tableHeight+4)

	newTable := table.New(
		table.WithColumns(columns[tableType]),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)
	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)

	styles.Selected.Foreground(lipgloss.Color("#FFFFFF"))
	styles.Selected.Background(lipgloss.Color(utils.ColorSpotifyGreenOpaque))
	newTable.SetStyles(styles)
	textTableViewport.SetContent(textView(tableType, textRows))

	return TableModel{newTable,
		tableType,
		updatable,
		previewText,
		showHelp,
		"table",
		textTableViewport,
		context,
	}
}

func (model TableModel) Init() tea.Cmd {
	return nil
}

func (model TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return model, tea.Quit
		case "enter":
			if model.updatable {
				searchModel := CreateSearchModel(true, model.table.SelectedRow()[0], "")
				msg := SelectedItemMsg{model.table.SelectedRow()[0]}

				return searchModel.Update(msg)
			}
		case "up":
			if model.table.Cursor() == 0 {
				model.table.GotoBottom()

				return model, nil
			}
		case "down":
			if model.table.Cursor() == len(model.table.Rows())-1 {
				model.table.GotoTop()

				return model, nil
			}
		}
		switch {
		case key.Matches(msg, tableKeys.newSearch):
			if model.showHelp {
				searchModel := CreateSearchModel(true, "", "")

				return searchModel, searchModel.Init()
			}
		case key.Matches(msg, tableKeys.newSearchInPlaylist):
			if model.showHelp {
				searchModel := CreateSearchModel(false, model.context.selectedPlaylist, "")

				return searchModel, searchModel.Init()
			}
		case key.Matches(msg, tableKeys.switchMode):
			if !model.updatable {
				if model.mode == utils.TableModeDefault {
					model.mode = utils.TableModeText
					tableHelpOptions[1].condition = false
					tableHelpOptions[2].condition = true
				} else {
					model.mode = utils.TableModeDefault
					tableHelpOptions[1].condition = true
					tableHelpOptions[2].condition = false
				}
			}
		}
	}

	if model.mode == utils.TableModeDefault {
		model.table, cmd = model.table.Update(msg)
	} else {
		model.viewport, cmd = model.viewport.Update(msg)
	}

	cmds = append(cmds, cmd)

	return model, tea.Batch(cmds...)
}

func (model TableModel) View() string {
	var content string

	if len(model.previewText) > 0 {
		content += fmt.Sprintf("%s\n\n", model.previewText)
	}

	if model.mode == utils.TableModeDefault {
		content += fmt.Sprintf("%s\n\n", tableBaseStyle.Render(model.table.View()))
	} else {
		content += fmt.Sprintf("%s\n\n", model.viewport.View())
	}

	if model.showHelp {
		content += model.helpView()
	}

	return fmt.Sprintf("\n%s\n\n", content)
}

func (model TableModel) helpView() string {
	var optionsToShow []string

	for i := range tableHelpOptions {
		if tableHelpOptions[i].condition {
			optionsToShow = append(optionsToShow, fmt.Sprintf("%s: %s", tableHelpOptions[i].character, tableHelpOptions[i].description))
		}
	}

	return utils.HelpStyle(" " + strings.Join(optionsToShow, " • "))
}

func textView(tableType string, rows []textTable.Row) string {
	var header textTable.Row
	newTable := textTable.NewWriter()
	currentColumns := columns[tableType]

	for i := range currentColumns {
		header = append(header, currentColumns[i].Title)
	}

	newTable.AppendHeader(header)
	newTable.AppendRows(rows)

	return newTable.Render()
}
