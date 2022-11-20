package commands

import (
	"math/rand"
	"time"

	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/spf13/cobra"

	_ "github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
)

var (
	// pprofs
	cpuprofile string
	memprofile string
)

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "decli",
	Short: "Differential evolution tool build in go",
	Long:  `A CLI for using the implementation of the differential evolution algorithm`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		rand.Seed(time.Now().UnixNano())
		cmd.SetContext(log.New().SetContext(cmd.Context()))
		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(modeCmd)

}
