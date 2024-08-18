package authcmd

import (
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	"github.com/spf13/cobra"
)

var cfg *config.Config

// authCmd encapsulates the authentication operations.
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "encapsulates authentication operations",
	RunE:  func(_ *cobra.Command, _ []string) error { return nil },
}

// RegisterCommands adds the subset of commands into the provided cobra.Command
func RegisterCommands(root *cobra.Command) { root.AddCommand(authCmd) }

// SetupConfig sets the config of this package to be the same as the geral config.
func SetupConfig(rootCfg *config.Config) { cfg = rootCfg }
