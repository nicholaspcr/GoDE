package commands

import (
	"log/slog"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"

	"github.com/nicholaspcr/GoDE/internal/log"

	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/utils"
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
		if err := config.InitializeRoot(cmd, &cfg); err != nil {
			return err
		}

		logCfg := cfg.Logger.Config
		if logCfg != nil && cfg.Logger.Filename != "" {
			f, err := os.Create(cfg.Logger.Filename)
			if err != nil {
				return err
			}
			logCfg.Writer = f
		}
		logger := log.New(utils.LogOptionsFromConfig(logCfg)...)
		slog.SetDefault(logger)

		slog.Info("Initialization of CLI:",
			slog.Any("Configuration", cfg),
		)

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
		slog.Debug("Initialization of CLI:",
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
