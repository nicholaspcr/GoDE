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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		logger := log.New()
		slog.SetDefault(logger)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
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
