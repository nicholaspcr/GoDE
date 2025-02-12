package authcmd

import (
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/state"
	"github.com/spf13/cobra"
)

var (
	cfg *config.Config
	db  state.Operations
)

// authCmd encapsulates the authentication operations.
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "encapsulates authentication operations",
	RunE:  func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
}

// RegisterCommands adds the subset of commands into the provided cobra.Command
func RegisterCommands(root *cobra.Command) { root.AddCommand(authCmd) }

// SetupConfig sets the config of this package.
func SetupConfig(rootCfg *config.Config) { cfg = rootCfg }

// SetupStateHandler sets the state handler of this package
func SetupStateHandler(rootDB state.Operations) { db = rootDB }
