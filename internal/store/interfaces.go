package store

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api"
)

// UserStore is the interface for the user store.
type UserStore interface {
	Create(context.Context, *api.User) error
	Read(context.Context, *api.UserID) (*api.User, error)
	Update(context.Context, *api.User) error
	Delete(context.Context, *api.UserID) error
}
