package commands

import (
	"log/slog"

	"github.com/nicholaspcr/GoDE/cmd/deserver/internal/config"
	"github.com/nicholaspcr/GoDE/cmd/deserver/internal/server"
	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/spf13/cobra"
)

var cfg *config.DeServer

// RootCmd is the root command for the deserver application.
var RootCmd = &cobra.Command{
	Use:   "deserver",
	Short: "deserver is a server for the de client",
	Long: `deserver is a server that implements the services described in the 
proto files. Requests can be made via gRPC or HTTP.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Checks for files, its contents overwrites the default config fields.
		cfg = config.Default()

		logger := log.New(
			log.WithWriter(cfg.Log.Writer),
			log.WithType(cfg.Log.Type),
			log.WithLevel(cfg.Log.Level),
			log.WithPrettyConfig(cfg.Log.Pretty),
		)
		slog.SetDefault(logger)

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
		ctx := cmd.Context()

		st, err := store.New(ctx)
		if err != nil {
			return err
		}

		srv, err := server.New(ctx, server.WithStore(st))
		if err != nil {
			return err
		}

		return srv.Start(ctx)
	},
}
