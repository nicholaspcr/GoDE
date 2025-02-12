package authcmd

import (
	"github.com/spf13/cobra"
)

// logoutCmd encapsulates the logout related operations
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Invalidates login token saved in the local CLI state",
	RunE: func(cmd *cobra.Command, args []string) error {
		return db.InvalidateAuthToken()
	},
}

func init() {
	authCmd.AddCommand(logoutCmd)
}
