package commands

import (
	"fmt"
	"log/slog"
	_ "net/http/pprof"
	"os"

	authcmd "github.com/nicholaspcr/GoDE/cmd/decli/internal/commands/auth"
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/commands/decmd"
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/state"
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/state/sqlite"
	sharedCfg "github.com/nicholaspcr/GoDE/internal/config"
	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/spf13/cobra"
)

var (
	cfg *config.Config
	db  state.Operations
)

// Execute runs the CLI root command and handles any errors.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "decli",
	Short: "Differential evolution tool built in go",
	Long: `
A CLI for using the implementation of the differential evolution algorithm, this
allows the usage of the algorithm locally and the ability to connect to a
server.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg = config.Default()
		if err := sharedCfg.Load("decli", cfg); err != nil {
			return err
		}

		ctx := cmd.Context()

		logOpts := []log.Option{
			log.WithType(cfg.Log.Type),
			log.WithLevel(cfg.Log.Level),
			log.WithPrettyConfig(cfg.Log.Pretty),
		}
		if cfg.Log.Filename != "" {
			f, err := os.Create(cfg.Log.Filename)
			if err != nil {
				return err
			}
			logOpts = append(logOpts, log.WithWriter(f))
		}
		logger := log.New(logOpts...)
		slog.SetDefault(logger)

		var err error
		db, err = sqlite.New(ctx, cfg.State)
		if err != nil {
			return err
		}

		// NOTE: this function call has to be on the end of the PersistentPreRun.
		setupCommands()
		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		slog.Debug("Initialization of CLI:",
			slog.Any("flags", cmd.Flags()),
			slog.Any("Configuration", cfg),
		)
		return cmd.Help()
	},
}

// Sets the config and state handler for isolated command packages.
func setupCommands() {
	authcmd.SetupConfig(cfg)
	authcmd.SetupStateHandler(db)

	decmd.SetupConfig(cfg)
	decmd.SetupStateHandler(db)
}

func init() {
	// Flags
	rootCmd.PersistentFlags().String(
		"config", "", "config file (default is $HOME/.decli.yaml)",
	)

	// Commands
	authcmd.RegisterCommands(rootCmd)
	decmd.RegisterCommands(rootCmd)
}
