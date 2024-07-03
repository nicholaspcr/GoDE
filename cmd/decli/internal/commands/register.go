package commands

import (
	"github.com/spf13/cobra"
)

// registerCmd allows registering an account in the deserver.
var registerCmd = &cobra.Command{
	Use:   "register <username> <password>",
	Short: "Creates an account",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return nil
	},
}
