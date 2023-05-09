package playlist

import (
	"fmt"
	"os"

	"github.com/CarlosGMI/Playlistify/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func SearchCommand() *cobra.Command {
	var playlistIdFlag string
	var searchTermFlag string
	command := &cobra.Command{
		Use:   "search",
		Short: "Command to look for a song/artist inside a specific playlist",
		Long: `This command will allow you to find possible duplicates in a playlist by looking up a specific term within all the tracks of the playlist.

		Usage:
		- playlistify search -p PLAYLIST_ID -s "TERM_YOU_WANT_TO_LOOK_FOR"
		- playlistify search -p PLAYLIST_ID -s "TERM_YOU_WANT_TO_LOOK_FOR"
		Example:
		  - playlistify search
		  - playlistify search -p 10 -s "Linkin"
		  - playlistify search -p 2 -s "two hearts"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var model tui.SearchModel
			hasPlaylist := cmd.Flags().Changed("playlist")
			hasSearchTerm := cmd.Flags().Changed("term")

			if !hasPlaylist && !hasSearchTerm {
				model = tui.CreateSearchModel(true, "", "")
			} else {
				model = tui.CreateSearchModel(false, playlistIdFlag, searchTermFlag)
			}

			if _, err := tea.NewProgram(model).Run(); err != nil {
				fmt.Println("could not run program:", err)
				os.Exit(1)
			}

			return nil
		},
	}

	command.Flags().StringVarP(&playlistIdFlag, "playlist", "p", "", "Playlist ID (required)")
	command.Flags().StringVarP(&searchTermFlag, "term", "t", "", "Term to search for (required)")
	command.MarkFlagsRequiredTogether("playlist", "term")

	return command
}
