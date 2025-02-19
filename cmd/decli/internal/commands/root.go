package commands

import (
	"fmt"
	"log/slog"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"

	"github.com/nicholaspcr/GoDE/internal/log"

	authcmd "github.com/nicholaspcr/GoDE/cmd/decli/internal/commands/auth"
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/commands/decmd"
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/state"
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/state/sqlite"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile        string
	memProfileFile string
	cpuProfileFile string
	cfg            *config.Config
	db             state.Operations
)

// Executes the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "decli",
	Short: "Differential evolution tool built in go",
	Long: `
A CLI for using the implementation of the differential evolution algorithm, this
allows the usage of the algorithm locally and the ability to connect to a
server.
`,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) (err error) {
		ctx := cmd.Context()

		logOpts := []log.Option{
			log.WithType(cfg.Log.Type),
			log.WithLevel(cfg.Log.Level),
			log.WithPrettyConfig(cfg.Log.Pretty),
		}
		if cfg.Log.Filename != "" {
			f, err := os.Create(cfg.Log.Filename)
			if err != nil {
				return err
			}
			logOpts = append(logOpts, log.WithWriter(f))
		}
		logger := log.New(logOpts...)
		slog.SetDefault(logger)

		db, err = sqlite.New(ctx, cfg.Sqlite)
		if err != nil {
			return err
		}

		cpuProfile, err := os.Create(cpuProfileFile)
		if err != nil {
			return err
		}
		if err := pprof.StartCPUProfile(cpuProfile); err != nil {
			return err
		}

		// NOTE: this function call has to be on the end of the PersistentPreRun.
		setupCommands()
		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		slog.Debug("Initialization of CLI:",
			slog.Any("flags", cmd.Flags()),
			slog.Any("Configuration", cfg),
		)
		return cmd.Help()
	},
	PersistentPostRunE: func(_ *cobra.Command, _ []string) error {
		pprof.StopCPUProfile()

		memProfile, err := os.Create(memProfileFile)
		if err != nil {
			return err
		}
		return pprof.WriteHeapProfile(memProfile)
	},
}

// Sets the config and state handler for isolated command packages.
func setupCommands() {
	authcmd.SetupConfig(cfg)
	authcmd.SetupStateHandler(db)

	decmd.SetupConfig(cfg)
	decmd.SetupStateHandler(db)
}

func init() {
	cobra.OnInitialize(initConfig)

	// Flags
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "", "config file (default is $HOME/.decli.yaml)",
	)
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	rootCmd.PersistentFlags().StringVar(
		&cpuProfileFile, "cpu-profile-file",
		".dev/cli/cpuprofile", "cpu profile filename",
	)
	rootCmd.PersistentFlags().StringVar(
		&memProfileFile, "mem-profile-file",
		".dev/cli/memprofile", "mem profile filename",
	)

	// Viper binds
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.BindPFlag(
		"cpuProfileFile", rootCmd.PersistentFlags().Lookup("cpu-profile-file"),
	)
	viper.BindPFlag(
		"memProfileFile", rootCmd.PersistentFlags().Lookup("mem-profile-file"),
	)

	// Commands
	authcmd.RegisterCommands(rootCmd)
	decmd.RegisterCommands(rootCmd)

	rootCmd.AddCommand(localCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".decli" (no extention)
		viper.AddConfigPath(home)
		viper.AddConfigPath(".env")
		viper.AddConfigPath(".")

		viper.SetConfigType("yaml")
		viper.SetConfigName(".decli")
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
