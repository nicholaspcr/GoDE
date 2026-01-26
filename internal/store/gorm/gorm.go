// Package gorm contains the gorm implementation of the interfaces defined in
// the store package.
package gorm

import (
	"context"
	"time"

	"github.com/nicholaspcr/GoDE/internal/store"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

// Dialector is a type alias for gorm.Dialector used for database connections.
type Dialector gorm.Dialector

// gormStore is the main store for the application. It contains implementations
// of all the interfaces defined in the store package.
type gormStore struct {
	db *gorm.DB
	*userStore
	*paretoStore
	*vectorStore
	*executionStore
}

// New returns a new GormStore.
func New(dialector gorm.Dialector, pool store.ConnectionPool) (*gormStore, error) {
	db, err := gorm.Open(dialector)
	if err != nil {
		return nil, err
	}

	// Configure connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool parameters from configuration
	sqlDB.SetMaxIdleConns(pool.MaxIdleConns)
	sqlDB.SetMaxOpenConns(pool.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(pool.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Keep hardcoded for now

	if err := db.Use(tracing.NewPlugin()); err != nil {
		return nil, err
	}

	store := &gormStore{
		db:             db,
		userStore:      newUserStore(db),
		paretoStore:    newParetoStore(db),
		vectorStore:    newVectorStore(db),
		executionStore: newExecutionStore(db),
	}

	return store, nil
}

func (s *gormStore) AutoMigrate() error {
	return s.db.AutoMigrate(
		&userModel{},
		&paretoModel{},
		&vectorModel{},
		&executionModel{},
	)
}

// HealthCheck verifies the database connection is alive by pinging it.
func (s *gormStore) HealthCheck(ctx context.Context) error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
