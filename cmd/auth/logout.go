package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

func LogoutCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "logout",
		Short: "A brief description of your command",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("logout called")
		},
	}

	return command
}
