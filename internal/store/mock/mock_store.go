// Package mock provides mock implementations of storage interfaces for testing.
package mock

import (
	"context"

	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
)

// MockStore is a mock implementation of store.Store for testing.
type MockStore struct {
	// User operations
	CreateUserFn func(ctx context.Context, user *api.User) error
	GetUserFn    func(ctx context.Context, userIDs *api.UserIDs) (*api.User, error)
	UpdateUserFn func(ctx context.Context, user *api.User, fields ...string) error
	DeleteUserFn func(ctx context.Context, userIDs *api.UserIDs) error

	// Pareto operations
	CreateParetoFn func(ctx context.Context, pareto *api.Pareto) error
	GetParetoFn    func(ctx context.Context, ids *api.ParetoIDs) (*api.Pareto, error)
	UpdateParetoFn func(ctx context.Context, pareto *api.Pareto, fields ...string) error
	DeleteParetoFn func(ctx context.Context, ids *api.ParetoIDs) error
	ListParetosFn  func(ctx context.Context, userIDs *api.UserIDs) ([]*api.Pareto, error)

	// Vector operations
	CreateVectorFn func(ctx context.Context, vector *api.Vector, paretoID int64) error
	GetVectorFn    func(ctx context.Context, id int64) (*api.Vector, error)
	UpdateVectorFn func(ctx context.Context, vector *api.Vector) error
	DeleteVectorFn func(ctx context.Context, id int64) error

	AutoMigrateFn func() error
	HealthCheckFn func(ctx context.Context) error
}

// Verify MockStore implements store.Store
var _ store.Store = (*MockStore)(nil)

// CreateUser implements store.Store
func (m *MockStore) CreateUser(ctx context.Context, user *api.User) error {
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, user)
	}
	return nil
}

// GetUser implements store.Store
func (m *MockStore) GetUser(ctx context.Context, userIDs *api.UserIDs) (*api.User, error) {
	if m.GetUserFn != nil {
		return m.GetUserFn(ctx, userIDs)
	}
	return nil, nil
}

// UpdateUser implements store.Store
func (m *MockStore) UpdateUser(ctx context.Context, user *api.User, fields ...string) error {
	if m.UpdateUserFn != nil {
		return m.UpdateUserFn(ctx, user, fields...)
	}
	return nil
}

// DeleteUser implements store.Store
func (m *MockStore) DeleteUser(ctx context.Context, userIDs *api.UserIDs) error {
	if m.DeleteUserFn != nil {
		return m.DeleteUserFn(ctx, userIDs)
	}
	return nil
}

// CreatePareto implements store.Store
func (m *MockStore) CreatePareto(ctx context.Context, pareto *api.Pareto) error {
	if m.CreateParetoFn != nil {
		return m.CreateParetoFn(ctx, pareto)
	}
	return nil
}

// GetPareto implements store.Store
func (m *MockStore) GetPareto(ctx context.Context, ids *api.ParetoIDs) (*api.Pareto, error) {
	if m.GetParetoFn != nil {
		return m.GetParetoFn(ctx, ids)
	}
	return nil, nil
}

// UpdatePareto implements store.Store
func (m *MockStore) UpdatePareto(ctx context.Context, pareto *api.Pareto, fields ...string) error {
	if m.UpdateParetoFn != nil {
		return m.UpdateParetoFn(ctx, pareto, fields...)
	}
	return nil
}

// DeletePareto implements store.Store
func (m *MockStore) DeletePareto(ctx context.Context, ids *api.ParetoIDs) error {
	if m.DeleteParetoFn != nil {
		return m.DeleteParetoFn(ctx, ids)
	}
	return nil
}

// ListParetos implements store.Store
func (m *MockStore) ListParetos(ctx context.Context, userIDs *api.UserIDs) ([]*api.Pareto, error) {
	if m.ListParetosFn != nil {
		return m.ListParetosFn(ctx, userIDs)
	}
	return nil, nil
}

// CreateVector implements store.Store
func (m *MockStore) CreateVector(ctx context.Context, vector *api.Vector, paretoID int64) error {
	if m.CreateVectorFn != nil {
		return m.CreateVectorFn(ctx, vector, paretoID)
	}
	return nil
}

// GetVector implements store.Store
func (m *MockStore) GetVector(ctx context.Context, id int64) (*api.Vector, error) {
	if m.GetVectorFn != nil {
		return m.GetVectorFn(ctx, id)
	}
	return nil, nil
}

// UpdateVector implements store.Store
func (m *MockStore) UpdateVector(ctx context.Context, vector *api.Vector) error {
	if m.UpdateVectorFn != nil {
		return m.UpdateVectorFn(ctx, vector)
	}
	return nil
}

// DeleteVector implements store.Store
func (m *MockStore) DeleteVector(ctx context.Context, id int64) error {
	if m.DeleteVectorFn != nil {
		return m.DeleteVectorFn(ctx, id)
	}
	return nil
}

// AutoMigrate implements store.Store
func (m *MockStore) AutoMigrate() error {
	if m.AutoMigrateFn != nil {
		return m.AutoMigrateFn()
	}
	return nil
}

// HealthCheck implements store.Store
func (m *MockStore) HealthCheck(ctx context.Context) error {
	if m.HealthCheckFn != nil {
		return m.HealthCheckFn(ctx)
	}
	return nil
}
