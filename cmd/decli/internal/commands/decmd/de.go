// Package de_cmd handles all decli commands related to the differential
// evolution API.
package decmd

import (
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/state"
	"github.com/spf13/cobra"
)

var (
	cfg *config.Config
	db  state.Operations
)

// deCmd encapsulates the DE operations.
var deCmd = &cobra.Command{
	Use:   "de",
	Short: "encapsulates DE operations",
	RunE:  func(cmd *cobra.Command, _ []string) error { return cmd.Help() },
}

// RegisterCommands adds the subset of commands into the provided cobra.Command
func RegisterCommands(root *cobra.Command) { root.AddCommand(deCmd) }

// SetupConfig sets the config of this package.
func SetupConfig(rootCfg *config.Config) { cfg = rootCfg }

// SetupStateHandler sets the state handler of this package
func SetupStateHandler(rootDB state.Operations) { db = rootDB }
