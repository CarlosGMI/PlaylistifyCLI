package playlist

import (
	"fmt"
	"os"

	"github.com/CarlosGMI/Playlistify/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func ListCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "List all your Spotify playlists, including collaborative playlists",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			model := tui.CreatePlaylistsModel()

			if _, err := tea.NewProgram(&model).Run(); err != nil {
				fmt.Println("could not run program:", err)
				os.Exit(1)
			}

			return nil
		},
	}

	return command
}
