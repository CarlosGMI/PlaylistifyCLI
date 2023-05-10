package account

import (
	"github.com/CarlosGMI/Playlistify/services"
	"github.com/spf13/cobra"
)

func LogoutCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "logout",
		Short: "Log out from your current Spotify account",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			services.EmptyTokenInformation()
		},
	}

	return command
}
