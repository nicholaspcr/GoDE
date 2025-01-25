// Package sqlite implements the methods defined in the CLI state interface.
package sqlite

import (
	"context"
	"errors"

	"github.com/glebarez/sqlite"
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/state"
	"gorm.io/gorm"
)

type store struct {
	*authTokenStore
}

// New returns a new store that handles state operations for the CLI.
func New(ctx context.Context, cfg Config) (state.Operations, error) {
	var dialector gorm.Dialector

	switch cfg.Provider {
	case "memory":
		dialector = sqlite.Open(":memory:")
	case "file":
		dialector = sqlite.Open(cfg.Filepath)
	default:
		return nil, errors.New("invalid store type")
	}

	db, err := gorm.Open(dialector)
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&authTokenModel{}); err != nil {
		return nil, err
	}

	return &store{
		authTokenStore: &authTokenStore{DB: db},
	}, nil
}
