package store

import (
	"context"
	"errors"
	"log/slog"

	"github.com/glebarez/sqlite"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nicholaspcr/GoDE/internal/migrations"
	"github.com/nicholaspcr/GoDE/internal/store/gorm"
	"gorm.io/driver/postgres"
)

// New returns a new Store instance
func New(ctx context.Context, cfg Config) (Store, error) {
	// Run migrations first (except for memory stores)
	if cfg.Type != "memory" {
		connStr := cfg.ConnectionString()
		if connStr != "" {
			slog.Info("Running database migrations before connecting...")
			if err := migrations.Run(connStr); err != nil {
				slog.Error("Migration failed", slog.String("error", err.Error()))
				return nil, err
			}
		}
	}

	var dialector gorm.Dialector

	switch cfg.Type {
	case "memory":
		dialector = sqlite.Open(":memory:")
	case "sqlite":
		dialector = sqlite.Open(cfg.Sqlite.Filepath)
	case "postgres":
		dialector = postgres.Open(cfg.Postgresql.DNS)
	default:
		return nil, errors.New("invalid store type")
	}

	st, err := gorm.New(dialector)
	if err != nil {
		return nil, err
	}

	// For memory stores, still use AutoMigrate since migrations don't work with :memory:
	if cfg.Type == "memory" {
		if err := st.AutoMigrate(); err != nil {
			return nil, err
		}
	}

	return st, nil
}
