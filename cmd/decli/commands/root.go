package commands

import (
	"fmt"
	"os"

	"github.com/nicholaspcr/GoDE/internal/config"
	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// loggers for the CLI
	logger log.Logger
	conf   *config.Config

	// pprofs
	cpuprofile string
	memprofile string
)

// Execute adds all child commands to the root command and sets flags
// appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if logger.IsNil() {
			fmt.Printf("decli ended unexpectedly, error: %s", err)
		} else {
			logger.Error("decli ended unexpectedly", "error", err)
		}
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "decli",
	Short: "Differential evolution tool build in go",
	Long:  `A CLI for using the implementation of the differential evolution algorithm`,
	PreRunE: func(*cobra.Command, []string) error {
		logger = log.New()
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(
		modeCmd,
	)
	conf = config.New()
	setDefaults(conf)
	conf.SetConfigFile(".config")
	conf.SetConfigType("yaml")
	_ = conf.ReadInConfig() // ignore if config file doesn't exist
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
