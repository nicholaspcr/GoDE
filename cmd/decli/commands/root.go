package commands

import (
	"math/rand"
	"time"

	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/spf13/cobra"
)

var cfg config.Config

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "decli",
	Short: "Differential evolution tool built in go",
	Long: `
A CLI for using the implementation of the differential evolution algorithm, this
allows the usage of the algorithm locally and the ability to connect to a
server.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		rand.Seed(time.Now().UnixNano())
		cmd.SetContext(log.SetContext(cmd.Context(), log.New()))
		return config.Unmarshal(&cfg)
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(localCmd)
	RootCmd.AddCommand(remoteCmd)
}
