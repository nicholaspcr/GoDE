package commands

import (
	"github.com/nicholaspcr/GoDE/cmd/deserver/internal/server"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/spf13/cobra"
)

// StartCmd runs the server
var StartCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"run"},
	Short:   "starts a server that implements the API services",
	Long: `starts a server that implements the services described in the 
proto files. Requests can be made via gRPC or HTTP.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		st, err := store.New(ctx, cfg.Store)
		if err != nil {
			return err
		}

		srv, err := server.New(ctx, cfg.Server, server.WithStore(st))
		if err != nil {
			return err
		}

		return srv.Start(ctx)
	},
}

func init() {
	rootCmd.AddCommand(StartCmd)
}
