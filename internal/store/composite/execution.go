// Package composite provides a composite store that combines Redis and database stores.
package composite

import (
	"context"
	"log/slog"

	"github.com/nicholaspcr/GoDE/internal/store"
)

// ExecutionStore combines Redis and database stores for execution operations.
type ExecutionStore struct {
	redis  store.ExecutionOperations
	db     store.ExecutionOperations
	logger *slog.Logger
}

// NewExecutionStore creates a new composite execution store.
func NewExecutionStore(redis, db store.ExecutionOperations) *ExecutionStore {
	return &ExecutionStore{
		redis:  redis,
		db:     db,
		logger: slog.Default(),
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
		s.logger.Warn("failed to populate cache on create",
			slog.String("execution_id", execution.ID),
			slog.Any("error", err))
		// Database write succeeded but Redis failed - log but don't fail the request
		// The execution will still work, just slower queries
		return nil
	}

	return nil
}

// GetExecution tries Redis first, validates freshness against database.
func (s *ExecutionStore) GetExecution(ctx context.Context, executionID, userID string) (*store.Execution, error) {
	// Try Redis first
	cachedExec, cacheErr := s.redis.GetExecution(ctx, executionID, userID)

	// Fetch from DB for comparison
	dbExec, dbErr := s.db.GetExecution(ctx, executionID, userID)
	if dbErr != nil {
		if cacheErr == nil {
			s.logger.Warn("database unavailable, using cached execution",
				slog.String("execution_id", executionID))
			return cachedExec, nil
		}
		return nil, dbErr
	}

	// If cache is fresh, use it
	if cacheErr == nil && !cachedExec.UpdatedAt.Before(dbExec.UpdatedAt) {
		return cachedExec, nil
	}

	// Refresh stale cache
	_ = s.redis.CreateExecution(ctx, dbExec)
	return dbExec, nil
}

// UpdateExecutionStatus updates status in database and invalidates cache.
func (s *ExecutionStore) UpdateExecutionStatus(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
	// Update database first (source of truth)
	if err := s.db.UpdateExecutionStatus(ctx, executionID, status, errorMsg); err != nil {
		return err
	}

	// Invalidate cache instead of updating to avoid stale data
	if err := s.redis.DeleteExecution(ctx, executionID, ""); err != nil {
		s.logger.Warn("failed to invalidate cache on status update",
			slog.String("execution_id", executionID),
			slog.Any("error", err))
	}

	return nil
}

// UpdateExecutionResult updates the pareto ID in database and invalidates cache.
func (s *ExecutionStore) UpdateExecutionResult(ctx context.Context, executionID string, paretoID uint64) error {
	// Update database first (source of truth)
	if err := s.db.UpdateExecutionResult(ctx, executionID, paretoID); err != nil {
		return err
	}

	// Invalidate cache instead of updating to avoid stale data
	if err := s.redis.DeleteExecution(ctx, executionID, ""); err != nil {
		s.logger.Warn("failed to invalidate cache on result update",
			slog.String("execution_id", executionID),
			slog.Any("error", err))
	}

	return nil
}

// ListExecutions queries the database (source of truth for listing).
func (s *ExecutionStore) ListExecutions(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
	return s.db.ListExecutions(ctx, userID, status, limit, offset)
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
