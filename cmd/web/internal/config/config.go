package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const appname = "web"

type (
	// Config is a set of values that are necessary to execute an Differential Evolutionary algorithm.
	Config struct {
		Server Server `json:"server" yaml:"server"`
		Logger Log    `json:"logger" yaml:"logger"`
	}

	// Server is a set of values that are necessary to configure the server.
	Server struct {
		Address string `json:"address" yaml:"address"`
		Cors    Cors
	}

	// Cors is a set of values that are necessary to configure the CORS.
	Cors struct {
		AllowOrigins     []string      `json:"allow_origins" yaml:"allow_origins"`
		AllowMethods     []string      `json:"allow_methods" yaml:"allow_methods"`
		AllowHeaders     []string      `json:"allow_headers" yaml:"allow_headers"`
		ExposeHeaders    []string      `json:"expose_headers" yaml:"expose_headers"`
		AllowCredentials bool          `json:"allow_credentials" yaml:"allow_credentials"`
		MaxAge           time.Duration `json:"max_age" yaml:"max_age"`
	}

	// Log is a set of values that are necessary to configure the logger.
	Log struct {
		*log.Config `json:"config" yaml:"config"`
		Filename    string `json:"filename" yaml:"filename"`
	}
)

// DefaultConfig is the default configuration for the web server.
var DefaultConfig = Config{
	Server: Server{
		Address: ":8080",
		// NOTE: This is a very permissive CORS configuration. You should change
		// these values to match your needs.
		Cors: Cors{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"*"},
			AllowCredentials: true,
			MaxAge:           300,
		},
	},
	Logger: Log{
		Config: &log.Config{Writer: os.Stdout},
	},
}

// InitializeRoot initializes the configuration for the root command.
func InitializeRoot(_ *cobra.Command, cfg *Config) error {
	v := viper.New()

	// Configuration filename and type.
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// Configuration search path.
	v.AddConfigPath("/etc/decli/")
	v.AddConfigPath("$HOME/.decli")
	v.AddConfigPath(".env")
	v.AddConfigPath(".")

	// Fetch configuration file contents.
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Warn("config file not found, using default configuration")
			return nil
		}
		return err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return err
	}

	v.AutomaticEnv()

	// TODO: Add cmd flags to override config values.
	return nil
}
