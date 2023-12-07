package store

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api"
)

// Store contains the methods to interact with the database
type Store interface {
	UserOperations
}

// UserOperations is the interface for the user store.
type UserOperations interface {
	Create(context.Context, *api.User) error
	Get(context.Context, *api.UserIDs) (*api.User, error)
	Update(context.Context, *api.User) error
	Delete(context.Context, *api.UserIDs) error
}
