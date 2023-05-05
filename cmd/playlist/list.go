package playlist

import (
	"github.com/spf13/cobra"
)

func ListCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "A brief description of your command",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
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
