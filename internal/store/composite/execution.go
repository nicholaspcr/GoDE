// Package composite provides a composite store that combines Redis and database stores.
package composite

import (
	"context"

	"github.com/nicholaspcr/GoDE/internal/store"
)

// ExecutionStore combines Redis and database stores for execution operations.
type ExecutionStore struct {
	redis store.ExecutionOperations
	db    store.ExecutionOperations
}

// NewExecutionStore creates a new composite execution store.
func NewExecutionStore(redis, db store.ExecutionOperations) *ExecutionStore {
	return &ExecutionStore{
		redis: redis,
		db:    db,
	}
}

// CreateExecution creates an execution in both Redis and database.
func (s *ExecutionStore) CreateExecution(ctx context.Context, execution *store.Execution) error {
	// Store in database for persistence
	if err := s.db.CreateExecution(ctx, execution); err != nil {
		return err
	}

	// Store in Redis for fast access
	if err := s.redis.CreateExecution(ctx, execution); err != nil {
		// Database write succeeded but Redis failed - log but don't fail the request
		// The execution will still work, just slower queries
		return nil
	}

	return nil
}

// GetExecution tries Redis first, falls back to database.
func (s *ExecutionStore) GetExecution(ctx context.Context, executionID, userID string) (*store.Execution, error) {
	// Try Redis first (fast path)
	execution, err := s.redis.GetExecution(ctx, executionID, userID)
	if err == nil {
		return execution, nil
	}

	// Fall back to database
	execution, err = s.db.GetExecution(ctx, executionID, userID)
	if err != nil {
		return nil, err
	}

	// Re-populate Redis cache
	_ = s.redis.CreateExecution(ctx, execution)

	return execution, nil
}

// UpdateExecutionStatus updates status in both stores.
func (s *ExecutionStore) UpdateExecutionStatus(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
	// Update database first (source of truth)
	if err := s.db.UpdateExecutionStatus(ctx, executionID, status, errorMsg); err != nil {
		return err
	}

	// Update Redis cache (best effort)
	_ = s.redis.UpdateExecutionStatus(ctx, executionID, status, errorMsg)

	return nil
}

// UpdateExecutionResult updates the pareto ID in both stores.
func (s *ExecutionStore) UpdateExecutionResult(ctx context.Context, executionID string, paretoID uint64) error {
	// Update database first (source of truth)
	if err := s.db.UpdateExecutionResult(ctx, executionID, paretoID); err != nil {
		return err
	}

	// Update Redis cache (best effort)
	_ = s.redis.UpdateExecutionResult(ctx, executionID, paretoID)

	return nil
}

// ListExecutions queries the database (source of truth for listing).
func (s *ExecutionStore) ListExecutions(ctx context.Context, userID string, status *store.ExecutionStatus) ([]*store.Execution, error) {
	return s.db.ListExecutions(ctx, userID, status)
}

// DeleteExecution removes from both stores.
func (s *ExecutionStore) DeleteExecution(ctx context.Context, executionID, userID string) error {
	// Delete from database
	if err := s.db.DeleteExecution(ctx, executionID, userID); err != nil {
		return err
	}

	// Delete from Redis (best effort)
	_ = s.redis.DeleteExecution(ctx, executionID, userID)

	return nil
}

// SaveProgress delegates to Redis only.
func (s *ExecutionStore) SaveProgress(ctx context.Context, progress *store.ExecutionProgress) error {
	return s.redis.SaveProgress(ctx, progress)
}

// GetProgress delegates to Redis only.
func (s *ExecutionStore) GetProgress(ctx context.Context, executionID string) (*store.ExecutionProgress, error) {
	return s.redis.GetProgress(ctx, executionID)
}

// MarkExecutionForCancellation delegates to Redis only.
func (s *ExecutionStore) MarkExecutionForCancellation(ctx context.Context, executionID, userID string) error {
	return s.redis.MarkExecutionForCancellation(ctx, executionID, userID)
}

// IsExecutionCancelled delegates to Redis only.
func (s *ExecutionStore) IsExecutionCancelled(ctx context.Context, executionID string) (bool, error) {
	return s.redis.IsExecutionCancelled(ctx, executionID)
}

// Subscribe delegates to Redis for pub/sub functionality.
func (s *ExecutionStore) Subscribe(ctx context.Context, channel string) (<-chan []byte, error) {
	return s.redis.Subscribe(ctx, channel)
}
