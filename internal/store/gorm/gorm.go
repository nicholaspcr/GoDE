package gorm

import (
	"context"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// TODO: Make this an option and disabled by default.

// MemoryEnabled defines if the database should be in memory or not, useful for
// testing and debugging. Can be enabled by in configuration.
var MemoryEnabled = false

// GormStore is the main store for the application. It contains implementations
// of all the interfaces defined in the store package.
type GormStore struct {
	db *gorm.DB
	*userStore
}

// New returns a new GormStore.
func New(ctx context.Context) (*GormStore, error) {
	var sqlitePath string = ".dev/sqlite.db"
	if MemoryEnabled {
		sqlitePath = ":memory:"
	}
	db, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	store := &GormStore{
		db:        db,
		userStore: newUserStore(db),
	}

	return store, nil
}

func (s *GormStore) AutoMigrate() error {
	return s.db.AutoMigrate(&userModel{})
}
