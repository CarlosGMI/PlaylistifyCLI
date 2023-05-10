package account

import (
	"fmt"
	"os"

	"github.com/CarlosGMI/Playlistify/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func LoginCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "login",
		Short: "Login to your Shopify account and authorize this app",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			authModel := tui.CreateAuthentication()

			if _, err := tea.NewProgram(&authModel).Run(); err != nil {
				fmt.Println("could not run program:", err)
				os.Exit(1)
			}

			return nil
		},
	}
	return command
}
