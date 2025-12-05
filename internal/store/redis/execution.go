package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nicholaspcr/GoDE/internal/cache/redis"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
)

// ExecutionStore implements ExecutionOperations using Redis.
type ExecutionStore struct {
	client       *redis.Client
	executionTTL time.Duration
	progressTTL  time.Duration
}

// NewExecutionStore creates a new Redis-backed execution store.
func NewExecutionStore(client *redis.Client, executionTTL, progressTTL time.Duration) *ExecutionStore {
	return &ExecutionStore{
		client:       client,
		executionTTL: executionTTL,
		progressTTL:  progressTTL,
	}
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

// CreateExecution stores a new execution in Redis.
func (s *ExecutionStore) CreateExecution(ctx context.Context, execution *store.Execution) error {
	key := s.executionKey(execution.ID)

	data, err := json.Marshal(execution)
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
		return nil, fmt.Errorf("execution not found: %w", err)
	}

	var execution store.Execution
	if err := json.Unmarshal([]byte(data), &execution); err != nil {
		return nil, fmt.Errorf("failed to unmarshal execution: %w", err)
	}

	// Verify ownership
	if execution.UserID != userID {
		return nil, fmt.Errorf("execution does not belong to user")
	}

	return &execution, nil
}

// UpdateExecutionStatus updates the status of an execution.
func (s *ExecutionStore) UpdateExecutionStatus(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
	key := s.executionKey(executionID)

	data, err := s.client.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("execution not found: %w", err)
	}

	var execution store.Execution
	if err := json.Unmarshal([]byte(data), &execution); err != nil {
		return fmt.Errorf("failed to unmarshal execution: %w", err)
	}

	execution.Status = status
	execution.UpdatedAt = time.Now()
	if errorMsg != "" {
		execution.Error = errorMsg
	}
	if status == store.ExecutionStatusCompleted || status == store.ExecutionStatusFailed || status == store.ExecutionStatusCancelled {
		now := time.Now()
		execution.CompletedAt = &now
	}

	updatedData, err := json.Marshal(execution)
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

// UpdateExecutionResult updates the pareto ID for a completed execution.
func (s *ExecutionStore) UpdateExecutionResult(ctx context.Context, executionID string, paretoID uint64) error {
	key := s.executionKey(executionID)

	data, err := s.client.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("execution not found: %w", err)
	}

	var execution store.Execution
	if err := json.Unmarshal([]byte(data), &execution); err != nil {
		return fmt.Errorf("failed to unmarshal execution: %w", err)
	}

	execution.ParetoID = &paretoID
	execution.UpdatedAt = time.Now()

	updatedData, err := json.Marshal(execution)
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

// ListExecutions retrieves all executions for a user, optionally filtered by status.
func (s *ExecutionStore) ListExecutions(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
	userKey := s.userExecutionsKey(userID)

	data, err := s.client.HGetAll(ctx, userKey)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user executions: %w", err)
	}

	// First, collect all matching executions
	allExecutions := make([]*store.Execution, 0, len(data))
	for _, executionData := range data {
		var execution store.Execution
		if err := json.Unmarshal([]byte(executionData), &execution); err != nil {
			continue // Skip invalid entries
		}

		// Filter by status if provided
		if status != nil && execution.Status != *status {
			continue
		}

		allExecutions = append(allExecutions, &execution)
	}

	// Apply defaults
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	totalCount := len(allExecutions)

	// Apply in-memory pagination
	start := offset
	if start > totalCount {
		return []*store.Execution{}, totalCount, nil
	}

	end := start + limit
	if end > totalCount {
		end = totalCount
	}

	return allExecutions[start:end], totalCount, nil
}

// DeleteExecution removes an execution from Redis.
func (s *ExecutionStore) DeleteExecution(ctx context.Context, executionID, userID string) error {
	// Verify ownership first
	execution, err := s.GetExecution(ctx, executionID, userID)
	if err != nil {
		return err
	}

	// Delete from Redis
	key := s.executionKey(executionID)
	if err := s.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("failed to delete execution: %w", err)
	}

	// Remove from user's execution set
	userKey := s.userExecutionsKey(execution.UserID)
	fields, err := s.client.HGetAll(ctx, userKey)
	if err == nil {
		delete(fields, executionID)
		// Clear and repopulate
		if err := s.client.Delete(ctx, userKey); err == nil {
			for id, data := range fields {
				if hsetErr := s.client.HSet(ctx, userKey, id, data); hsetErr != nil {
					slog.Error("failed to repopulate user execution set",
						slog.String("user_id", execution.UserID),
						slog.String("execution_id", id),
						slog.Any("error", hsetErr),
					)
				}
			}
		}
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

// SaveProgress stores execution progress in Redis.
func (s *ExecutionStore) SaveProgress(ctx context.Context, progress *store.ExecutionProgress) error {
	key := s.progressKey(progress.ExecutionID)

	progress.UpdatedAt = time.Now()

	data, err := json.Marshal(progress)
	if err != nil {
		return fmt.Errorf("failed to marshal progress: %w", err)
	}

	if err := s.client.Set(ctx, key, data, s.progressTTL); err != nil {
		return fmt.Errorf("failed to store progress: %w", err)
	}

	// Publish progress update to pub/sub channel
	channel := fmt.Sprintf("execution:%s:updates", progress.ExecutionID)
	if pubErr := s.client.Publish(ctx, channel, data); pubErr != nil {
		slog.Warn("failed to publish progress update",
			slog.String("execution_id", progress.ExecutionID),
			slog.String("channel", channel),
			slog.Any("error", pubErr),
		)
	}

	return nil
}

// GetProgress retrieves the current progress of an execution.
func (s *ExecutionStore) GetProgress(ctx context.Context, executionID string) (*store.ExecutionProgress, error) {
	key := s.progressKey(executionID)

	data, err := s.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("progress not found: %w", err)
	}

	var progress store.ExecutionProgress
	if err := json.Unmarshal([]byte(data), &progress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal progress: %w", err)
	}

	return &progress, nil
}

// MarkExecutionForCancellation sets a cancellation flag for an execution.
func (s *ExecutionStore) MarkExecutionForCancellation(ctx context.Context, executionID, userID string) error {
	// Verify ownership
	if _, err := s.GetExecution(ctx, executionID, userID); err != nil {
		return err
	}

	key := s.cancelKey(executionID)
	if err := s.client.Set(ctx, key, "1", s.executionTTL); err != nil {
		return fmt.Errorf("failed to set cancellation flag: %w", err)
	}

	// Publish cancellation event
	channel := fmt.Sprintf("execution:%s:cancel", executionID)
	if pubErr := s.client.Publish(ctx, channel, "cancel"); pubErr != nil {
		slog.Warn("failed to publish cancellation event",
			slog.String("execution_id", executionID),
			slog.String("channel", channel),
			slog.Any("error", pubErr),
		)
	}

	return nil
}

// IsExecutionCancelled checks if an execution has been marked for cancellation.
func (s *ExecutionStore) IsExecutionCancelled(ctx context.Context, executionID string) (bool, error) {
	key := s.cancelKey(executionID)

	_, err := s.client.Get(ctx, key)
	if err != nil {
		// Key doesn't exist, not cancelled
		return false, nil
	}

	return true, nil
}

// Subscribe subscribes to real-time updates on a channel.
func (s *ExecutionStore) Subscribe(ctx context.Context, channel string) (<-chan []byte, error) {
	pubsub := s.client.Subscribe(ctx, channel)

	// Create output channel
	ch := make(chan []byte, 100) // Buffer to prevent blocking

	// Start goroutine to receive messages and forward to channel
	go func() {
		defer close(ch)
		defer func() {
			if closeErr := pubsub.Close(); closeErr != nil {
				slog.Warn("failed to close pubsub connection",
					slog.String("channel", channel),
					slog.Any("error", closeErr),
				)
			}
		}()

		// Subscribe to the PubSub channel
		msgChan := pubsub.Channel()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgChan:
				if !ok {
					return
				}
				if msg != nil {
					ch <- []byte(msg.Payload)
				}
			}
		}
	}()

	return ch, nil
}

// ConfigToProto converts a DEConfig proto to a store-compatible format for JSON marshaling.
func ConfigToProto(config *api.DEConfig) map[string]interface{} {
	result := map[string]interface{}{
		"executions":      config.Executions,
		"generations":     config.Generations,
		"population_size": config.PopulationSize,
		"dimensions_size": config.DimensionsSize,
		"objectives_size":  config.ObjectivesSize,
		"floor_limiter":   config.FloorLimiter,
		"ceil_limiter":    config.CeilLimiter,
	}

	if config.AlgorithmConfig != nil {
		if algConfig, ok := config.AlgorithmConfig.(*api.DEConfig_Gde3); ok && algConfig.Gde3 != nil {
			result["gde3_config"] = map[string]interface{}{
				"cr": algConfig.Gde3.Cr,
				"f":  algConfig.Gde3.F,
				"p":  algConfig.Gde3.P,
			}
		}
	}

	return result
}
