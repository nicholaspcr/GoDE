package store

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
)

// Store contains the methods to interact with the database
type Store interface {
	UserOperations
	TenantOperations
}

// UserOperations is the interface for the user store.
type UserOperations interface {
	CreateUser(context.Context, *api.User) error
	GetUser(context.Context, *api.UserIDs) (*api.User, error)
	UpdateUser(context.Context, *api.User, ...string) error
	DeleteUser(context.Context, *api.UserIDs) error
}

type TenantOperations interface {
	CreateTenant(context.Context, *api.Tenant) error
	GetTenant(context.Context, *api.TenantIDs) (*api.Tenant, error)
	DeleteTenant(context.Context, *api.TenantIDs) error
}
