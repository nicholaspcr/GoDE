// Package migrations handles database schema migrations.
package migrations

import (
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	// postgres driver for migrate
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// sqlite3 driver for migrate
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/nicholaspcr/GoDE/internal/store/migrations"
)

// Run executes database migrations.
func Run(databaseURL string) error {
	slog.Info("Running database migrations", slog.String("url", maskPassword(databaseURL)))

	// Create migrate instance with embedded filesystem source
	sourceDriver, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to create iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			slog.Error("failed to close migrate source", slog.String("error", srcErr.Error()))
		}
		if dbErr != nil {
			slog.Error("failed to close migrate database", slog.String("error", dbErr.Error()))
		}
	}()

	// Run migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			slog.Info("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if dirty {
		slog.Warn("Database is in dirty state", slog.Uint64("version", uint64(version)))
		return fmt.Errorf("database is in dirty state at version %d", version)
	}

	slog.Info("Migrations applied successfully", slog.Uint64("version", uint64(version)))
	return nil
}

// Rollback rolls back the last migration.
func Rollback(databaseURL string, steps int) error {
	slog.Info("Rolling back database migrations",
		slog.String("url", maskPassword(databaseURL)),
		slog.Int("steps", steps))

	// Create migrate instance with embedded filesystem source
	sourceDriver, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("failed to create iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			slog.Error("failed to close migrate source", slog.String("error", srcErr.Error()))
		}
		if dbErr != nil {
			slog.Error("failed to close migrate database", slog.String("error", dbErr.Error()))
		}
	}()

	if err := m.Steps(-steps); err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			slog.Info("Rolled back all migrations")
			return nil
		}
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if dirty {
		slog.Warn("Database is in dirty state after rollback", slog.Uint64("version", uint64(version)))
		return fmt.Errorf("database is in dirty state at version %d", version)
	}

	slog.Info("Migrations rolled back successfully", slog.Uint64("version", uint64(version)))
	return nil
}

// Version returns the current migration version.
func Version(databaseURL string) (uint, bool, error) {
	// Create migrate instance with embedded filesystem source
	sourceDriver, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return 0, false, fmt.Errorf("failed to create iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, databaseURL)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			slog.Error("failed to close migrate source", slog.String("error", srcErr.Error()))
		}
		if dbErr != nil {
			slog.Error("failed to close migrate database", slog.String("error", dbErr.Error()))
		}
	}()

	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

// maskPassword masks the password in database URL for logging.
func maskPassword(url string) string {
	// Simple masking for security - in production you'd want more sophisticated masking
	if len(url) > 20 {
		return url[:10] + "***" + url[len(url)-10:]
	}
	return "***"
}
