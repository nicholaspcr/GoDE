package commands

import (
	"github.com/spf13/cobra"
)

// localCmd represents the de command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Local operations related to Differential Evolutionary algorithm",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

// localRunCmd represents the run command for local operations.
var localRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a local Differential Evolutionary algorithm",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}
