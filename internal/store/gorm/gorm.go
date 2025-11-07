// Package gorm contains the gorm implementation of the interfaces defined in
// the store package.
package gorm

import (
	"time"

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

	// Configure connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool parameters for optimal performance and resource management
	sqlDB.SetMaxIdleConns(10)                  // Maximum number of idle connections in the pool
	sqlDB.SetMaxOpenConns(100)                 // Maximum number of open connections to the database
	sqlDB.SetConnMaxLifetime(time.Hour)        // Maximum amount of time a connection may be reused
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Maximum amount of time a connection may be idle

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
