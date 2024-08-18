package commands

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/pprof"

	"github.com/nicholaspcr/GoDE/cmd/deserver/internal/config"
	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile        string
	memProfileFile string
	cpuProfileFile string
	cfg            *config.Config
)

// Executes the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// rootCmd is the root command for the deserver application.
var rootCmd = &cobra.Command{
	Use:   "deserver",
	Short: "deserver is API to create and administer differential algorithms",
	Long: `deserver is a server that implements the services described in the
proto files found on the API folder. Requests can be made via gRPC or HTTP.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logger := log.New(
			log.WithWriter(cfg.Log.Writer),
			log.WithType(cfg.Log.Type),
			log.WithLevel(cfg.Log.Level),
			log.WithPrettyConfig(cfg.Log.Pretty),
		)
		slog.SetDefault(logger)

		cpuProfile, err := os.Create(cpuProfileFile)
		if err != nil {
			return err
		}
		return pprof.StartCPUProfile(cpuProfile)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		pprof.StopCPUProfile()

		memProfile, err := os.Create(memProfileFile)
		if err != nil {
			return err
		}
		return pprof.WriteHeapProfile(memProfile)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Flags
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "", "config file (default is $HOME/.deserver.yaml)",
	)
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	rootCmd.PersistentFlags().StringVar(
		&cpuProfileFile, "cpu-profile-file",
		".dev/server/cpuprofile", "cpu profile filename",
	)
	rootCmd.PersistentFlags().StringVar(
		&memProfileFile, "mem-profile-file",
		".dev/server/memprofile", "mem profile filename",
	)

	// Viper binds
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.BindPFlag(
		"cpuProfileFile", rootCmd.PersistentFlags().Lookup("cpu-profile-file"),
	)
	viper.BindPFlag(
		"memProfileFile", rootCmd.PersistentFlags().Lookup("mem-profile-file"),
	)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".deserver" (no extention)
		viper.AddConfigPath(home)
		viper.AddConfigPath(".env")
		viper.AddConfigPath(".")

		viper.SetConfigType("yaml")
		viper.SetConfigName(".deserver")
	}

	viper.AutomaticEnv()
	cfg = config.Default()

	err := viper.ReadInConfig()

	// No error, using the config config from viper. Unmarshal and return.
	if err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		cobra.CheckErr(viper.Unmarshal(&cfg))
		return
	}
	// If the error isn't config not found then CheckErr.
	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		cobra.CheckErr(err)
	}
}
