package commands

import (
	"log/slog"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// authCmd encapsulates the authentication operations.
	authCmd = &cobra.Command{
		Use:   "auth",
		Short: "encapsulates authentication operations",
		RunE:  func(_ *cobra.Command, _ []string) error { return nil },
	}

	// registerCmd allows registering an account in the deserver.
	registerCmd = &cobra.Command{
		Use:   "register <username> <password>",
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
						Ids:      &api.UserIDs{Email: "nicholaspcr@gmail.com"},
						Name:     "nicholaspcr",
						Password: "password",
					},
				},
			)
			if err != nil {
				return err
			}

			return nil
		},
	}

	// loginCmd encapsulates the login related operations
	loginCmd = &cobra.Command{
		Use:   "login <username> <password>",
		Short: "Logs in the user's account",
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

			resp, err := client.Login(
				ctx,
				&api.AuthServiceLoginRequest{
					Email:    "nicholaspcr@gmail.com",
					Password: "password",
				},
			)
			if err != nil {
				return err
			}

			slog.Debug(
				"Logged in successfully",
				slog.String("token", resp.Token),
			)

			return nil
		},
	}
)

func init() {
	authCmd.AddCommand(registerCmd)
	authCmd.AddCommand(loginCmd)
	RootCmd.AddCommand(authCmd)
}
