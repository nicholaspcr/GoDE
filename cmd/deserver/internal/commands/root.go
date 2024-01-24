package commands

import (
	"log/slog"

	"github.com/nicholaspcr/GoDE/cmd/deserver/internal/server"
	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/spf13/cobra"
)

// RootCmd is the root command for the deserver application.
var RootCmd = &cobra.Command{
	Use:   "deserver",
	Short: "deserver is a server for the de client",
	Long: `deserver is a server that implements the services described in the 
proto files. Requests can be made via gRPC or HTTP.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		// Set default logger.
		logger := log.New()
		slog.SetDefault(logger)

		st, err := store.New(ctx)
		if err != nil {
			return err
		}
		srv := server.New(ctx, st)

		return srv.Start(ctx)
	},
}
