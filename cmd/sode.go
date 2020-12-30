package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/nicholaspcr/go-de/so"
)

// local flags
var pConst float64

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
			P:     pConst,
		}
		so.Run(params)
	},
}

func init() {
	rootCmd.AddCommand(sodeCmd)
	sodeCmd.Flags().Float64Var(
		&pConst,
		"P",
		0.5,
		"P -> DE constant",
	)
}
