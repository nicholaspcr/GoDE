package authcmd

import (
	"log/slog"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// loginCmd encapsulates the login related operations
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in the user's account",
	RunE: func(cmd *cobra.Command, args []string) error {
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
				Username: username,
				Password: password,
			},
		)
		if err != nil {
			return err
		}

		slog.Debug("Logged in successfully", slog.String("token", resp.Token))
		return db.SaveAuthToken(resp.Token)
	},
}

func init() {
	// Flags
	loginCmd.Flags().StringVar(&username, "username", "", "user's name")
	loginCmd.Flags().StringVar(&password, "password", "", "user's password")

	// Requirements
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")

	// Commands
	authCmd.AddCommand(loginCmd)
}
