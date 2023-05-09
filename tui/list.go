package tui

import (
	"fmt"

	"github.com/CarlosGMI/Playlistify/services"
	"github.com/CarlosGMI/Playlistify/utils"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type PlaylistsModel struct {
	state       string
	loader      spinner.Model
	loaderText  string
	results     TableModel
	resultsText string
}

func CreatePlaylistsModel() PlaylistsModel {
	return PlaylistsModel{
		state:      "",
		loader:     CreateSpinner(),
		loaderText: "Refreshing token...",
	}
}

func (model PlaylistsModel) Init() tea.Cmd {
	return services.InitAuthentication
}

func (model *PlaylistsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return model, tea.Quit
		default:
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
	case services.LoggedInMsg:
		model.state = utils.LoadingState
		model.loaderText = "Fetching playlists..."
		model.loader, cmd = model.loader.Update(spinner.TickMsg{})
		cmds = append(cmds, fetchPlaylists(), cmd)
	case services.PlaylistsMsg, services.AuthErrorMsg:
		model.state = "table"
		playlists, err := services.PrintPlaylists()

		if err != nil {
			model.state = utils.ErrorState
			model.resultsText = err.Error()

			return model, tea.Quit
		}

		model.results = CreateTable("PLAYLISTS", playlists, false, "")

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

func (model PlaylistsModel) View() string {
	if model.state == utils.LoadingState {
		return fmt.Sprintf("\n %s %s\n\n", model.loader.View(), model.loaderText)
	} else if model.state == utils.ErrorState {
		return fmt.Sprintf("\n %s%s\n\n", utils.ErrorStyle("Error: "), model.resultsText)
	}

	return fmt.Sprintf("\n %s\n\n", model.results.View())
}

func fetchPlaylists() tea.Cmd {
	return services.GetPlaylists
}
