package tui

import (
	"fmt"

	"github.com/CarlosGMI/Playlistify/services"
	"github.com/CarlosGMI/Playlistify/utils"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type AuthModel struct {
	state       string
	loader      spinner.Model
	loaderText  string
	resultsText string
}

func CreateAuthentication() AuthModel {
	return AuthModel{
		state:      utils.LoadingState,
		loader:     CreateSpinner(),
		loaderText: "Authenticating...",
	}
}

func (model AuthModel) Init() tea.Cmd {
	return tea.Batch(model.loader.Tick, services.InitAuthentication)
}

func (model *AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return model, tea.Quit
		default:
			if model.state == utils.ErrorState {
				return model, tea.Quit
			}

			return model, nil
		}
	case services.NotAuthenticatedMsg:
		if msg.ErrorType == utils.NotLoggedInCode {
			model.loaderText = "Authorizing..."
			model.loader, cmd = model.loader.Update(spinner.TickMsg{})
			cmds = append(cmds, authenticate(), cmd)
		} else if msg.ErrorType == utils.ExpiredTokenCode {
			model.loaderText = "Refreshing token..."
			model.loader, cmd = model.loader.Update(spinner.TickMsg{})
			cmds = append(cmds, refreshAuth(), cmd)
		}
	case services.AuthorizedMsg:
		model.loaderText = "Logging in..."
		model.loader, cmd = model.loader.Update(spinner.TickMsg{})
		cmds = append(cmds, login(), cmd)
	case services.LoggedInMsg:
		model.loaderText = "Fetching user information..."
		model.loader, cmd = model.loader.Update(spinner.TickMsg{})
		cmds = append(cmds, fetchUser(), cmd)
	case services.LoggedInUserMsg:
		model.state = utils.SuccessState
		model.resultsText = msg.Message

		return model, tea.Quit
	case services.AuthErrorMsg:
		if msg.ErrorType == utils.AlreadyLoggedInCode {
			model.loaderText = "Fetching user information..."
			model.loader, cmd = model.loader.Update(spinner.TickMsg{})
			cmds = append(cmds, fetchUser(), cmd)
		} else {
			model.state = utils.ErrorState
			model.resultsText = msg.Message

			return model.Update(tea.KeyMsg{})
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		model.loader, cmd = model.loader.Update(msg)

		return model, cmd
	default:
		return model, nil
	}

	return model, tea.Batch(cmds...)
}

func (model AuthModel) View() string {
	if model.state == utils.LoadingState {
		return fmt.Sprintf("\n %s %s\n\n", model.loader.View(), model.loaderText)
	} else if model.state == utils.ErrorState {
		return fmt.Sprintf("\n %s%s\n\n", utils.ErrorStyle("Error: "), model.resultsText)
	}

	return fmt.Sprintf("\n %s\n\n", model.resultsText)
}

func authenticate() tea.Cmd {
	return services.Authenticate
}

func login() tea.Cmd {
	return services.Login
}

func fetchUser() tea.Cmd {
	return services.FetchAuthenticatedUser
}

func refreshAuth() tea.Cmd {
	return services.RefreshAuthorization
}
