package commands

import (
	"fmt"
	"log/slog"

	"github.com/nicholaspcr/GoDE/internal/migrations"
	"github.com/spf13/cobra"
)

var (
	migrateSteps int
)

// migrateCmd handles database migrations
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  `Manage database schema migrations (up, down, version)`,
}

// migrateUpCmd applies all pending migrations
var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	Long:  `Apply all pending database migrations to bring schema up to date`,
	RunE: func(cmd *cobra.Command, args []string) error {
		databaseURL := cfg.Store.ConnectionString()
		if databaseURL == "" {
			return fmt.Errorf("database URL not configured")
		}

		slog.Info("Applying migrations...")
		if err := migrations.Run(databaseURL); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}

		slog.Info("Migrations completed successfully")
		return nil
	},
}

// migrateDownCmd rolls back migrations
var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback migrations",
	Long:  `Rollback the specified number of database migrations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		databaseURL := cfg.Store.ConnectionString()
		if databaseURL == "" {
			return fmt.Errorf("database URL not configured")
		}

		if migrateSteps <= 0 {
			return fmt.Errorf("steps must be greater than 0")
		}

		slog.Info(fmt.Sprintf("Rolling back %d migration(s)...", migrateSteps))
		if err := migrations.Rollback(databaseURL, migrateSteps); err != nil {
			return fmt.Errorf("rollback failed: %w", err)
		}

		slog.Info("Rollback completed successfully")
		return nil
	},
}

// migrateVersionCmd shows current migration version
var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	Long:  `Display the current database migration version and status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		databaseURL := cfg.Store.ConnectionString()
		if databaseURL == "" {
			return fmt.Errorf("database URL not configured")
		}

		version, dirty, err := migrations.Version(databaseURL)
		if err != nil {
			return fmt.Errorf("failed to get version: %w", err)
		}

		if version == 0 {
			slog.Info("No migrations applied yet")
		} else {
			status := "clean"
			if dirty {
				status = "dirty"
			}
			slog.Info("Current migration version",
				slog.Uint64("version", uint64(version)),
				slog.String("status", status))
		}

		return nil
	},
}

func init() {
	migrateDownCmd.Flags().IntVarP(&migrateSteps, "steps", "n", 1, "Number of migrations to rollback")

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateVersionCmd)

	rootCmd.AddCommand(migrateCmd)
}
