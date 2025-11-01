package store

import (
	"context"
	"errors"

	"github.com/glebarez/sqlite"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nicholaspcr/GoDE/internal/store/gorm"
	"gorm.io/driver/postgres"
)

// New returns a new Store instance
func New(ctx context.Context, cfg Config) (Store, error) {
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

	if err := st.AutoMigrate(); err != nil {
		return nil, err
	}

	return st, nil
}
