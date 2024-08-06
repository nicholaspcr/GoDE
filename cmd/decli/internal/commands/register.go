package commands

import (
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

// registerCmd allows registering an account in the deserver.
var registerCmd = &cobra.Command{
	Use:   "register <username> <password>",
	Short: "Creates an account",
	RunE: func(cmd *cobra.Command, _ []string) error {
		registerAddr := cfg.Server.HTTPAddr + "/register"

		_, err := http.PostForm(registerAddr, url.Values{
			"name":     []string{"nicholaspcr"},
			"password": []string{"password"},
			"email":    []string{"nicholaspcr@gmail.com"},
		})
		return err
	},
}

func init() {
	RootCmd.AddCommand(registerCmd)
}
