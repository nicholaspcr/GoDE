package commands

import (
	"fmt"
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
		if err := config.InitializeRoot(cmd); err != nil {
			return err
		}
		//if err := config.Unmarshal(&cfg); err != nil {
		//	return err
		//}
		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		fmt.Println("FLAGS:", cmd.Flags())
		fmt.Println("Config:", cfg)
		return cmd.Help()
	},
}

func init() {
	// Definition of commands
	RootCmd.AddCommand(localCmd)
	RootCmd.AddCommand(remoteCmd)

}
