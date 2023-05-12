package tui

import (
	"fmt"
	"strings"

	"github.com/CarlosGMI/Playlistify/services"
	"github.com/CarlosGMI/Playlistify/utils"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SearchModel struct {
	state            string
	showPlaylists    bool
	loader           spinner.Model
	loaderText       string
	playlists        TableModel
	selectedPlaylist string
	searchInput      textinput.Model
	searchTerm       string
	searchTermError  string
	results          TableModel
	resultsText      string
}

func CreateSearchModel(showPlaylists bool, playlistId string, searchTerm string) SearchModel {
	model := SearchModel{
		state:            "",
		showPlaylists:    showPlaylists,
		loader:           CreateSpinner(),
		loaderText:       "Refreshing token...",
		selectedPlaylist: playlistId,
		searchTerm:       searchTerm,
	}

	return model
}

func (model SearchModel) Init() tea.Cmd {
	return services.InitAuthentication
}

func (model SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			if msg.String() == "q" && model.searchInput.Focused() {
				model.searchInput, cmd = model.searchInput.Update(msg)

				return model, cmd
			}

			return model, tea.Quit
		case "enter":
			model.searchTerm = model.searchInput.Value()

			if !isSearchTermValid(model.searchTerm) {
				model.searchInput, cmd = model.searchInput.Update(msg)
				model.searchTermError = "the search term must be longer than 3 characters (excluding white spaces)"

				return model, cmd
			} else {
				model.searchTermError = ""
				model.state = utils.LoadingState
				model.loaderText = utils.SearchingText
				model.loader, cmd = model.loader.Update(spinner.TickMsg{})
				cmds = append(cmds, executeSearch(model.selectedPlaylist, model.searchTerm), cmd)
			}
		default:
			if model.searchInput.Focused() {
				model.searchInput, cmd = model.searchInput.Update(msg)

				return model, cmd
			}

			return model, tea.Quit
		}
	case services.NotAuthenticatedMsg:
		if msg.ErrorType == utils.NotLoggedInCode {
			model.state = utils.ErrorState
			model.resultsText = msg.Message

			return model.Update(tea.KeyMsg{})
		} else if msg.ErrorType == utils.ExpiredTokenCode {
			model.state = utils.LoadingState
			model.loader, cmd = model.loader.Update(spinner.TickMsg{})
			cmds = append(cmds, refreshAuth(), cmd)
		}
	case services.LoggedInMsg, services.AuthErrorMsg:
		model.state = utils.LoadingState
		model.loader, cmd = model.loader.Update(spinner.TickMsg{})

		if model.showPlaylists {
			model.loaderText = "Fetching playlists..."
			cmds = append(cmds, fetchPlaylists(), cmd)
		} else {
			model.loaderText = utils.SearchingText
			cmds = append(cmds, executeSearch(model.selectedPlaylist, model.searchTerm), cmd)
		}
	case services.PlaylistsMsg:
		model.state = "table"
		playlists, textPlaylists, err := services.PrintPlaylists()

		if err != nil {
			model.state = utils.ErrorState
			model.resultsText = err.Error()

			return model, tea.Quit
		}

		model.playlists = CreateTable(utils.PlaylistsTable, playlists, textPlaylists, true, "Select the playlist to search:")

		return model.playlists.Update(msg)
	case SelectedItemMsg:
		model.selectedPlaylist = msg.Item
		model.state = utils.InputState
		model.searchInput = createSearchInput()
		model.searchInput, cmd = model.searchInput.Update(msg)
		cmds = append(cmds, cmd)
	case services.SearchResultsMsg:
		model.state = "table"
		model.results = CreateTable(utils.SongsTable, msg.Results, msg.TextResults, false, "")

		return model.results.Update(msg)
	case services.PlaylistsErrorMsg:
		model.state = utils.ErrorState
		model.resultsText = msg.Message

		return model, tea.Quit
	case spinner.TickMsg:
		model.loader, cmd = model.loader.Update(msg)

		return model, cmd
	default:
		return model, nil
	}

	return model, tea.Batch(cmds...)
}

func (model SearchModel) View() string {
	if model.state == utils.LoadingState {
		return fmt.Sprintf("\n %s %s\n\n", model.loader.View(), model.loaderText)
	} else if model.state == utils.ErrorState {
		return fmt.Sprintf("\n %s%s\n\n", utils.ErrorStyle("Error: "), model.resultsText)
	} else if model.state == utils.InputState {
		commonText := fmt.Sprintf("\n %s \n\n%s\n\n", "Enter your search term:", model.searchInput.View())

		if len(model.searchTermError) > 0 {
			return commonText + fmt.Sprintf(" %s%s\n\n", utils.ErrorStyle("Error: "), model.searchTermError)
		}

		return commonText
	}

	return ""
}

func createSearchInput() textinput.Model {
	input := textinput.New()
	input.Focus()
	input.CharLimit = 156
	input.Width = 20

	return input
}

func isSearchTermValid(term string) bool {
	termNoWhitespace := strings.ReplaceAll(term, " ", "")

	return len(termNoWhitespace) >= 3
}

func executeSearch(playlistId string, term string) tea.Cmd {
	return func() tea.Msg {
		return services.SearchInPlaylist(playlistId, term)
	}
}
