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
	Use:   "login <username> <password>",
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
			username = args[0]
			password = args[1]
		default:
			return errors.New("invalid amount of arguments")
		}

		// Validate
		if password == "" || username == "" {
			return errors.New("missing neccessary fields (username,password)")
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
	// Requirements
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")

	// Flags
	loginCmd.Flags().StringVar(&username, "username", "", "user's name")
	loginCmd.Flags().StringVar(&password, "password", "", "user's password")

	// Viper binds
	viper.BindPFlag("username", loginCmd.Flags().Lookup("username"))
	viper.BindPFlag("password", loginCmd.Flags().Lookup("password"))

	// Commands
	authCmd.AddCommand(loginCmd)
}
