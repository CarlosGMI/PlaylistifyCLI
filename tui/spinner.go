package tui

import (
	"github.com/CarlosGMI/Playlistify/utils"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

func CreateSpinner() spinner.Model {
	newSpinner := spinner.New()
	newSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(utils.ColorSpotifyGreen))
	newSpinner.Spinner = spinner.Dot

	return newSpinner
}
