package store

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
)

// Store contains the methods to interact with the database
type Store interface {
	UserOperations
	ParetoOperations
	HealthCheck(context.Context) error
}

// UserOperations is the interface for the user store.
type UserOperations interface {
	CreateUser(context.Context, *api.User) error
	GetUser(context.Context, *api.UserIDs) (*api.User, error)
	UpdateUser(context.Context, *api.User, ...string) error
	DeleteUser(context.Context, *api.UserIDs) error
}

// TenantOperations is the interface for the tenant store.
type TenantOperations interface {
	CreateTenant(context.Context, *api.Tenant) error
	GetTenant(context.Context, *api.TenantIDs) (*api.Tenant, error)
	DeleteTenant(context.Context, *api.TenantIDs) error
}

// ParetoOperations is the interface for the pareto store.
type ParetoOperations interface {
	CreatePareto(context.Context, *api.Pareto) error
	GetPareto(context.Context, *api.ParetoIDs) (*api.Pareto, error)
	UpdatePareto(context.Context, *api.Pareto, ...string) error
	DeletePareto(context.Context, *api.ParetoIDs) error
	ListParetos(context.Context, *api.UserIDs) ([]*api.Pareto, error)
}
