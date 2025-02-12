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

// loginCmd encapsulates the login related operations
var loginCmd = &cobra.Command{
	Use:   "login <email> <password>",
	Short: "Log in the user's account",
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.ValidateRequiredFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// Process arguments.
		switch len(args) {
		case 0:
			break
		case 2:
			email = args[0]
			password = args[1]
		default:
			return errors.New("invalid amount of arguments")
		}

		// Validate
		if password == "" || email == "" {
			return errors.New("missing neccessary fields (email,password)")
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

		resp, err := client.Login(
			ctx,
			&api.AuthServiceLoginRequest{
				Email:    email,
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
	// Requirements
	loginCmd.MarkFlagRequired("email")
	loginCmd.MarkFlagRequired("password")

	// Flags
	loginCmd.Flags().StringVar(&email, "email", "", "user's email")
	loginCmd.Flags().StringVar(&password, "password", "", "user's password")

	// Viper binds
	viper.BindPFlag("email", loginCmd.Flags().Lookup("email"))
	viper.BindPFlag("password", loginCmd.Flags().Lookup("password"))

	// Commands
	authCmd.AddCommand(loginCmd)
}
