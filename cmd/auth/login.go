package auth

import (
	"fmt"

	"github.com/spf13/cobra"
)

func LoginCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "login",
		Short: "A brief description of your command",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("login called")
		},
	}

	return command
}
