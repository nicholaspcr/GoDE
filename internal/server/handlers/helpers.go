package handlers

import (
	"context"

	"github.com/nicholaspcr/GoDE/internal/server/middleware"
	"github.com/nicholaspcr/GoDE/internal/store"
	api "github.com/nicholaspcr/GoDE/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// usernameFromContext extracts the authenticated username from the context.
// Returns a gRPC Unauthenticated error if the username is not present.
func usernameFromContext(ctx context.Context) (string, error) {
	userID := middleware.UsernameFromContext(ctx)
	if userID == "" {
		return "", status.Error(codes.Unauthenticated, "user not authenticated")
	}
	return userID, nil
}

// authDB is the minimal store interface required by authHandler.
type authDB interface {
	CreateUser(context.Context, *api.User) error
	GetUser(context.Context, *api.UserIDs) (*api.User, error)
}

// userDB is the minimal store interface required by userHandler.
type userDB interface {
	CreateUser(context.Context, *api.User) error
	GetUser(context.Context, *api.UserIDs) (*api.User, error)
	UpdateUser(context.Context, *api.User, ...string) error
	DeleteUser(context.Context, *api.UserIDs) error
}

// paretoDB is the minimal store interface required by paretoHandler.
type paretoDB interface {
	GetPareto(context.Context, *api.ParetoIDs) (*api.Pareto, error)
	DeletePareto(context.Context, *api.ParetoIDs) error
	ListParetos(ctx context.Context, userIDs *api.UserIDs, limit, offset int) ([]*api.Pareto, int, error)
}

// deStore is the minimal store interface required by deHandler.
type deStore interface {
	ListExecutions(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error)
	DeleteExecution(ctx context.Context, executionID, userID string) error
	GetExecution(ctx context.Context, executionID, userID string) (*store.Execution, error)
	GetProgress(ctx context.Context, executionID string) (*store.ExecutionProgress, error)
	GetParetoSetByID(context.Context, uint64) (*store.ParetoSet, error)
	Subscribe(ctx context.Context, channel string) (<-chan []byte, error)
}
