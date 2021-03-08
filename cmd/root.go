package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// global flags
var (
	np, dim                 int
	gen, execs              int
	floor, ceil             float64
	crConst, fConst, pConst float64

	// pprofs
	cpuprofile string
	memprofile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gode",
	Short: "differential evolution tool build in go",
	Long:  `A CLI for using the implementation of the differential evolution algorithm`,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	// persistent flags
	rootCmd.PersistentFlags().IntVarP(&np,
		"np",
		"n",
		100,
		"amout of elements.")

	rootCmd.PersistentFlags().IntVarP(&dim,
		"dim",
		"d",
		5,
		"quantity of dimension used for the problem.")

	rootCmd.PersistentFlags().IntVarP(&gen,
		"gen",
		"g",
		300,
		"generations of the DE")

	rootCmd.PersistentFlags().IntVarP(&execs,
		"execs",
		"e",
		1,
		"amount of times to run DE")

	rootCmd.PersistentFlags().Float64Var(&floor,
		"floor",
		0.0,
		"floor of the float64 generator (default 0)")

	rootCmd.PersistentFlags().Float64Var(&ceil,
		"ceil",
		1.0,
		"ceil of the float64 generator")

	rootCmd.PersistentFlags().Float64Var(&crConst,
		"CR",
		0.9,
		"CR -> DE constant")

	rootCmd.PersistentFlags().Float64Var(&fConst,
		"F",
		0.5,
		"F -> DE constant")

	rootCmd.PersistentFlags().Float64Var(
		&pConst,
		"P",
		0.2,
		"P -> DE constant",
	)

	rootCmd.PersistentFlags().StringVar(&cpuprofile,
		"cpuprofile",
		"",
		"write cpu profile to `file`")

	rootCmd.PersistentFlags().StringVar(&memprofile,
		"memprofile",
		"",
		"write memory profile to `file`")

}
