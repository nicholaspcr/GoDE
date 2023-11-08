package store

import (
	"context"

	"github.com/nicholaspcr/GoDE/internal/store/gorm"
)

// Store contains the methods to interact with the database
type Store interface {
	UserStore
}

// New returns a new Store instance
func New(ctx context.Context) (Store, error) {
	st, err := gorm.New(ctx)
	if err != nil {
		return nil, err
	}

	if err := st.AutoMigrate(); err != nil {
		return nil, err
	}
	return st, nil
}
