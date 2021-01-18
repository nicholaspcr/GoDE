package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// global flags
var np, dim, gen, execs int
var floor, ceil, crConst, fConst, pConst float64

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gode",
	Short: "differential evolution tool build in go",
	Long:  `A CLI for using the implementation of the differential evolution algorithm`,
	// todo:  Allow the user to insert his own

	//	Run: func(cmd *cobra.Command, args []string) { },
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
	cobra.OnInitialize(initConfig)

	// persistent flags
	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config",
		"",
		"config file (default is $HOME/.gode.yaml)")
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

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gode" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gode")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
