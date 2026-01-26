package redis

import (
	"context"
	"fmt"

	"github.com/nicholaspcr/GoDE/internal/store"
)

// ListExecutions returns a paginated list of executions for a user, optionally filtered by status.
// Uses HSCAN for memory-efficient iteration and supports early termination optimization.
func (s *ExecutionStore) ListExecutions(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
	// Apply defaults first
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	userKey := s.userExecutionsKey(userID)

	// Optimization: when not filtering by status, get total count upfront using HLEN
	// This allows early termination after collecting enough items
	var totalCount int
	earlyTerminationEnabled := (status == nil)

	if earlyTerminationEnabled {
		count, err := s.client.HLen(ctx, userKey)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get execution count: %w", err)
		}
		totalCount = int(count)
	}

	// Use HSCAN to iterate without loading all into memory
	var cursor uint64
	var executions []*store.Execution
	seen := 0
	collected := 0

	for {
		// HScan returns alternating field/value pairs
		pairs, nextCursor, err := s.client.HScan(ctx, userKey, cursor, "*", 100)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user executions: %w", err)
		}

		// Process pairs (alternating key, value)
		for i := 1; i < len(pairs); i += 2 {
			executionData := pairs[i]

			execution, err := unmarshalExecution([]byte(executionData))
			if err != nil {
				continue // Skip invalid entries
			}

			// Filter by status if provided
			if status != nil && execution.Status != *status {
				continue
			}

			seen++
			// Only collect items within our pagination window
			if seen > offset && collected < limit {
				executions = append(executions, execution)
				collected++

				// Early termination: if we've collected enough and we're not filtering,
				// we can stop scanning since we already know the total count
				if earlyTerminationEnabled && collected >= limit {
					return executions, totalCount, nil
				}
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	// When filtering by status, return the filtered count
	if !earlyTerminationEnabled {
		totalCount = seen
	}

	return executions, totalCount, nil
}
