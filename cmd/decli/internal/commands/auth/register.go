package authcmd

import (
	"log/slog"

	"github.com/nicholaspcr/GoDE/cmd/decli/internal/utils"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var regPassword string

// registerCmd allows registering an account in the deserver.
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Creates an account",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()

		// Use provided password or prompt for it
		pwd := regPassword
		if pwd == "" {
			var err error
			pwd, err = utils.ReadPasswordWithConfirmation(
				"Password: ",
				"Confirm password: ",
			)
			if err != nil {
				return err
			}
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

		_, err = client.Register(
			ctx,
			&api.AuthServiceRegisterRequest{
				User: &api.User{
					Ids:      &api.UserIDs{Username: username},
					Email:    email,
					Password: pwd,
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
	registerCmd.Flags().StringVar(&email, "email", "", "user's email")
	registerCmd.Flags().StringVar(&regPassword, "password", "", "user's password (optional, will prompt if not provided)")

	// Requirements
	_ = registerCmd.MarkFlagRequired("username")
	_ = registerCmd.MarkFlagRequired("email")

	// Commands
	authCmd.AddCommand(registerCmd)
}
