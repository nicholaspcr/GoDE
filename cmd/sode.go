package cmd

import (
	so "github.com/nicholaspcr/IC-GDE3/pkg/sode"
	"github.com/spf13/cobra"
)

// sodeCmd represents the sode command
var sodeCmd = &cobra.Command{
	Use:   "single",
	Short: "Single-objective implementation of DE",
	Long:  `Implementation of the DE algorithm in the simpler sense, where there is only one objective funtion to be minimized.`,
	Run: func(cmd *cobra.Command, args []string) {
		params := so.Params{
			NP:    np,
			DIM:   dim,
			GEN:   gen,
			EXECS: execs,
			FLOOR: floor,
			CEIL:  ceil,
			CR:    crConst,
			F:     fConst,
		}
		so.Run(params)
	},
}

func init() {
	rootCmd.AddCommand(sodeCmd)
}
