package commands

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/spf13/cobra"
)

var (
	// pprofs
	cpuprofile string
	memprofile string
)

// Execute adds all child commands to root and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("decli ended unexpectedly, error: %s", err)
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "decli",
	Short: "Differential evolution tool build in go",
	Long:  `A CLI for using the implementation of the differential evolution algorithm`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		rand.Seed(time.Now().UnixNano())
		cmd.SetContext(log.New().SetContext(cmd.Context()))
		// TODO NICK: Add logger
		//if err := viper.ReadInConfig(); err != nil {
		//	return err
		//}
		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(modeCmd)
	addGlobalFlags(rootCmd)
	// TODO NICK: Add viper config to flagset
	// conf = config.New()
	// setDefaults(conf)
	// conf.SetConfigFile(".config")
	// conf.SetConfigType("yaml")
	// _ = conf.ReadInConfig() // ignore if config file doesn't exist
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
		DefaultDimensions.Floors,
		"floor of the float64 generator (default 0)")
	cmd.PersistentFlags().Float64SliceVarP(&ceil,
		"ceil",
		"",
		DefaultDimensions.Ceils,
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
	cmd.PersistentFlags().StringVar(&problemName,
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
