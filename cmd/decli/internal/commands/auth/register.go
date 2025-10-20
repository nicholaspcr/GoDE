package authcmd

import (
	"log/slog"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// registerCmd allows registering an account in the deserver.
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Creates an account",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()

		conn, err := grpc.NewClient(
			cfg.Server.GRPCAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := api.NewAuthServiceClient(conn)

		_, err = client.Register(
			ctx,
			&api.AuthServiceRegisterRequest{
				User: &api.User{
					Ids:      &api.UserIDs{Username: username},
					Email:    email,
					Password: password,
				},
			},
		)
		if err != nil {
			return err
		}

		slog.Info("Create account successfully")
		return nil
	},
}

func init() {
	// Flags
	registerCmd.Flags().StringVar(&username, "username", "", "user's name")
	registerCmd.Flags().StringVar(&password, "password", "", "user's password")
	registerCmd.Flags().StringVar(&email, "email", "", "user's email")

	// Requirements
	registerCmd.MarkFlagRequired("username")
	registerCmd.MarkFlagRequired("email")
	registerCmd.MarkFlagRequired("password")

	// Commands
	authCmd.AddCommand(registerCmd)
}
