package local

import (
	"github.com/spf13/cobra"
)

// LocalCmd represents the de command
var LocalCmd = &cobra.Command{
	Use:   "local",
	Short: "Local operations related to Differential Evolutionary algorithm",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}
