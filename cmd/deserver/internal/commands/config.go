package commands

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var (
	ofJson bool
	ofYaml bool
)

// configCmd handles commands related to the configuration of the server.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Handles configuration for the deserver",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := cmd.ValidateFlagGroups(); err != nil {
			return err
		}
		if err := cmd.ValidateRequiredFlags(); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfgString string
		var err error

		if ofJson {
			cfgString, err = cfg.StringifyJSON()
		}
		if ofYaml {
			cfgString, err = cfg.StringifyYAML()
		}

		if err != nil {
			return err
		}
		slog.Debug("Printing the deserver configuration")
		fmt.Println(cfgString)
		return nil
	},
}

func init() {
	configCmd.Flags().BoolVar(&ofJson, "json", true, "Output in JSON")
	configCmd.Flags().BoolVar(&ofYaml, "yaml", false, "Output in YAML")
	configCmd.MarkFlagsMutuallyExclusive("json", "yaml")
	rootCmd.AddCommand(configCmd)
}
