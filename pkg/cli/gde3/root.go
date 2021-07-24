package gde3

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// global flags
var (
	np, dim                 int
	gen, execs              int
	floor, ceil             []float64
	crConst, fConst, pConst float64

	mConst       int
	functionName string
	disablePlot  bool

	// pprofs
	cpuprofile string
	memprofile string

	// filename for the yaml file
	filename string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "gode",
		Short: "differential evolution tool build in go",
		Long:  `A CLI for using the implementation of the differential evolution algorithm`,

		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func init() {
	addGlobalFlags(rootCmd)

	rootCmd.AddCommand(
		modeCmd,
		scriptCmd,
	)

	modeCmd.Flags().StringVar(
		&variantName,
		"vr",
		"rand1",
		"name fo the variant to be used",
	)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addGlobalFlags(cmd *cobra.Command) {
	// persistent flags
	cmd.PersistentFlags().IntVarP(&np,
		"np",
		"n",
		100,
		"amout of elements.")

	cmd.PersistentFlags().IntVarP(&dim,
		"dim",
		"d",
		7,
		"quantity of dimension used for the problem.")

	cmd.PersistentFlags().IntVarP(&gen,
		"gen",
		"g",
		300,
		"generations of the DE")

	cmd.PersistentFlags().IntVarP(&execs,
		"execs",
		"e",
		1,
		"amount of times to run DE")

	cmd.PersistentFlags().Float64SliceVarP(&floor,
		"floor",
		"",
		[]float64{1.0},
		"floor of the float64 generator (default 0)")

	cmd.PersistentFlags().Float64SliceVarP(&ceil,
		"ceil",
		"",
		[]float64{1.0},
		"ceil of the float64 generator")

	cmd.PersistentFlags().Float64Var(&crConst,
		"CR",
		0.9,
		"CR -> DE constant")

	cmd.PersistentFlags().Float64Var(&fConst,
		"F",
		0.5,
		"F -> DE constant")

	cmd.PersistentFlags().Float64Var(
		&pConst,
		"P",
		0.2,
		"P -> DE constant",
	)

	cmd.PersistentFlags().IntVar(&mConst,
		"M",
		3,
		"M -> DE constant")

	cmd.PersistentFlags().StringVar(&functionName,
		"fn",
		"DTLZ1",
		"name of the problem to be used.")

	cmd.PersistentFlags().BoolVar(&disablePlot,
		"disable-plot",
		false,
		"to write in files the result of the gde3 to be able to plot it with the python scripts")

	cmd.PersistentFlags().StringVar(&cpuprofile,
		"cpuprofile",
		"",
		"write cpu profile to `file`")

	cmd.PersistentFlags().StringVar(&memprofile,
		"memprofile",
		"",
		"write memory profile to `file`")

	cmd.PersistentFlags().StringVar(&filename,
		"filename",
		"",
		"filename path to the yaml file that contains the values of the problem")
}
