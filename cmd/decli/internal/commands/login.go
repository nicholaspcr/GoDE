package commands

import (
	"github.com/spf13/cobra"
)

// loginCmd encapsulates the login related operations
var loginCmd = &cobra.Command{
	Use:   "login <username> <password>",
	Short: "Logs in the user's account",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return nil
	},
}
