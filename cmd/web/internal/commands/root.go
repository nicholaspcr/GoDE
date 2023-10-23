package commands

import (
	"log/slog"
	_ "net/http/pprof"
	"os"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nicholaspcr/GoDE/cmd/web/internal"
	"github.com/nicholaspcr/GoDE/cmd/web/internal/config"
	"github.com/nicholaspcr/GoDE/cmd/web/internal/routes"
	"github.com/nicholaspcr/GoDE/internal/log"
	slogecho "github.com/samber/slog-echo"
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
		r := echo.New()

		// Add static files
		r.Static("/static", path.Join(internal.ProjectPath(), "static"))

		r.Use(
			slogecho.New(slog.Default()),
			middleware.Recover(),
		)

		// Create routes for each group.
		for groupName, groupRoutes := range routes.RouteGroups {
			group := r.Group(groupName)
			groupRoutes(group)
		}

		slog.Info(
			"Starting server on",
			slog.String("Address", cfg.Server.Address),
		)
		return r.Start(cfg.Server.Address)
	},
}

// NOTE: Add extra commands to the web CLI in here.
func init() {
	RootCmd.AddCommand(configCmd)
}
