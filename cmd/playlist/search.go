package playlist

import (
	"errors"
	"strconv"
	"strings"

	"github.com/CarlosGMI/Playlistify/services"
	"github.com/spf13/cobra"
)

func SearchCommand() *cobra.Command {
	var playlistIdFlag string
	var searchTermFlag string
	command := &cobra.Command{
		Use:   "search",
		Short: "Command to look for a song/artist inside a specific playlist",
		Long: `This command will allow you to find possible duplicates in a playlist by looking up a specific term within all the tracks of the playlist.

		Usage: playlistify search -p PLAYLIST_ID -s "TERM_YOU_WANT_TO_LOOK_FOR"
		Example:
		  - playlistify search -p 10 -s "Linkin"
		  - playlistify search -p 2 -s "two hearts"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			termNoWhitespace := strings.ReplaceAll(searchTermFlag, " ", "")

			if len(termNoWhitespace) < 3 {
				return errors.New("the search term must be longer than 3 characters (excluding white spaces)")
			}

			if err := services.IsAuthenticated(); err != nil {
				return err
			}

			playlistId, err := strconv.Atoi(playlistIdFlag)

			if err != nil {
				return err
			}

			if err := services.SearchInPlaylist(playlistId, strings.ToLower(searchTermFlag)); err != nil {
				return err
			}

			return nil
		},
	}

	command.Flags().StringVarP(&playlistIdFlag, "playlist", "p", "", "Playlist ID (required)")
	command.MarkFlagRequired("playlist")
	command.Flags().StringVarP(&searchTermFlag, "term", "t", "", "Term to search for (required)")
	command.MarkFlagRequired("term")

	return command
}
