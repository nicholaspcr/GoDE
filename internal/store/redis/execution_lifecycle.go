package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
)

// SaveProgress stores execution progress in Redis and publishes update via pub/sub.
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

	// Handle nil pubsub (e.g., during testing with mocks)
	if pubsub == nil {
		go func() {
			<-ctx.Done()
			close(ch)
		}()
		return ch, nil
	}

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
func ConfigToProto(config *api.DEConfig) map[string]any {
	result := map[string]any{
		"executions":      config.Executions,
		"generations":     config.Generations,
		"population_size": config.PopulationSize,
		"dimensions_size": config.DimensionsSize,
		"objectives_size": config.ObjectivesSize,
		"floor_limiter":   config.FloorLimiter,
		"ceil_limiter":    config.CeilLimiter,
	}

	if config.AlgorithmConfig != nil {
		if algConfig, ok := config.AlgorithmConfig.(*api.DEConfig_Gde3); ok && algConfig.Gde3 != nil {
			result["gde3_config"] = map[string]any{
				"cr": algConfig.Gde3.Cr,
				"f":  algConfig.Gde3.F,
				"p":  algConfig.Gde3.P,
			}
		}
	}

	return result
}
