package commands

import (
	"log/slog"
	_ "net/http/pprof"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nicholaspcr/GoDE/cmd/web/internal/config"
	"github.com/nicholaspcr/GoDE/cmd/web/internal/routes"
	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/spf13/cobra"
)

var cfg = config.DefaultConfig

// RootCmd represents the base command when called without any subcommands.
var RootCmd = &cobra.Command{
	Use:   "web",
	Short: "A CLI to initialize the web server",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		// Set default logger.
		logger := log.New()
		slog.SetDefault(logger)

		// Read configuration from file and environment variables.
		if err := config.InitializeRoot(cmd, &cfg); err != nil {
			slog.Error("Failed to initialize configuration:", err)
			return err
		}

		// Set logger from configuration.
		logCfg := cfg.Logger.Config
		if logCfg != nil && cfg.Logger.Filename != "" {
			f, err := os.Create(cfg.Logger.Filename)
			if err != nil {
				return err
			}
			logCfg.Writer = f
		}
		// Create new logger from configuration and set it as default.
		logger = log.New(logOptionsFromConfig(logCfg)...)
		slog.SetDefault(logger)

		slog.Info(
			"Initialization of Web server:",
			slog.Any("Configuration", cfg),
		)
		return nil
	},
	RunE: func(*cobra.Command, []string) error {
		r := gin.Default()

		basePath := r.Group("")

		// Set default middlewares.
		basePath.Use(
			cors.New(cors.Config{
				AllowMethods:     cfg.Server.Cors.AllowHeaders,
				AllowOrigins:     cfg.Server.Cors.AllowOrigins,
				AllowHeaders:     cfg.Server.Cors.AllowHeaders,
				ExposeHeaders:    cfg.Server.Cors.ExposeHeaders,
				AllowCredentials: cfg.Server.Cors.AllowCredentials,
				MaxAge:           cfg.Server.Cors.MaxAge,
			}),
		)

		// Invoke route definitions from routes package.
		for _, route := range routes.Routes {
			route(basePath)
		}

		slog.Info(
			"Starting server on",
			slog.String("Address", cfg.Server.Address),
		)
		return r.Run(cfg.Server.Address)
	},
}

// NOTE: Add extra commands to the web CLI in here.
func init() {
	RootCmd.AddCommand(configCmd)
}
