package store

import (
	"context"

	"github.com/nicholaspcr/GoDE/internal/store/gorm"
)

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
