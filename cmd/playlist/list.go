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
		Short: "A brief description of your command",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			m := tui.CreateSpinner("Loading...")

			if _, err := tea.NewProgram(m).Run(); err != nil {
				fmt.Println("could not run program:", err)
				os.Exit(1)
			}

			fmt.Println("perfect ed sehraan")
			// if err := services.IsAuthenticated(); err != nil {
			// 	return err
			// }

			// if err := services.GetPlaylists(); err != nil {
			// 	return err
			// }

			// if err := services.PrintPlaylists(); err != nil {
			// 	return err
			// }

			return nil
		},
	}

	return command
}
