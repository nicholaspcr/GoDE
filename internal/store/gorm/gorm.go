// Package gorm contains the gorm implementation of the interfaces defined in
// the store package.
package gorm

import (
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type Dialector gorm.Dialector

// gormStore is the main store for the application. It contains implementations
// of all the interfaces defined in the store package.
type gormStore struct {
	db *gorm.DB
	*userStore
	*paretoStore
	*vectorStore
}

// New returns a new GormStore.
func New(dialector gorm.Dialector) (*gormStore, error) {
	db, err := gorm.Open(dialector)
	if err != nil {
		return nil, err
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		return nil, err
	}

	store := &gormStore{
		db:          db,
		userStore:   newUserStore(db),
		paretoStore: newParetoStore(db),
		vectorStore: newVectorStore(db),
	}

	return store, nil
}

func (s *gormStore) AutoMigrate() error {
	return s.db.AutoMigrate(
		&userModel{},
		&paretoModel{},
		&vectorModel{},
	)
}
