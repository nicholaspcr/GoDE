package store

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
)

// Store contains the methods to interact with the database
type Store interface {
	UserOperations
	ParetoOperations
	ExecutionOperations
	HealthCheck(context.Context) error
}

// UserOperations is the interface for the user store.
type UserOperations interface {
	CreateUser(context.Context, *api.User) error
	GetUser(context.Context, *api.UserIDs) (*api.User, error)
	UpdateUser(context.Context, *api.User, ...string) error
	DeleteUser(context.Context, *api.UserIDs) error
}

// ParetoOperations is the interface for the pareto store.
type ParetoOperations interface {
	CreatePareto(context.Context, *api.Pareto) error
	GetPareto(context.Context, *api.ParetoIDs) (*api.Pareto, error)
	UpdatePareto(context.Context, *api.Pareto, ...string) error
	DeletePareto(context.Context, *api.ParetoIDs) error
	ListParetos(ctx context.Context, userIDs *api.UserIDs, limit, offset int) ([]*api.Pareto, int, error)
	CreateParetoSet(context.Context, *ParetoSet) error
	GetParetoSetByID(context.Context, uint64) (*ParetoSet, error)
}

// ExecutionOperations is the interface for the execution store.
type ExecutionOperations interface {
	CreateExecution(ctx context.Context, execution *Execution) error
	GetExecution(ctx context.Context, executionID, userID string) (*Execution, error)
	UpdateExecutionStatus(ctx context.Context, executionID string, status ExecutionStatus, errorMsg string) error
	UpdateExecutionResult(ctx context.Context, executionID string, paretoID uint64) error
	ListExecutions(ctx context.Context, userID string, status *ExecutionStatus, limit, offset int) ([]*Execution, int, error)
	DeleteExecution(ctx context.Context, executionID, userID string) error

	// Idempotency: returns existing executionID or ErrExecutionNotFound.
	GetExecutionByIdempotencyKey(ctx context.Context, userID, idempotencyKey string) (string, error)

	// Progress tracking
	SaveProgress(ctx context.Context, progress *ExecutionProgress) error
	GetProgress(ctx context.Context, executionID string) (*ExecutionProgress, error)

	// Cancellation support
	MarkExecutionForCancellation(ctx context.Context, executionID, userID string) error
	IsExecutionCancelled(ctx context.Context, executionID string) (bool, error)

	// Real-time updates
	Subscribe(ctx context.Context, channel string) (<-chan []byte, error)
}
