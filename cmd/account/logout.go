package account

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func LogoutCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "logout",
		Short: "A brief description of your command",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			inOneHour := time.Now().Unix()
			fmt.Println(inOneHour)
		},
	}

	return command
}
