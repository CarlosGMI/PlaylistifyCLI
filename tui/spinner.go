package tui

import (
	"fmt"

	"github.com/CarlosGMI/Playlistify/utils"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SpinnerModel struct {
	spinner spinner.Model
	text    string
}

func CreateSpinner(spinnerText string) SpinnerModel {
	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(utils.ColorSpotifyGreen))
	newSpinner.Spinner = spinner.Dot

	return SpinnerModel{
		spinner: newSpinner,
		text:    spinnerText,
	}
}

func (model SpinnerModel) Init() tea.Cmd {
	return model.spinner.Tick
}

func (model SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return model, tea.Quit
		default:
			return model, nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		model.spinner, cmd = model.spinner.Update(msg)
		return model, cmd
	default:
		return model, nil
	}
}

func (model SpinnerModel) View() (s string) {
	s += fmt.Sprintf("\n %s %s\n\n", model.spinner.View(), model.text)

	return
}
