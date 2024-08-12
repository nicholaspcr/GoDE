package commands

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"
)

// registerCmd allows registering an account in the deserver.
var registerCmd = &cobra.Command{
	Use:   "register <username> <password>",
	Short: "Creates an account",
	RunE: func(cmd *cobra.Command, _ []string) error {
		registerAddr := cfg.Server.HTTPAddr + "/register"

		data := struct {
			Name     string `json:"name"`
			Password string `json:"password"`
			Email    string `json:"email"`
		}{"nicholaspcr", "password", "nicholaspcr@gmail.com"}
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}

		res, err := http.DefaultClient.Post(
			registerAddr,
			"application/json",
			bytes.NewBuffer(b),
		)
		if err != nil {
			return err
		}
		defer func() { res.Body.Close() }()

		return nil
	},
}

func init() {
	RootCmd.AddCommand(registerCmd)
}
