package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ConfigCmd handles commands related to the configuration of the server.
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Handles configuration for the deserver",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgString, err := cfg.StringifyJSON()
		if err != nil {
			return err
		}
		fmt.Println("Printing the deserver configuration in JSON format")
		fmt.Println(cfgString)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(ConfigCmd)
}
