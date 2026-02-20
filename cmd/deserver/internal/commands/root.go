package commands

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/nicholaspcr/GoDE/cmd/deserver/internal/config"
	sharedCfg "github.com/nicholaspcr/GoDE/internal/config"
	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/spf13/cobra"
)

var (
	cfg *config.Config
)

// Execute runs the server root command and handles any errors.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// rootCmd is the root command for the deserver application.
var rootCmd = &cobra.Command{
	Use:   "deserver",
	Short: "deserver is API to create and administer differential algorithms",
	Long: `deserver is a server that implements the services described in the
proto files found on the API folder. Requests can be made via gRPC or HTTP.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg = config.Default()
		configPath, _ := cmd.Flags().GetString("config")
		if err := sharedCfg.Load("deserver", configPath, cfg); err != nil {
			return err
		}

		logger := log.New(
			log.WithType(cfg.Log.Type),
			log.WithLevel(cfg.Log.Level),
			log.WithPrettyConfig(cfg.Log.Pretty),
		)
		slog.SetDefault(logger)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	// Flags
	rootCmd.PersistentFlags().StringP(
		"config", "c", "", "config file (default is $HOME/.deserver.yaml)",
	)
}
