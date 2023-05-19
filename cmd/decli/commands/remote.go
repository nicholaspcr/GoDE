package commands

import (
	"github.com/spf13/cobra"
)

// remoteCmd is the command for remote operations, involves the auth procedure
// and all the necessary operations to interact with the remote server that runs
// DE related operations.
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Remote operations for interacting with the server that runs DE related operations",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}
