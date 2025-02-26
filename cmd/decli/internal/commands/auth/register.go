package authcmd

import (
	"errors"
	"log/slog"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// registerCmd allows registering an account in the deserver.
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Creates an account",
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.ValidateRequiredFlags()
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()

		if username == "" || password == "" || email == "" {
			return errors.New("missing neccessary fields (name,password,email)")
		}

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
	// Requirements
	registerCmd.MarkFlagRequired("username")
	registerCmd.MarkFlagRequired("email")
	registerCmd.MarkFlagRequired("password")

	// Flags
	registerCmd.Flags().StringVar(&username, "username", "", "user's name")
	registerCmd.Flags().StringVar(&password, "password", "", "user's password")
	registerCmd.Flags().StringVar(&email, "email", "", "user's email")

	// Viper binds
	viper.BindPFlag("email", registerCmd.Flags().Lookup("email"))
	viper.BindPFlag("password", registerCmd.Flags().Lookup("password"))
	viper.BindPFlag("username", registerCmd.Flags().Lookup("username"))

	// Commands
	authCmd.AddCommand(registerCmd)
}
