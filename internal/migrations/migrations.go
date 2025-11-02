// Package migrations handles database schema migrations.
package migrations

import (
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Run executes database migrations.
func Run(databaseURL string) error {
	slog.Info("Running database migrations", slog.String("url", maskPassword(databaseURL)))

	// Create migrate instance with file source
	m, err := migrate.New("file://db/migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

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

	m, err := migrate.New("file://db/migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

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
	m, err := migrate.New("file://db/migrations", databaseURL)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

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
