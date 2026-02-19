package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/nicholaspcr/GoDE/internal/cache/redis"
	"github.com/nicholaspcr/GoDE/internal/store"
	storerrors "github.com/nicholaspcr/GoDE/internal/store/errors"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// executionJSON is a helper struct for JSON serialization that handles protobuf Config field.
type executionJSON struct {
	ID          string                `json:"id"`
	UserID      string                `json:"user_id"`
	Status      store.ExecutionStatus `json:"status"`
	ConfigJSON  string                `json:"config_json"` // protojson encoded DEConfig
	Algorithm   string                `json:"algorithm,omitempty"`
	Variant     string                `json:"variant,omitempty"`
	Problem     string                `json:"problem,omitempty"`
	ParetoID    *uint64               `json:"pareto_id,omitempty"`
	Error       string                `json:"error,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	CompletedAt *time.Time            `json:"completed_at,omitempty"`
}

func marshalExecution(exec *store.Execution) ([]byte, error) {
	configJSON := ""
	if exec.Config != nil {
		data, err := protojson.Marshal(exec.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal config: %w", err)
		}
		configJSON = string(data)
	}

	helper := executionJSON{
		ID:          exec.ID,
		UserID:      exec.UserID,
		Status:      exec.Status,
		ConfigJSON:  configJSON,
		Algorithm:   exec.Algorithm,
		Variant:     exec.Variant,
		Problem:     exec.Problem,
		ParetoID:    exec.ParetoID,
		Error:       exec.Error,
		CreatedAt:   exec.CreatedAt,
		UpdatedAt:   exec.UpdatedAt,
		CompletedAt: exec.CompletedAt,
	}

	return json.Marshal(helper)
}

func unmarshalExecution(data []byte) (*store.Execution, error) {
	var helper executionJSON
	if err := json.Unmarshal(data, &helper); err != nil {
		return nil, fmt.Errorf("failed to unmarshal execution: %w", err)
	}

	var config *api.DEConfig
	if helper.ConfigJSON != "" {
		config = &api.DEConfig{}
		if err := protojson.Unmarshal([]byte(helper.ConfigJSON), config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}
	}

	return &store.Execution{
		ID:          helper.ID,
		UserID:      helper.UserID,
		Status:      helper.Status,
		Config:      config,
		Algorithm:   helper.Algorithm,
		Variant:     helper.Variant,
		Problem:     helper.Problem,
		ParetoID:    helper.ParetoID,
		Error:       helper.Error,
		CreatedAt:   helper.CreatedAt,
		UpdatedAt:   helper.UpdatedAt,
		CompletedAt: helper.CompletedAt,
	}, nil
}

// ExecutionStore implements ExecutionOperations using Redis.
type ExecutionStore struct {
	client       redis.ClientInterface
	executionTTL time.Duration
	progressTTL  time.Duration
	// updateLocks serializes concurrent read-modify-write on the same execution key.
	// This prevents lost updates when multiple goroutines update the same execution.
	updateLocks sync.Map
}

// NewExecutionStore creates a new Redis-backed execution store.
func NewExecutionStore(client redis.ClientInterface, executionTTL, progressTTL time.Duration) *ExecutionStore {
	return &ExecutionStore{
		client:       client,
		executionTTL: executionTTL,
		progressTTL:  progressTTL,
	}
}

// getUpdateLock returns a per-execution-ID mutex from the sync.Map, creating one if needed.
func (s *ExecutionStore) getUpdateLock(executionID string) *sync.Mutex {
	mu, _ := s.updateLocks.LoadOrStore(executionID, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

func (s *ExecutionStore) executionKey(executionID string) string {
	return fmt.Sprintf("execution:%s", executionID)
}

func (s *ExecutionStore) progressKey(executionID string) string {
	return fmt.Sprintf("execution:%s:progress", executionID)
}

func (s *ExecutionStore) cancelKey(executionID string) string {
	return fmt.Sprintf("execution:%s:cancel", executionID)
}

func (s *ExecutionStore) userExecutionsKey(userID string) string {
	return fmt.Sprintf("user:%s:executions", userID)
}

// updateExecution is a helper that implements the read-modify-write pattern for execution updates.
// It retrieves an execution, applies the modifier function, updates the timestamp, and saves back to Redis.
// A per-key mutex serializes concurrent updates to the same execution to prevent lost updates.
func (s *ExecutionStore) updateExecution(ctx context.Context, executionID string, modifier func(*store.Execution) error) error {
	// Acquire per-key lock to serialize concurrent read-modify-write
	mu := s.getUpdateLock(executionID)
	mu.Lock()
	defer mu.Unlock()

	key := s.executionKey(executionID)

	// Get current execution
	data, err := s.client.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("%w: %s", storerrors.ErrExecutionNotFound, executionID)
	}

	execution, err := unmarshalExecution([]byte(data))
	if err != nil {
		return err
	}

	// Apply modifications
	if err := modifier(execution); err != nil {
		return err
	}

	// Update timestamp
	execution.UpdatedAt = time.Now()

	// Save updated execution
	updatedData, err := marshalExecution(execution)
	if err != nil {
		return fmt.Errorf("failed to marshal execution: %w", err)
	}

	if err := s.client.Set(ctx, key, updatedData, s.executionTTL); err != nil {
		return fmt.Errorf("failed to update execution: %w", err)
	}

	// Update in user's execution set
	userKey := s.userExecutionsKey(execution.UserID)
	if err := s.client.HSet(ctx, userKey, executionID, string(updatedData)); err != nil {
		return fmt.Errorf("failed to update execution in user set: %w", err)
	}

	return nil
}

// CreateExecution stores a new execution in Redis.
func (s *ExecutionStore) CreateExecution(ctx context.Context, execution *store.Execution) error {
	key := s.executionKey(execution.ID)

	data, err := marshalExecution(execution)
	if err != nil {
		return fmt.Errorf("failed to marshal execution: %w", err)
	}

	if err := s.client.Set(ctx, key, data, s.executionTTL); err != nil {
		return fmt.Errorf("failed to store execution: %w", err)
	}

	// Add to user's execution set
	userKey := s.userExecutionsKey(execution.UserID)
	if err := s.client.HSet(ctx, userKey, execution.ID, string(data)); err != nil {
		return fmt.Errorf("failed to add execution to user set: %w", err)
	}
	if err := s.client.Expire(ctx, userKey, s.executionTTL); err != nil {
		return fmt.Errorf("failed to set TTL on user executions: %w", err)
	}

	return nil
}

// GetExecution retrieves an execution from Redis.
func (s *ExecutionStore) GetExecution(ctx context.Context, executionID, userID string) (*store.Execution, error) {
	key := s.executionKey(executionID)

	data, err := s.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", storerrors.ErrExecutionNotFound, executionID)
	}

	execution, err := unmarshalExecution([]byte(data))
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if execution.UserID != userID {
		return nil, fmt.Errorf("execution does not belong to user")
	}

	return execution, nil
}

// UpdateExecutionStatus updates the status of an execution.
func (s *ExecutionStore) UpdateExecutionStatus(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
	return s.updateExecution(ctx, executionID, func(exec *store.Execution) error {
		exec.Status = status
		if errorMsg != "" {
			exec.Error = errorMsg
		}
		if status == store.ExecutionStatusCompleted || status == store.ExecutionStatusFailed || status == store.ExecutionStatusCancelled {
			now := time.Now()
			exec.CompletedAt = &now
		}
		return nil
	})
}

// UpdateExecutionResult updates the pareto ID for a completed execution.
func (s *ExecutionStore) UpdateExecutionResult(ctx context.Context, executionID string, paretoID uint64) error {
	return s.updateExecution(ctx, executionID, func(exec *store.Execution) error {
		exec.ParetoID = &paretoID
		return nil
	})
}


// DeleteExecution removes an execution from Redis.
// If userID is empty, ownership verification is skipped (used for cache invalidation).
func (s *ExecutionStore) DeleteExecution(ctx context.Context, executionID, userID string) error {
	// Get execution - skip ownership check if userID is empty (cache invalidation)
	key := s.executionKey(executionID)
	data, err := s.client.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("%w: %s", storerrors.ErrExecutionNotFound, executionID)
	}

	execution, err := unmarshalExecution([]byte(data))
	if err != nil {
		return err
	}

	// Verify ownership only if userID is provided
	if userID != "" && execution.UserID != userID {
		return fmt.Errorf("execution does not belong to user")
	}

	// Delete from Redis
	if err := s.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete execution: %w", err)
	}

	// Remove from user's execution set using atomic HDel
	userKey := s.userExecutionsKey(execution.UserID)
	if err := s.client.HDel(ctx, userKey, executionID); err != nil {
		slog.Error("failed to remove execution from user set",
			slog.String("user_id", execution.UserID),
			slog.String("execution_id", executionID),
			slog.Any("error", err),
		)
		// Continue cleanup even if this fails
	}

	// Delete progress
	progressKey := s.progressKey(executionID)
	if delErr := s.client.Delete(ctx, progressKey); delErr != nil {
		slog.Warn("failed to delete progress key",
			slog.String("execution_id", executionID),
			slog.Any("error", delErr),
		)
	}

	// Delete cancellation flag
	cancelKey := s.cancelKey(executionID)
	if delErr := s.client.Delete(ctx, cancelKey); delErr != nil {
		slog.Warn("failed to delete cancellation key",
			slog.String("execution_id", executionID),
			slog.Any("error", delErr),
		)
	}

	return nil
}
