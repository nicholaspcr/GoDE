// Package composite provides a composite store that combines Redis and database stores.
package composite

import (
	"context"

	"github.com/nicholaspcr/GoDE/internal/cache/redis"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
)

// Store implements store.Store by combining database and Redis stores.
type Store struct {
	db          store.Store
	redis       *redis.Client
	execStore   *ExecutionStore
}

// New creates a new composite store.
func New(db store.Store, redisClient *redis.Client, redisExecStore store.ExecutionOperations) *Store {
	return &Store{
		db:        db,
		redis:     redisClient,
		execStore: NewExecutionStore(redisExecStore, db),
	}
}

// User operations delegate to database
func (s *Store) CreateUser(ctx context.Context, user *api.User) error {
	return s.db.CreateUser(ctx, user)
}

func (s *Store) GetUser(ctx context.Context, userIDs *api.UserIDs) (*api.User, error) {
	return s.db.GetUser(ctx, userIDs)
}

func (s *Store) UpdateUser(ctx context.Context, user *api.User, fields ...string) error {
	return s.db.UpdateUser(ctx, user, fields...)
}

func (s *Store) DeleteUser(ctx context.Context, userIDs *api.UserIDs) error {
	return s.db.DeleteUser(ctx, userIDs)
}

// Pareto operations delegate to database
func (s *Store) CreatePareto(ctx context.Context, pareto *api.Pareto) error {
	return s.db.CreatePareto(ctx, pareto)
}

func (s *Store) GetPareto(ctx context.Context, paretoIDs *api.ParetoIDs) (*api.Pareto, error) {
	return s.db.GetPareto(ctx, paretoIDs)
}

func (s *Store) UpdatePareto(ctx context.Context, pareto *api.Pareto, fields ...string) error {
	return s.db.UpdatePareto(ctx, pareto, fields...)
}

func (s *Store) DeletePareto(ctx context.Context, paretoIDs *api.ParetoIDs) error {
	return s.db.DeletePareto(ctx, paretoIDs)
}

func (s *Store) ListParetos(ctx context.Context, userIDs *api.UserIDs) ([]*api.Pareto, error) {
	return s.db.ListParetos(ctx, userIDs)
}

// ParetoSet operations delegate to database
func (s *Store) CreateParetoSet(ctx context.Context, paretoSet *store.ParetoSet) error {
	return s.db.CreateParetoSet(ctx, paretoSet)
}

func (s *Store) GetParetoSetByID(ctx context.Context, id uint64) (*store.ParetoSet, error) {
	return s.db.GetParetoSetByID(ctx, id)
}

// Execution operations use composite execution store
func (s *Store) CreateExecution(ctx context.Context, execution *store.Execution) error {
	return s.execStore.CreateExecution(ctx, execution)
}

func (s *Store) GetExecution(ctx context.Context, executionID, userID string) (*store.Execution, error) {
	return s.execStore.GetExecution(ctx, executionID, userID)
}

func (s *Store) UpdateExecutionStatus(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
	return s.execStore.UpdateExecutionStatus(ctx, executionID, status, errorMsg)
}

func (s *Store) UpdateExecutionResult(ctx context.Context, executionID string, paretoID uint64) error {
	return s.execStore.UpdateExecutionResult(ctx, executionID, paretoID)
}

func (s *Store) ListExecutions(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
	return s.execStore.ListExecutions(ctx, userID, status, limit, offset)
}

func (s *Store) DeleteExecution(ctx context.Context, executionID, userID string) error {
	return s.execStore.DeleteExecution(ctx, executionID, userID)
}

func (s *Store) SaveProgress(ctx context.Context, progress *store.ExecutionProgress) error {
	return s.execStore.SaveProgress(ctx, progress)
}

func (s *Store) GetProgress(ctx context.Context, executionID string) (*store.ExecutionProgress, error) {
	return s.execStore.GetProgress(ctx, executionID)
}

func (s *Store) MarkExecutionForCancellation(ctx context.Context, executionID, userID string) error {
	return s.execStore.MarkExecutionForCancellation(ctx, executionID, userID)
}

func (s *Store) IsExecutionCancelled(ctx context.Context, executionID string) (bool, error) {
	return s.execStore.IsExecutionCancelled(ctx, executionID)
}

func (s *Store) Subscribe(ctx context.Context, channel string) (<-chan []byte, error) {
	return s.execStore.Subscribe(ctx, channel)
}

// HealthCheck checks both database and Redis health.
func (s *Store) HealthCheck(ctx context.Context) error {
	// Check database health
	if err := s.db.HealthCheck(ctx); err != nil {
		return err
	}

	// Check Redis health
	if err := s.redis.HealthCheck(ctx); err != nil {
		return err
	}

	return nil
}

// RedisClient returns the underlying Redis client for separate health checks.
func (s *Store) RedisClient() *redis.Client {
	return s.redis
}
