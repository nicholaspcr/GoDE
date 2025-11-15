package authcmd

import (
	"log/slog"

	"github.com/nicholaspcr/GoDE/cmd/decli/internal/utils"
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

		// Prompt for password securely
		password, err := utils.ReadPassword("Password: ")
		if err != nil {
			return err
		}

		conn, err := grpc.NewClient(
			cfg.Server.GRPCAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := conn.Close(); cerr != nil {
				slog.Warn("Failed to close connection", slog.String("error", cerr.Error()))
			}
		}()

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

		slog.Debug("Logged in successfully", slog.String("access_token", resp.AccessToken))
		// For now, we only store the access token. In the future, we could store the refresh token
		// and implement automatic token refresh in the CLI
		return db.SaveAuthToken(resp.AccessToken)
	},
}

func init() {
	// Flags
	loginCmd.Flags().StringVar(&username, "username", "", "user's name")

	// Requirements
	_ = loginCmd.MarkFlagRequired("username")

	// Commands
	authCmd.AddCommand(loginCmd)
}
