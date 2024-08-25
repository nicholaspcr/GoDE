package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/glebarez/sqlite"
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
		sqlDB, err := sql.Open("pgx", "mydb_dsn")
		if err != nil {
			return nil, err
		}
		dialector = postgres.New(postgres.Config{Conn: sqlDB})
	default:
		return nil, errors.New("invalid store type")
	}

	st, err := gorm.New(ctx, dialector)
	if err != nil {
		return nil, err
	}

	if err := st.AutoMigrate(); err != nil {
		return nil, err
	}

	return st, nil
}
