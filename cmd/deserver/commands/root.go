package commands

import (
	"fmt"
	"os"

	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/spf13/cobra"
)

var logger *log.Logger

// Execute adds all child commands to the root command and sets flags
// appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if logger.SugaredLogger == nil {
			fmt.Printf("decli ended unexpectedly, error: %s", err)
		} else {
			logger.Error("decli ended unexpectedly", "error", err)
		}
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "deserver",
	Short: "Server for executing Differential Evolution algorithms",
	Long: `Server capable of serving multiple requests of Differential Evolution
requests, storing the values of each step and the end result in a database and
making it possible to retrieve those at any point.`,
	PersistentPreRun: func(*cobra.Command, []string) {
		logger = log.New()
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

func init() {
	testCmd.AddCommand(dbServerCmd)
	rootCmd.AddCommand(
		startCmd,
		testCmd,
	)
}
