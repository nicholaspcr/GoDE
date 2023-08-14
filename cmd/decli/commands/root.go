package commands

import (
	"log/slog"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"

	"github.com/nicholaspcr/GoDE/internal/log"

	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	"github.com/spf13/cobra"
)

var cfg config.Config

var memProfile, cpuProfile *os.File

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "decli",
	Short: "Differential evolution tool built in go",
	Long: `
A CLI for using the implementation of the differential evolution algorithm, this
allows the usage of the algorithm locally and the ability to connect to a
server.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Add configuration to the logger
		// Configuration has to be parsed and checked if any of the logger
		// fields are set, if so, then the logger has to be configured using
		// log.WithX field.
		logger := log.New()
		slog.SetDefault(logger)

		if err := config.InitializeRoot(cmd); err != nil {
			return err
		}
		cfg = config.DefaultConfig
		if err := config.Unmarshal(&cfg); err != nil {
			return err
		}
		var err error
		cpuProfile, err = os.Create("cpuprofile")
		if err != nil {
			return err
		}
		memProfile, err = os.Create("memprofile")
		if err != nil {
			return err
		}
		return pprof.StartCPUProfile(cpuProfile)
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		logger := slog.Default()
		logger.Debug("Initialization of CLI:",
			slog.Any("flags", cmd.Flags()),
			slog.Any("Configuration", cfg),
		)
		return cmd.Help()
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		pprof.StopCPUProfile()
		return pprof.WriteHeapProfile(memProfile)
	},
}

func init() {
	// Definition of commands
	RootCmd.AddCommand(localCmd)
	RootCmd.AddCommand(remoteCmd)
}
