package gorm

import (
	"context"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// memoryEnabled defines if the database should be in memory or not, useful for
// testing and debugging. Can be enabled by in configuration.
//
// TODO: Make this an environment option and disabled by default.
var memoryEnabled = true

// gormStore is the main store for the application. It contains implementations
// of all the interfaces defined in the store package.
type gormStore struct {
	db *gorm.DB
	*userStore
}

// New returns a new GormStore.
func New(_ context.Context) (*gormStore, error) {
	sqlitePath := ".dev/sqlite.db"
	if memoryEnabled {
		sqlitePath = ":memory:"
	}
	db, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	store := &gormStore{
		db:        db,
		userStore: newUserStore(db),
	}

	return store, nil
}

func (s *gormStore) AutoMigrate() error {
	return s.db.AutoMigrate(
		&userModel{},
	)
}
