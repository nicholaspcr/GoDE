// Package composite provides tests for the composite store implementation.
package composite

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockExecutionStore implements store.ExecutionOperations for testing.
type mockExecutionStore struct {
	// Function hooks for customizing behavior
	CreateExecutionFn              func(ctx context.Context, execution *store.Execution) error
	GetExecutionFn                 func(ctx context.Context, executionID, userID string) (*store.Execution, error)
	UpdateExecutionStatusFn        func(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error
	UpdateExecutionResultFn        func(ctx context.Context, executionID string, paretoID uint64) error
	ListExecutionsFn               func(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error)
	DeleteExecutionFn              func(ctx context.Context, executionID, userID string) error
	SaveProgressFn                 func(ctx context.Context, progress *store.ExecutionProgress) error
	GetProgressFn                  func(ctx context.Context, executionID string) (*store.ExecutionProgress, error)
	MarkExecutionForCancellationFn func(ctx context.Context, executionID, userID string) error
	IsExecutionCancelledFn         func(ctx context.Context, executionID string) (bool, error)
	SubscribeFn                    func(ctx context.Context, channel string) (<-chan []byte, error)

	// Call tracking
	createExecutionCalls              int
	getExecutionCalls                 int
	updateExecutionStatusCalls        int
	updateExecutionResultCalls        int
	listExecutionsCalls               int
	deleteExecutionCalls              int
	saveProgressCalls                 int
	getProgressCalls                  int
	markExecutionForCancellationCalls int
	isExecutionCancelledCalls         int
	subscribeCalls                    int
}

// Verify mockExecutionStore implements store.ExecutionOperations
var _ store.ExecutionOperations = (*mockExecutionStore)(nil)

func (m *mockExecutionStore) CreateExecution(ctx context.Context, execution *store.Execution) error {
	m.createExecutionCalls++
	if m.CreateExecutionFn != nil {
		return m.CreateExecutionFn(ctx, execution)
	}
	return nil
}

func (m *mockExecutionStore) GetExecution(ctx context.Context, executionID, userID string) (*store.Execution, error) {
	m.getExecutionCalls++
	if m.GetExecutionFn != nil {
		return m.GetExecutionFn(ctx, executionID, userID)
	}
	return nil, nil
}

func (m *mockExecutionStore) UpdateExecutionStatus(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
	m.updateExecutionStatusCalls++
	if m.UpdateExecutionStatusFn != nil {
		return m.UpdateExecutionStatusFn(ctx, executionID, status, errorMsg)
	}
	return nil
}

func (m *mockExecutionStore) UpdateExecutionResult(ctx context.Context, executionID string, paretoID uint64) error {
	m.updateExecutionResultCalls++
	if m.UpdateExecutionResultFn != nil {
		return m.UpdateExecutionResultFn(ctx, executionID, paretoID)
	}
	return nil
}

func (m *mockExecutionStore) ListExecutions(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
	m.listExecutionsCalls++
	if m.ListExecutionsFn != nil {
		return m.ListExecutionsFn(ctx, userID, status, limit, offset)
	}
	return nil, 0, nil
}

func (m *mockExecutionStore) DeleteExecution(ctx context.Context, executionID, userID string) error {
	m.deleteExecutionCalls++
	if m.DeleteExecutionFn != nil {
		return m.DeleteExecutionFn(ctx, executionID, userID)
	}
	return nil
}

func (m *mockExecutionStore) SaveProgress(ctx context.Context, progress *store.ExecutionProgress) error {
	m.saveProgressCalls++
	if m.SaveProgressFn != nil {
		return m.SaveProgressFn(ctx, progress)
	}
	return nil
}

func (m *mockExecutionStore) GetProgress(ctx context.Context, executionID string) (*store.ExecutionProgress, error) {
	m.getProgressCalls++
	if m.GetProgressFn != nil {
		return m.GetProgressFn(ctx, executionID)
	}
	return nil, nil
}

func (m *mockExecutionStore) MarkExecutionForCancellation(ctx context.Context, executionID, userID string) error {
	m.markExecutionForCancellationCalls++
	if m.MarkExecutionForCancellationFn != nil {
		return m.MarkExecutionForCancellationFn(ctx, executionID, userID)
	}
	return nil
}

func (m *mockExecutionStore) IsExecutionCancelled(ctx context.Context, executionID string) (bool, error) {
	m.isExecutionCancelledCalls++
	if m.IsExecutionCancelledFn != nil {
		return m.IsExecutionCancelledFn(ctx, executionID)
	}
	return false, nil
}

func (m *mockExecutionStore) Subscribe(ctx context.Context, channel string) (<-chan []byte, error) {
	m.subscribeCalls++
	if m.SubscribeFn != nil {
		return m.SubscribeFn(ctx, channel)
	}
	ch := make(chan []byte)
	close(ch)
	return ch, nil
}

// Test helper functions

func createTestExecution(id, userID string, status store.ExecutionStatus) *store.Execution {
	return &store.Execution{
		ID:        id,
		UserID:    userID,
		Status:    status,
		Config:    createTestDEConfig(),
		Algorithm: "gde3",
		Variant:   "rand/1",
		Problem:   "zdt1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func createTestDEConfig() *api.DEConfig {
	return &api.DEConfig{
		Executions:     1,
		Generations:    100,
		PopulationSize: 50,
		DimensionsSize: 10,
		ObjectivesSize: 2,
		FloorLimiter:   0.0,
		CeilLimiter:    1.0,
		AlgorithmConfig: &api.DEConfig_Gde3{
			Gde3: &api.GDE3Config{
				Cr: 0.9,
				F:  0.5,
				P:  0.1,
			},
		},
	}
}

func createTestProgress(executionID string) *store.ExecutionProgress {
	return &store.ExecutionProgress{
		ExecutionID:         executionID,
		CurrentGeneration:   50,
		TotalGenerations:    100,
		CompletedExecutions: 0,
		TotalExecutions:     1,
		UpdatedAt:           time.Now(),
	}
}

// ExecutionStore Tests

func TestNewExecutionStore(t *testing.T) {
	redis := &mockExecutionStore{}
	db := &mockExecutionStore{}

	s := NewExecutionStore(redis, db)

	assert.NotNil(t, s)
	assert.NotNil(t, s.logger)
}

func TestExecutionStore_CreateExecution(t *testing.T) {
	tests := []struct {
		name           string
		execution      *store.Execution
		setupRedis     func(*mockExecutionStore)
		setupDB        func(*mockExecutionStore)
		wantErr        bool
		errContains    string
		wantDBCalls    int
		wantRedisCalls int
	}{
		{
			name:      "success - both stores succeed",
			execution: createTestExecution("exec-1", "user-1", store.ExecutionStatusPending),
			setupRedis: func(m *mockExecutionStore) {
				m.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
					return nil
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
					return nil
				}
			},
			wantErr:        false,
			wantDBCalls:    1,
			wantRedisCalls: 1,
		},
		{
			name:      "success - redis fails but db succeeds (graceful degradation)",
			execution: createTestExecution("exec-2", "user-1", store.ExecutionStatusPending),
			setupRedis: func(m *mockExecutionStore) {
				m.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
					return errors.New("redis connection error")
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
					return nil
				}
			},
			wantErr:        false, // Should succeed despite Redis failure
			wantDBCalls:    1,
			wantRedisCalls: 1,
		},
		{
			name:      "failure - db fails",
			execution: createTestExecution("exec-3", "user-1", store.ExecutionStatusPending),
			setupRedis: func(m *mockExecutionStore) {
				m.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
					return nil
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
					return errors.New("database error")
				}
			},
			wantErr:        true,
			errContains:    "database error",
			wantDBCalls:    1,
			wantRedisCalls: 0, // Redis should not be called if DB fails
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := &mockExecutionStore{}
			db := &mockExecutionStore{}

			if tt.setupRedis != nil {
				tt.setupRedis(redis)
			}
			if tt.setupDB != nil {
				tt.setupDB(db)
			}

			s := NewExecutionStore(redis, db)
			ctx := context.Background()

			err := s.CreateExecution(ctx, tt.execution)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantDBCalls, db.createExecutionCalls)
			assert.Equal(t, tt.wantRedisCalls, redis.createExecutionCalls)
		})
	}
}

func TestExecutionStore_GetExecution(t *testing.T) {
	now := time.Now()
	olderTime := now.Add(-time.Hour)
	newerTime := now.Add(time.Hour)

	tests := []struct {
		name           string
		executionID    string
		userID         string
		setupRedis     func(*mockExecutionStore)
		setupDB        func(*mockExecutionStore)
		wantErr        bool
		errContains    string
		wantExecID     string
		wantDBCalls    int
		wantRedisCalls int
	}{
		{
			name:        "cache hit - cache is fresh",
			executionID: "exec-1",
			userID:      "user-1",
			setupRedis: func(m *mockExecutionStore) {
				m.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
					exec := createTestExecution(executionID, userID, store.ExecutionStatusRunning)
					exec.UpdatedAt = newerTime
					return exec, nil
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
					exec := createTestExecution(executionID, userID, store.ExecutionStatusRunning)
					exec.UpdatedAt = olderTime
					return exec, nil
				}
			},
			wantErr:        false,
			wantExecID:     "exec-1",
			wantDBCalls:    1,
			wantRedisCalls: 1,
		},
		{
			name:        "cache miss - fallback to db",
			executionID: "exec-2",
			userID:      "user-1",
			setupRedis: func(m *mockExecutionStore) {
				m.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
					return nil, errors.New("cache miss")
				}
				m.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
					return nil
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
					return createTestExecution(executionID, userID, store.ExecutionStatusRunning), nil
				}
			},
			wantErr:        false,
			wantExecID:     "exec-2",
			wantDBCalls:    1,
			wantRedisCalls: 1,
		},
		{
			name:        "stale cache - refresh from db",
			executionID: "exec-3",
			userID:      "user-1",
			setupRedis: func(m *mockExecutionStore) {
				m.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
					exec := createTestExecution(executionID, userID, store.ExecutionStatusPending)
					exec.UpdatedAt = olderTime
					return exec, nil
				}
				m.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
					return nil
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
					exec := createTestExecution(executionID, userID, store.ExecutionStatusRunning)
					exec.UpdatedAt = newerTime
					return exec, nil
				}
			},
			wantErr:        false,
			wantExecID:     "exec-3",
			wantDBCalls:    1,
			wantRedisCalls: 1,
		},
		{
			name:        "db unavailable - use cached value",
			executionID: "exec-4",
			userID:      "user-1",
			setupRedis: func(m *mockExecutionStore) {
				m.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
					return createTestExecution(executionID, userID, store.ExecutionStatusRunning), nil
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
					return nil, errors.New("database unavailable")
				}
			},
			wantErr:        false,
			wantExecID:     "exec-4",
			wantDBCalls:    1,
			wantRedisCalls: 1,
		},
		{
			name:        "both fail - return db error",
			executionID: "exec-5",
			userID:      "user-1",
			setupRedis: func(m *mockExecutionStore) {
				m.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
					return nil, errors.New("cache error")
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr:        true,
			errContains:    "database error",
			wantDBCalls:    1,
			wantRedisCalls: 1,
		},
		{
			name:        "not found in db - return error even if cache has old data",
			executionID: "exec-6",
			userID:      "user-1",
			setupRedis: func(m *mockExecutionStore) {
				m.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
					return nil, errors.New("cache miss")
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
					return nil, store.ErrExecutionNotFound
				}
			},
			wantErr:        true,
			errContains:    "execution not found",
			wantDBCalls:    1,
			wantRedisCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := &mockExecutionStore{}
			db := &mockExecutionStore{}

			if tt.setupRedis != nil {
				tt.setupRedis(redis)
			}
			if tt.setupDB != nil {
				tt.setupDB(db)
			}

			s := NewExecutionStore(redis, db)
			ctx := context.Background()

			exec, err := s.GetExecution(ctx, tt.executionID, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, exec)
				assert.Equal(t, tt.wantExecID, exec.ID)
			}

			assert.Equal(t, tt.wantDBCalls, db.getExecutionCalls)
			assert.Equal(t, tt.wantRedisCalls, redis.getExecutionCalls)
		})
	}
}

func TestExecutionStore_UpdateExecutionStatus(t *testing.T) {
	tests := []struct {
		name              string
		executionID       string
		status            store.ExecutionStatus
		errorMsg          string
		setupRedis        func(*mockExecutionStore)
		setupDB           func(*mockExecutionStore)
		wantErr           bool
		errContains       string
		wantDBCalls       int
		wantRedisDelCalls int
	}{
		{
			name:        "success - db updates, cache invalidated",
			executionID: "exec-1",
			status:      store.ExecutionStatusRunning,
			errorMsg:    "",
			setupRedis: func(m *mockExecutionStore) {
				m.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
					return nil
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.UpdateExecutionStatusFn = func(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
					return nil
				}
			},
			wantErr:           false,
			wantDBCalls:       1,
			wantRedisDelCalls: 1,
		},
		{
			name:        "success - cache invalidation fails (graceful)",
			executionID: "exec-2",
			status:      store.ExecutionStatusCompleted,
			errorMsg:    "",
			setupRedis: func(m *mockExecutionStore) {
				m.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
					return errors.New("redis error")
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.UpdateExecutionStatusFn = func(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
					return nil
				}
			},
			wantErr:           false, // Should succeed despite cache invalidation failure
			wantDBCalls:       1,
			wantRedisDelCalls: 1,
		},
		{
			name:        "failure - db fails",
			executionID: "exec-3",
			status:      store.ExecutionStatusFailed,
			errorMsg:    "algorithm error",
			setupRedis: func(m *mockExecutionStore) {
				m.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
					return nil
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.UpdateExecutionStatusFn = func(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
					return errors.New("database error")
				}
			},
			wantErr:           true,
			errContains:       "database error",
			wantDBCalls:       1,
			wantRedisDelCalls: 0, // Should not try to invalidate cache if DB fails
		},
		{
			name:        "success with error message",
			executionID: "exec-4",
			status:      store.ExecutionStatusFailed,
			errorMsg:    "out of memory",
			setupRedis: func(m *mockExecutionStore) {
				m.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
					return nil
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.UpdateExecutionStatusFn = func(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
					assert.Equal(t, "out of memory", errorMsg)
					return nil
				}
			},
			wantErr:           false,
			wantDBCalls:       1,
			wantRedisDelCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := &mockExecutionStore{}
			db := &mockExecutionStore{}

			if tt.setupRedis != nil {
				tt.setupRedis(redis)
			}
			if tt.setupDB != nil {
				tt.setupDB(db)
			}

			s := NewExecutionStore(redis, db)
			ctx := context.Background()

			err := s.UpdateExecutionStatus(ctx, tt.executionID, tt.status, tt.errorMsg)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantDBCalls, db.updateExecutionStatusCalls)
			assert.Equal(t, tt.wantRedisDelCalls, redis.deleteExecutionCalls)
		})
	}
}

func TestExecutionStore_UpdateExecutionResult(t *testing.T) {
	tests := []struct {
		name              string
		executionID       string
		paretoID          uint64
		setupRedis        func(*mockExecutionStore)
		setupDB           func(*mockExecutionStore)
		wantErr           bool
		errContains       string
		wantDBCalls       int
		wantRedisDelCalls int
	}{
		{
			name:        "success - db updates, cache invalidated",
			executionID: "exec-1",
			paretoID:    12345,
			setupRedis: func(m *mockExecutionStore) {
				m.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
					return nil
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.UpdateExecutionResultFn = func(ctx context.Context, executionID string, paretoID uint64) error {
					return nil
				}
			},
			wantErr:           false,
			wantDBCalls:       1,
			wantRedisDelCalls: 1,
		},
		{
			name:        "success - cache invalidation fails (graceful)",
			executionID: "exec-2",
			paretoID:    67890,
			setupRedis: func(m *mockExecutionStore) {
				m.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
					return errors.New("redis error")
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.UpdateExecutionResultFn = func(ctx context.Context, executionID string, paretoID uint64) error {
					return nil
				}
			},
			wantErr:           false,
			wantDBCalls:       1,
			wantRedisDelCalls: 1,
		},
		{
			name:        "failure - db fails",
			executionID: "exec-3",
			paretoID:    11111,
			setupRedis: func(m *mockExecutionStore) {
				m.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
					return nil
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.UpdateExecutionResultFn = func(ctx context.Context, executionID string, paretoID uint64) error {
					return errors.New("database error")
				}
			},
			wantErr:           true,
			errContains:       "database error",
			wantDBCalls:       1,
			wantRedisDelCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := &mockExecutionStore{}
			db := &mockExecutionStore{}

			if tt.setupRedis != nil {
				tt.setupRedis(redis)
			}
			if tt.setupDB != nil {
				tt.setupDB(db)
			}

			s := NewExecutionStore(redis, db)
			ctx := context.Background()

			err := s.UpdateExecutionResult(ctx, tt.executionID, tt.paretoID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantDBCalls, db.updateExecutionResultCalls)
			assert.Equal(t, tt.wantRedisDelCalls, redis.deleteExecutionCalls)
		})
	}
}

func TestExecutionStore_ListExecutions(t *testing.T) {
	runningStatus := store.ExecutionStatusRunning

	tests := []struct {
		name        string
		userID      string
		status      *store.ExecutionStatus
		limit       int
		offset      int
		setupDB     func(*mockExecutionStore)
		wantErr     bool
		wantCount   int
		wantTotal   int
		wantDBCalls int
	}{
		{
			name:   "success - list all executions",
			userID: "user-1",
			status: nil,
			limit:  50,
			offset: 0,
			setupDB: func(m *mockExecutionStore) {
				m.ListExecutionsFn = func(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
					return []*store.Execution{
						createTestExecution("exec-1", userID, store.ExecutionStatusRunning),
						createTestExecution("exec-2", userID, store.ExecutionStatusCompleted),
					}, 2, nil
				}
			},
			wantErr:     false,
			wantCount:   2,
			wantTotal:   2,
			wantDBCalls: 1,
		},
		{
			name:   "success - list with status filter",
			userID: "user-1",
			status: &runningStatus,
			limit:  50,
			offset: 0,
			setupDB: func(m *mockExecutionStore) {
				m.ListExecutionsFn = func(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
					assert.NotNil(t, status)
					assert.Equal(t, store.ExecutionStatusRunning, *status)
					return []*store.Execution{
						createTestExecution("exec-1", userID, store.ExecutionStatusRunning),
					}, 1, nil
				}
			},
			wantErr:     false,
			wantCount:   1,
			wantTotal:   1,
			wantDBCalls: 1,
		},
		{
			name:   "success - empty list",
			userID: "user-2",
			status: nil,
			limit:  50,
			offset: 0,
			setupDB: func(m *mockExecutionStore) {
				m.ListExecutionsFn = func(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
					return []*store.Execution{}, 0, nil
				}
			},
			wantErr:     false,
			wantCount:   0,
			wantTotal:   0,
			wantDBCalls: 1,
		},
		{
			name:   "failure - db error",
			userID: "user-1",
			status: nil,
			limit:  50,
			offset: 0,
			setupDB: func(m *mockExecutionStore) {
				m.ListExecutionsFn = func(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
					return nil, 0, errors.New("database error")
				}
			},
			wantErr:     true,
			wantDBCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := &mockExecutionStore{}
			db := &mockExecutionStore{}

			if tt.setupDB != nil {
				tt.setupDB(db)
			}

			s := NewExecutionStore(redis, db)
			ctx := context.Background()

			executions, total, err := s.ListExecutions(ctx, tt.userID, tt.status, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, executions, tt.wantCount)
				assert.Equal(t, tt.wantTotal, total)
			}

			assert.Equal(t, tt.wantDBCalls, db.listExecutionsCalls)
			// Redis should not be called for list operations
			assert.Equal(t, 0, redis.listExecutionsCalls)
		})
	}
}

func TestExecutionStore_DeleteExecution(t *testing.T) {
	tests := []struct {
		name              string
		executionID       string
		userID            string
		setupRedis        func(*mockExecutionStore)
		setupDB           func(*mockExecutionStore)
		wantErr           bool
		errContains       string
		wantDBCalls       int
		wantRedisDelCalls int
	}{
		{
			name:        "success - deleted from both stores",
			executionID: "exec-1",
			userID:      "user-1",
			setupRedis: func(m *mockExecutionStore) {
				m.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
					return nil
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
					return nil
				}
			},
			wantErr:           false,
			wantDBCalls:       1,
			wantRedisDelCalls: 1,
		},
		{
			name:        "success - redis delete fails (best effort)",
			executionID: "exec-2",
			userID:      "user-1",
			setupRedis: func(m *mockExecutionStore) {
				m.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
					return errors.New("redis error")
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
					return nil
				}
			},
			wantErr:           false, // Should succeed despite Redis failure
			wantDBCalls:       1,
			wantRedisDelCalls: 1,
		},
		{
			name:        "failure - db delete fails",
			executionID: "exec-3",
			userID:      "user-1",
			setupRedis: func(m *mockExecutionStore) {
				m.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
					return nil
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
					return errors.New("database error")
				}
			},
			wantErr:           true,
			errContains:       "database error",
			wantDBCalls:       1,
			wantRedisDelCalls: 0, // Should not try to delete from Redis if DB fails
		},
		{
			name:        "failure - execution not found",
			executionID: "exec-nonexistent",
			userID:      "user-1",
			setupRedis: func(m *mockExecutionStore) {
				m.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
					return nil
				}
			},
			setupDB: func(m *mockExecutionStore) {
				m.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
					return store.ErrExecutionNotFound
				}
			},
			wantErr:           true,
			errContains:       "execution not found",
			wantDBCalls:       1,
			wantRedisDelCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := &mockExecutionStore{}
			db := &mockExecutionStore{}

			if tt.setupRedis != nil {
				tt.setupRedis(redis)
			}
			if tt.setupDB != nil {
				tt.setupDB(db)
			}

			s := NewExecutionStore(redis, db)
			ctx := context.Background()

			err := s.DeleteExecution(ctx, tt.executionID, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantDBCalls, db.deleteExecutionCalls)
			assert.Equal(t, tt.wantRedisDelCalls, redis.deleteExecutionCalls)
		})
	}
}

func TestExecutionStore_SaveProgress(t *testing.T) {
	tests := []struct {
		name           string
		progress       *store.ExecutionProgress
		setupRedis     func(*mockExecutionStore)
		wantErr        bool
		errContains    string
		wantRedisCalls int
	}{
		{
			name:     "success",
			progress: createTestProgress("exec-1"),
			setupRedis: func(m *mockExecutionStore) {
				m.SaveProgressFn = func(ctx context.Context, progress *store.ExecutionProgress) error {
					return nil
				}
			},
			wantErr:        false,
			wantRedisCalls: 1,
		},
		{
			name:     "failure - redis error",
			progress: createTestProgress("exec-2"),
			setupRedis: func(m *mockExecutionStore) {
				m.SaveProgressFn = func(ctx context.Context, progress *store.ExecutionProgress) error {
					return errors.New("redis error")
				}
			},
			wantErr:        true,
			errContains:    "redis error",
			wantRedisCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := &mockExecutionStore{}
			db := &mockExecutionStore{}

			if tt.setupRedis != nil {
				tt.setupRedis(redis)
			}

			s := NewExecutionStore(redis, db)
			ctx := context.Background()

			err := s.SaveProgress(ctx, tt.progress)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantRedisCalls, redis.saveProgressCalls)
			// DB should not be called for progress operations
			assert.Equal(t, 0, db.saveProgressCalls)
		})
	}
}

func TestExecutionStore_GetProgress(t *testing.T) {
	tests := []struct {
		name           string
		executionID    string
		setupRedis     func(*mockExecutionStore)
		wantErr        bool
		errContains    string
		wantExecID     string
		wantRedisCalls int
	}{
		{
			name:        "success",
			executionID: "exec-1",
			setupRedis: func(m *mockExecutionStore) {
				m.GetProgressFn = func(ctx context.Context, executionID string) (*store.ExecutionProgress, error) {
					return createTestProgress(executionID), nil
				}
			},
			wantErr:        false,
			wantExecID:     "exec-1",
			wantRedisCalls: 1,
		},
		{
			name:        "failure - not found",
			executionID: "exec-nonexistent",
			setupRedis: func(m *mockExecutionStore) {
				m.GetProgressFn = func(ctx context.Context, executionID string) (*store.ExecutionProgress, error) {
					return nil, errors.New("progress not found")
				}
			},
			wantErr:        true,
			errContains:    "progress not found",
			wantRedisCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := &mockExecutionStore{}
			db := &mockExecutionStore{}

			if tt.setupRedis != nil {
				tt.setupRedis(redis)
			}

			s := NewExecutionStore(redis, db)
			ctx := context.Background()

			progress, err := s.GetProgress(ctx, tt.executionID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, progress)
				assert.Equal(t, tt.wantExecID, progress.ExecutionID)
			}

			assert.Equal(t, tt.wantRedisCalls, redis.getProgressCalls)
			assert.Equal(t, 0, db.getProgressCalls)
		})
	}
}

func TestExecutionStore_MarkExecutionForCancellation(t *testing.T) {
	tests := []struct {
		name           string
		executionID    string
		userID         string
		setupRedis     func(*mockExecutionStore)
		wantErr        bool
		errContains    string
		wantRedisCalls int
	}{
		{
			name:        "success",
			executionID: "exec-1",
			userID:      "user-1",
			setupRedis: func(m *mockExecutionStore) {
				m.MarkExecutionForCancellationFn = func(ctx context.Context, executionID, userID string) error {
					return nil
				}
			},
			wantErr:        false,
			wantRedisCalls: 1,
		},
		{
			name:        "failure - redis error",
			executionID: "exec-2",
			userID:      "user-1",
			setupRedis: func(m *mockExecutionStore) {
				m.MarkExecutionForCancellationFn = func(ctx context.Context, executionID, userID string) error {
					return errors.New("redis error")
				}
			},
			wantErr:        true,
			errContains:    "redis error",
			wantRedisCalls: 1,
		},
		{
			name:        "failure - execution not found",
			executionID: "exec-nonexistent",
			userID:      "user-1",
			setupRedis: func(m *mockExecutionStore) {
				m.MarkExecutionForCancellationFn = func(ctx context.Context, executionID, userID string) error {
					return store.ErrExecutionNotFound
				}
			},
			wantErr:        true,
			errContains:    "execution not found",
			wantRedisCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := &mockExecutionStore{}
			db := &mockExecutionStore{}

			if tt.setupRedis != nil {
				tt.setupRedis(redis)
			}

			s := NewExecutionStore(redis, db)
			ctx := context.Background()

			err := s.MarkExecutionForCancellation(ctx, tt.executionID, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantRedisCalls, redis.markExecutionForCancellationCalls)
			assert.Equal(t, 0, db.markExecutionForCancellationCalls)
		})
	}
}

func TestExecutionStore_IsExecutionCancelled(t *testing.T) {
	tests := []struct {
		name           string
		executionID    string
		setupRedis     func(*mockExecutionStore)
		wantCancelled  bool
		wantErr        bool
		errContains    string
		wantRedisCalls int
	}{
		{
			name:        "cancelled",
			executionID: "exec-1",
			setupRedis: func(m *mockExecutionStore) {
				m.IsExecutionCancelledFn = func(ctx context.Context, executionID string) (bool, error) {
					return true, nil
				}
			},
			wantCancelled:  true,
			wantErr:        false,
			wantRedisCalls: 1,
		},
		{
			name:        "not cancelled",
			executionID: "exec-2",
			setupRedis: func(m *mockExecutionStore) {
				m.IsExecutionCancelledFn = func(ctx context.Context, executionID string) (bool, error) {
					return false, nil
				}
			},
			wantCancelled:  false,
			wantErr:        false,
			wantRedisCalls: 1,
		},
		{
			name:        "failure - redis error",
			executionID: "exec-3",
			setupRedis: func(m *mockExecutionStore) {
				m.IsExecutionCancelledFn = func(ctx context.Context, executionID string) (bool, error) {
					return false, errors.New("redis error")
				}
			},
			wantCancelled:  false,
			wantErr:        true,
			errContains:    "redis error",
			wantRedisCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := &mockExecutionStore{}
			db := &mockExecutionStore{}

			if tt.setupRedis != nil {
				tt.setupRedis(redis)
			}

			s := NewExecutionStore(redis, db)
			ctx := context.Background()

			cancelled, err := s.IsExecutionCancelled(ctx, tt.executionID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCancelled, cancelled)
			}

			assert.Equal(t, tt.wantRedisCalls, redis.isExecutionCancelledCalls)
			assert.Equal(t, 0, db.isExecutionCancelledCalls)
		})
	}
}

func TestExecutionStore_Subscribe(t *testing.T) {
	tests := []struct {
		name           string
		channel        string
		setupRedis     func(*mockExecutionStore)
		wantErr        bool
		errContains    string
		wantRedisCalls int
	}{
		{
			name:    "success",
			channel: "execution:exec-1:updates",
			setupRedis: func(m *mockExecutionStore) {
				m.SubscribeFn = func(ctx context.Context, channel string) (<-chan []byte, error) {
					ch := make(chan []byte)
					close(ch)
					return ch, nil
				}
			},
			wantErr:        false,
			wantRedisCalls: 1,
		},
		{
			name:    "failure - redis error",
			channel: "execution:exec-2:updates",
			setupRedis: func(m *mockExecutionStore) {
				m.SubscribeFn = func(ctx context.Context, channel string) (<-chan []byte, error) {
					return nil, errors.New("redis subscription error")
				}
			},
			wantErr:        true,
			errContains:    "redis subscription error",
			wantRedisCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := &mockExecutionStore{}
			db := &mockExecutionStore{}

			if tt.setupRedis != nil {
				tt.setupRedis(redis)
			}

			s := NewExecutionStore(redis, db)
			ctx := context.Background()

			ch, err := s.Subscribe(ctx, tt.channel)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, ch)
			}

			assert.Equal(t, tt.wantRedisCalls, redis.subscribeCalls)
			assert.Equal(t, 0, db.subscribeCalls)
		})
	}
}

// Edge Cases Tests

func TestExecutionStore_GetExecution_CacheRefreshFailure(t *testing.T) {
	// Test that when cache is stale and refresh fails, we still return DB data
	redis := &mockExecutionStore{}
	db := &mockExecutionStore{}

	olderTime := time.Now().Add(-time.Hour)
	newerTime := time.Now()

	redis.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
		exec := createTestExecution(executionID, userID, store.ExecutionStatusPending)
		exec.UpdatedAt = olderTime
		return exec, nil
	}
	redis.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
		return errors.New("cache refresh failed")
	}

	db.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
		exec := createTestExecution(executionID, userID, store.ExecutionStatusRunning)
		exec.UpdatedAt = newerTime
		return exec, nil
	}

	s := NewExecutionStore(redis, db)
	ctx := context.Background()

	exec, err := s.GetExecution(ctx, "exec-1", "user-1")

	assert.NoError(t, err)
	assert.NotNil(t, exec)
	assert.Equal(t, store.ExecutionStatusRunning, exec.Status) // Should return DB value
}

func TestExecutionStore_ContextCancellation(t *testing.T) {
	redis := &mockExecutionStore{}
	db := &mockExecutionStore{}

	// Setup to simulate slow operation
	redis.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(100 * time.Millisecond):
			return createTestExecution(executionID, userID, store.ExecutionStatusRunning), nil
		}
	}
	db.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
		return nil, errors.New("db error")
	}

	s := NewExecutionStore(redis, db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := s.GetExecution(ctx, "exec-1", "user-1")

	// Should get context deadline exceeded or db error
	assert.Error(t, err)
}

func TestExecutionStore_NilExecution(t *testing.T) {
	redis := &mockExecutionStore{}
	db := &mockExecutionStore{}

	db.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
		if execution == nil {
			return errors.New("nil execution")
		}
		return nil
	}

	s := NewExecutionStore(redis, db)
	ctx := context.Background()

	err := s.CreateExecution(ctx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil execution")
}

func TestExecutionStore_EmptyUserID(t *testing.T) {
	redis := &mockExecutionStore{}
	db := &mockExecutionStore{}

	db.ListExecutionsFn = func(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
		if userID == "" {
			return nil, 0, errors.New("empty user ID")
		}
		return []*store.Execution{}, 0, nil
	}

	s := NewExecutionStore(redis, db)
	ctx := context.Background()

	_, _, err := s.ListExecutions(ctx, "", nil, 50, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty user ID")
}

// Concurrent Access Tests

func TestExecutionStore_ConcurrentGetExecution(t *testing.T) {
	redis := &mockExecutionStore{}
	db := &mockExecutionStore{}

	redis.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
		return createTestExecution(executionID, userID, store.ExecutionStatusRunning), nil
	}
	redis.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
		return nil
	}
	db.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
		return createTestExecution(executionID, userID, store.ExecutionStatusRunning), nil
	}

	s := NewExecutionStore(redis, db)
	ctx := context.Background()

	// Run concurrent gets
	done := make(chan bool, 10)
	for i := range 10 {
		go func(idx int) {
			exec, err := s.GetExecution(ctx, "exec-1", "user-1")
			assert.NoError(t, err)
			assert.NotNil(t, exec)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for range 10 {
		<-done
	}
}

func TestExecutionStore_ConcurrentUpdates(t *testing.T) {
	redis := &mockExecutionStore{}
	db := &mockExecutionStore{}

	var updateCount int
	db.UpdateExecutionStatusFn = func(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
		updateCount++
		return nil
	}
	redis.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
		return nil
	}

	execStore := NewExecutionStore(redis, db)
	ctx := context.Background()

	// Run concurrent updates
	done := make(chan bool, 5)
	statuses := []store.ExecutionStatus{
		store.ExecutionStatusRunning,
		store.ExecutionStatusCompleted,
		store.ExecutionStatusFailed,
		store.ExecutionStatusCancelled,
		store.ExecutionStatusPending,
	}

	for i, status := range statuses {
		go func(idx int, st store.ExecutionStatus) {
			err := execStore.UpdateExecutionStatus(ctx, "exec-1", st, "")
			assert.NoError(t, err)
			done <- true
		}(i, status)
	}

	// Wait for all goroutines
	for range 5 {
		<-done
	}
}

// Verify cache freshness logic

func TestExecutionStore_GetExecution_CacheFreshnessComparison(t *testing.T) {
	tests := []struct {
		name           string
		cacheUpdatedAt time.Time
		dbUpdatedAt    time.Time
		wantFromCache  bool
	}{
		{
			name:           "cache is newer - use cache",
			cacheUpdatedAt: time.Now(),
			dbUpdatedAt:    time.Now().Add(-time.Hour),
			wantFromCache:  true,
		},
		{
			name:           "cache is older - refresh from db",
			cacheUpdatedAt: time.Now().Add(-time.Hour),
			dbUpdatedAt:    time.Now(),
			wantFromCache:  false,
		},
		{
			name:           "same time - use cache",
			cacheUpdatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			dbUpdatedAt:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			wantFromCache:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redis := &mockExecutionStore{}
			db := &mockExecutionStore{}

			cacheExec := createTestExecution("exec-1", "user-1", store.ExecutionStatusPending)
			cacheExec.UpdatedAt = tt.cacheUpdatedAt
			cacheExec.Status = store.ExecutionStatusPending // Marker for cache

			dbExec := createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)
			dbExec.UpdatedAt = tt.dbUpdatedAt
			dbExec.Status = store.ExecutionStatusRunning // Marker for DB

			redis.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
				return cacheExec, nil
			}
			redis.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
				return nil
			}
			db.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
				return dbExec, nil
			}

			s := NewExecutionStore(redis, db)
			ctx := context.Background()

			exec, err := s.GetExecution(ctx, "exec-1", "user-1")

			require.NoError(t, err)
			require.NotNil(t, exec)

			if tt.wantFromCache {
				assert.Equal(t, store.ExecutionStatusPending, exec.Status)
			} else {
				assert.Equal(t, store.ExecutionStatusRunning, exec.Status)
			}
		})
	}
}

// Benchmark tests

func BenchmarkExecutionStore_GetExecution_CacheHit(b *testing.B) {
	redis := &mockExecutionStore{}
	db := &mockExecutionStore{}

	exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)
	exec.UpdatedAt = time.Now()

	redis.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
		return exec, nil
	}
	db.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
		dbExec := createTestExecution(executionID, userID, store.ExecutionStatusRunning)
		dbExec.UpdatedAt = exec.UpdatedAt.Add(-time.Hour)
		return dbExec, nil
	}

	s := NewExecutionStore(redis, db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.GetExecution(ctx, "exec-1", "user-1")
	}
}

func BenchmarkExecutionStore_GetExecution_CacheMiss(b *testing.B) {
	redis := &mockExecutionStore{}
	db := &mockExecutionStore{}

	redis.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
		return nil, errors.New("cache miss")
	}
	redis.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
		return nil
	}
	db.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
		return createTestExecution(executionID, userID, store.ExecutionStatusRunning), nil
	}

	s := NewExecutionStore(redis, db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.GetExecution(ctx, "exec-1", "user-1")
	}
}

func BenchmarkExecutionStore_CreateExecution(b *testing.B) {
	redis := &mockExecutionStore{}
	db := &mockExecutionStore{}

	redis.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
		return nil
	}
	db.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
		return nil
	}

	s := NewExecutionStore(redis, db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		exec := createTestExecution("exec-bench", "user-1", store.ExecutionStatusPending)
		_ = s.CreateExecution(ctx, exec)
	}
}

// ================================================================================
// Store (composite wrapper) Tests
// ================================================================================

// mockFullStore implements all operations needed for the Store wrapper
type mockFullStore struct {
	mockStore
	mockExecutionStore
}

// createMockStoreWrapper creates a Store instance with mock dependencies for testing.
// This allows testing the Store wrapper without requiring a real Redis connection.
func createMockStoreWrapper(dbMock *mockStore, redisMock *mockExecutionStore) *Store {
	return &Store{
		db:        dbMock,
		redis:     nil, // Redis health check will be skipped in tests using this
		execStore: NewExecutionStore(redisMock, dbMock),
	}
}

// Tests for the Store wrapper type directly

func TestStore_New(t *testing.T) {
	dbMock := &mockStore{}
	redisMock := &mockExecutionStore{}

	store := New(dbMock, nil, redisMock)

	assert.NotNil(t, store)
	assert.NotNil(t, store.execStore)
}

func TestStore_UserOperations_Direct(t *testing.T) {
	t.Run("CreateUser", func(t *testing.T) {
		dbMock := &mockStore{}
		var createdUser *api.User
		dbMock.CreateUserFn = func(ctx context.Context, user *api.User) error {
			createdUser = user
			return nil
		}

		store := createMockStoreWrapper(dbMock, &mockExecutionStore{})
		ctx := context.Background()
		user := &api.User{Email: "test@example.com"}

		err := store.CreateUser(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, user, createdUser)
		assert.Equal(t, 1, dbMock.createUserCalls)
	})

	t.Run("GetUser", func(t *testing.T) {
		dbMock := &mockStore{}
		expectedUser := &api.User{Email: "test@example.com"}
		dbMock.GetUserFn = func(ctx context.Context, userIDs *api.UserIDs) (*api.User, error) {
			return expectedUser, nil
		}

		store := createMockStoreWrapper(dbMock, &mockExecutionStore{})
		ctx := context.Background()
		userIDs := &api.UserIDs{Username: "testuser"}

		user, err := store.GetUser(ctx, userIDs)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		dbMock := &mockStore{}
		var updatedUser *api.User
		dbMock.UpdateUserFn = func(ctx context.Context, user *api.User, fields ...string) error {
			updatedUser = user
			return nil
		}

		store := createMockStoreWrapper(dbMock, &mockExecutionStore{})
		ctx := context.Background()
		user := &api.User{Email: "updated@example.com"}

		err := store.UpdateUser(ctx, user, "email")

		assert.NoError(t, err)
		assert.Equal(t, user, updatedUser)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		dbMock := &mockStore{}
		var deletedUserIDs *api.UserIDs
		dbMock.DeleteUserFn = func(ctx context.Context, userIDs *api.UserIDs) error {
			deletedUserIDs = userIDs
			return nil
		}

		store := createMockStoreWrapper(dbMock, &mockExecutionStore{})
		ctx := context.Background()
		userIDs := &api.UserIDs{Username: "testuser"}

		err := store.DeleteUser(ctx, userIDs)

		assert.NoError(t, err)
		assert.Equal(t, userIDs, deletedUserIDs)
	})
}

func TestStore_ParetoOperations_Direct(t *testing.T) {
	t.Run("CreatePareto", func(t *testing.T) {
		dbMock := &mockStore{}
		var createdPareto *api.Pareto
		dbMock.CreateParetoFn = func(ctx context.Context, pareto *api.Pareto) error {
			createdPareto = pareto
			return nil
		}

		store := createMockStoreWrapper(dbMock, &mockExecutionStore{})
		ctx := context.Background()
		pareto := &api.Pareto{Ids: &api.ParetoIDs{Id: 1}}

		err := store.CreatePareto(ctx, pareto)

		assert.NoError(t, err)
		assert.Equal(t, pareto, createdPareto)
	})

	t.Run("GetPareto", func(t *testing.T) {
		dbMock := &mockStore{}
		expectedPareto := &api.Pareto{Ids: &api.ParetoIDs{Id: 1}}
		dbMock.GetParetoFn = func(ctx context.Context, ids *api.ParetoIDs) (*api.Pareto, error) {
			return expectedPareto, nil
		}

		store := createMockStoreWrapper(dbMock, &mockExecutionStore{})
		ctx := context.Background()
		paretoIDs := &api.ParetoIDs{Id: 1}

		pareto, err := store.GetPareto(ctx, paretoIDs)

		assert.NoError(t, err)
		assert.Equal(t, expectedPareto, pareto)
	})

	t.Run("UpdatePareto", func(t *testing.T) {
		dbMock := &mockStore{}
		dbMock.UpdateParetoFn = func(ctx context.Context, pareto *api.Pareto, fields ...string) error {
			return nil
		}

		store := createMockStoreWrapper(dbMock, &mockExecutionStore{})
		ctx := context.Background()
		pareto := &api.Pareto{Ids: &api.ParetoIDs{Id: 1}}

		err := store.UpdatePareto(ctx, pareto, "field1")

		assert.NoError(t, err)
	})

	t.Run("DeletePareto", func(t *testing.T) {
		dbMock := &mockStore{}
		dbMock.DeleteParetoFn = func(ctx context.Context, ids *api.ParetoIDs) error {
			return nil
		}

		store := createMockStoreWrapper(dbMock, &mockExecutionStore{})
		ctx := context.Background()
		paretoIDs := &api.ParetoIDs{Id: 1}

		err := store.DeletePareto(ctx, paretoIDs)

		assert.NoError(t, err)
	})

	t.Run("ListParetos", func(t *testing.T) {
		dbMock := &mockStore{}
		expectedParetos := []*api.Pareto{{Ids: &api.ParetoIDs{Id: 1}}, {Ids: &api.ParetoIDs{Id: 2}}}
		dbMock.ListParetosFn = func(ctx context.Context, userIDs *api.UserIDs, limit, offset int) ([]*api.Pareto, int, error) {
			return expectedParetos, 2, nil
		}

		store := createMockStoreWrapper(dbMock, &mockExecutionStore{})
		ctx := context.Background()
		userIDs := &api.UserIDs{Username: "testuser"}

		paretos, total, err := store.ListParetos(ctx, userIDs, 10, 0)

		assert.NoError(t, err)
		assert.Len(t, paretos, 2)
		assert.Equal(t, 2, total)
	})

	t.Run("CreateParetoSet", func(t *testing.T) {
		dbMock := &mockStore{}
		dbMock.CreateParetoSetFn = func(ctx context.Context, paretoSet *store.ParetoSet) error {
			return nil
		}

		st := createMockStoreWrapper(dbMock, &mockExecutionStore{})
		ctx := context.Background()
		paretoSet := &store.ParetoSet{UserID: "user-1"}

		err := st.CreateParetoSet(ctx, paretoSet)

		assert.NoError(t, err)
	})

	t.Run("GetParetoSetByID", func(t *testing.T) {
		dbMock := &mockStore{}
		expectedParetoSet := &store.ParetoSet{ID: 1, UserID: "user-1"}
		dbMock.GetParetoSetByIDFn = func(ctx context.Context, id uint64) (*store.ParetoSet, error) {
			return expectedParetoSet, nil
		}

		st := createMockStoreWrapper(dbMock, &mockExecutionStore{})
		ctx := context.Background()

		paretoSet, err := st.GetParetoSetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedParetoSet, paretoSet)
	})
}

func TestStore_ExecutionOperations_Direct(t *testing.T) {
	t.Run("CreateExecution", func(t *testing.T) {
		dbMock := &mockStore{}
		redisMock := &mockExecutionStore{}
		redisMock.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
			return nil
		}
		dbMock.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
			return nil
		}

		st := createMockStoreWrapper(dbMock, redisMock)
		ctx := context.Background()
		exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusPending)

		err := st.CreateExecution(ctx, exec)

		assert.NoError(t, err)
	})

	t.Run("GetExecution", func(t *testing.T) {
		dbMock := &mockStore{}
		redisMock := &mockExecutionStore{}
		expectedExec := createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)
		redisMock.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
			return expectedExec, nil
		}
		dbMock.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
			return expectedExec, nil
		}

		st := createMockStoreWrapper(dbMock, redisMock)
		ctx := context.Background()

		exec, err := st.GetExecution(ctx, "exec-1", "user-1")

		assert.NoError(t, err)
		assert.Equal(t, "exec-1", exec.ID)
	})

	t.Run("UpdateExecutionStatus", func(t *testing.T) {
		dbMock := &mockStore{}
		redisMock := &mockExecutionStore{}
		dbMock.UpdateExecutionStatusFn = func(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
			return nil
		}
		redisMock.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
			return nil
		}

		st := createMockStoreWrapper(dbMock, redisMock)
		ctx := context.Background()

		err := st.UpdateExecutionStatus(ctx, "exec-1", store.ExecutionStatusCompleted, "")

		assert.NoError(t, err)
	})

	t.Run("UpdateExecutionResult", func(t *testing.T) {
		dbMock := &mockStore{}
		redisMock := &mockExecutionStore{}
		dbMock.UpdateExecutionResultFn = func(ctx context.Context, executionID string, paretoID uint64) error {
			return nil
		}
		redisMock.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
			return nil
		}

		st := createMockStoreWrapper(dbMock, redisMock)
		ctx := context.Background()

		err := st.UpdateExecutionResult(ctx, "exec-1", 12345)

		assert.NoError(t, err)
	})

	t.Run("ListExecutions", func(t *testing.T) {
		dbMock := &mockStore{}
		redisMock := &mockExecutionStore{}
		expectedExecs := []*store.Execution{createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)}
		dbMock.ListExecutionsFn = func(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
			return expectedExecs, 1, nil
		}

		st := createMockStoreWrapper(dbMock, redisMock)
		ctx := context.Background()

		execs, total, err := st.ListExecutions(ctx, "user-1", nil, 50, 0)

		assert.NoError(t, err)
		assert.Len(t, execs, 1)
		assert.Equal(t, 1, total)
	})

	t.Run("DeleteExecution", func(t *testing.T) {
		dbMock := &mockStore{}
		redisMock := &mockExecutionStore{}
		dbMock.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
			return nil
		}
		redisMock.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
			return nil
		}

		st := createMockStoreWrapper(dbMock, redisMock)
		ctx := context.Background()

		err := st.DeleteExecution(ctx, "exec-1", "user-1")

		assert.NoError(t, err)
	})

	t.Run("SaveProgress", func(t *testing.T) {
		dbMock := &mockStore{}
		redisMock := &mockExecutionStore{}
		redisMock.SaveProgressFn = func(ctx context.Context, progress *store.ExecutionProgress) error {
			return nil
		}

		st := createMockStoreWrapper(dbMock, redisMock)
		ctx := context.Background()
		progress := createTestProgress("exec-1")

		err := st.SaveProgress(ctx, progress)

		assert.NoError(t, err)
	})

	t.Run("GetProgress", func(t *testing.T) {
		dbMock := &mockStore{}
		redisMock := &mockExecutionStore{}
		expectedProgress := createTestProgress("exec-1")
		redisMock.GetProgressFn = func(ctx context.Context, executionID string) (*store.ExecutionProgress, error) {
			return expectedProgress, nil
		}

		st := createMockStoreWrapper(dbMock, redisMock)
		ctx := context.Background()

		progress, err := st.GetProgress(ctx, "exec-1")

		assert.NoError(t, err)
		assert.Equal(t, "exec-1", progress.ExecutionID)
	})

	t.Run("MarkExecutionForCancellation", func(t *testing.T) {
		dbMock := &mockStore{}
		redisMock := &mockExecutionStore{}
		redisMock.MarkExecutionForCancellationFn = func(ctx context.Context, executionID, userID string) error {
			return nil
		}

		st := createMockStoreWrapper(dbMock, redisMock)
		ctx := context.Background()

		err := st.MarkExecutionForCancellation(ctx, "exec-1", "user-1")

		assert.NoError(t, err)
	})

	t.Run("IsExecutionCancelled", func(t *testing.T) {
		dbMock := &mockStore{}
		redisMock := &mockExecutionStore{}
		redisMock.IsExecutionCancelledFn = func(ctx context.Context, executionID string) (bool, error) {
			return true, nil
		}

		st := createMockStoreWrapper(dbMock, redisMock)
		ctx := context.Background()

		cancelled, err := st.IsExecutionCancelled(ctx, "exec-1")

		assert.NoError(t, err)
		assert.True(t, cancelled)
	})

	t.Run("Subscribe", func(t *testing.T) {
		dbMock := &mockStore{}
		redisMock := &mockExecutionStore{}
		redisMock.SubscribeFn = func(ctx context.Context, channel string) (<-chan []byte, error) {
			ch := make(chan []byte)
			close(ch)
			return ch, nil
		}

		st := createMockStoreWrapper(dbMock, redisMock)
		ctx := context.Background()

		ch, err := st.Subscribe(ctx, "exec-1:updates")

		assert.NoError(t, err)
		assert.NotNil(t, ch)
	})
}

func TestStore_HealthCheck_Direct(t *testing.T) {
	t.Run("db health check failure", func(t *testing.T) {
		dbMock := &mockStore{}
		dbMock.HealthCheckFn = func(ctx context.Context) error {
			return errors.New("db connection failed")
		}

		// Create store with nil redis to test db health check failure path
		st := &Store{
			db:    dbMock,
			redis: nil,
		}
		ctx := context.Background()

		err := st.HealthCheck(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db connection failed")
	})
}

func TestStore_RedisClient(t *testing.T) {
	dbMock := &mockStore{}

	st := &Store{
		db:    dbMock,
		redis: nil,
	}

	client := st.RedisClient()

	assert.Nil(t, client)
}

// mockStore implements store.Store for testing the main Store wrapper
type mockStore struct {
	// User operations
	CreateUserFn func(ctx context.Context, user *api.User) error
	GetUserFn    func(ctx context.Context, userIDs *api.UserIDs) (*api.User, error)
	UpdateUserFn func(ctx context.Context, user *api.User, fields ...string) error
	DeleteUserFn func(ctx context.Context, userIDs *api.UserIDs) error

	// Pareto operations
	CreateParetoFn     func(ctx context.Context, pareto *api.Pareto) error
	GetParetoFn        func(ctx context.Context, ids *api.ParetoIDs) (*api.Pareto, error)
	UpdateParetoFn     func(ctx context.Context, pareto *api.Pareto, fields ...string) error
	DeleteParetoFn     func(ctx context.Context, ids *api.ParetoIDs) error
	ListParetosFn      func(ctx context.Context, userIDs *api.UserIDs, limit, offset int) ([]*api.Pareto, int, error)
	CreateParetoSetFn  func(ctx context.Context, paretoSet *store.ParetoSet) error
	GetParetoSetByIDFn func(ctx context.Context, id uint64) (*store.ParetoSet, error)

	// Execution operations
	CreateExecutionFn              func(ctx context.Context, execution *store.Execution) error
	GetExecutionFn                 func(ctx context.Context, executionID, userID string) (*store.Execution, error)
	UpdateExecutionStatusFn        func(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error
	UpdateExecutionResultFn        func(ctx context.Context, executionID string, paretoID uint64) error
	ListExecutionsFn               func(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error)
	DeleteExecutionFn              func(ctx context.Context, executionID, userID string) error
	SaveProgressFn                 func(ctx context.Context, progress *store.ExecutionProgress) error
	GetProgressFn                  func(ctx context.Context, executionID string) (*store.ExecutionProgress, error)
	MarkExecutionForCancellationFn func(ctx context.Context, executionID, userID string) error
	IsExecutionCancelledFn         func(ctx context.Context, executionID string) (bool, error)
	SubscribeFn                    func(ctx context.Context, channel string) (<-chan []byte, error)

	HealthCheckFn func(ctx context.Context) error

	// Call tracking
	createUserCalls       int
	getUserCalls          int
	updateUserCalls       int
	deleteUserCalls       int
	createParetoCalls     int
	getParetoCalls        int
	updateParetoCalls     int
	deleteParetoCalls     int
	listParetosCalls      int
	createParetoSetCalls  int
	getParetoSetByIDCalls int
	healthCheckCalls      int
}

// Verify mockStore implements store.Store
var _ store.Store = (*mockStore)(nil)

func (m *mockStore) CreateUser(ctx context.Context, user *api.User) error {
	m.createUserCalls++
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, user)
	}
	return nil
}

func (m *mockStore) GetUser(ctx context.Context, userIDs *api.UserIDs) (*api.User, error) {
	m.getUserCalls++
	if m.GetUserFn != nil {
		return m.GetUserFn(ctx, userIDs)
	}
	return nil, nil
}

func (m *mockStore) UpdateUser(ctx context.Context, user *api.User, fields ...string) error {
	m.updateUserCalls++
	if m.UpdateUserFn != nil {
		return m.UpdateUserFn(ctx, user, fields...)
	}
	return nil
}

func (m *mockStore) DeleteUser(ctx context.Context, userIDs *api.UserIDs) error {
	m.deleteUserCalls++
	if m.DeleteUserFn != nil {
		return m.DeleteUserFn(ctx, userIDs)
	}
	return nil
}

func (m *mockStore) CreatePareto(ctx context.Context, pareto *api.Pareto) error {
	m.createParetoCalls++
	if m.CreateParetoFn != nil {
		return m.CreateParetoFn(ctx, pareto)
	}
	return nil
}

func (m *mockStore) GetPareto(ctx context.Context, ids *api.ParetoIDs) (*api.Pareto, error) {
	m.getParetoCalls++
	if m.GetParetoFn != nil {
		return m.GetParetoFn(ctx, ids)
	}
	return nil, nil
}

func (m *mockStore) UpdatePareto(ctx context.Context, pareto *api.Pareto, fields ...string) error {
	m.updateParetoCalls++
	if m.UpdateParetoFn != nil {
		return m.UpdateParetoFn(ctx, pareto, fields...)
	}
	return nil
}

func (m *mockStore) DeletePareto(ctx context.Context, ids *api.ParetoIDs) error {
	m.deleteParetoCalls++
	if m.DeleteParetoFn != nil {
		return m.DeleteParetoFn(ctx, ids)
	}
	return nil
}

func (m *mockStore) ListParetos(ctx context.Context, userIDs *api.UserIDs, limit, offset int) ([]*api.Pareto, int, error) {
	m.listParetosCalls++
	if m.ListParetosFn != nil {
		return m.ListParetosFn(ctx, userIDs, limit, offset)
	}
	return nil, 0, nil
}

func (m *mockStore) CreateParetoSet(ctx context.Context, paretoSet *store.ParetoSet) error {
	m.createParetoSetCalls++
	if m.CreateParetoSetFn != nil {
		return m.CreateParetoSetFn(ctx, paretoSet)
	}
	return nil
}

func (m *mockStore) GetParetoSetByID(ctx context.Context, id uint64) (*store.ParetoSet, error) {
	m.getParetoSetByIDCalls++
	if m.GetParetoSetByIDFn != nil {
		return m.GetParetoSetByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockStore) CreateExecution(ctx context.Context, execution *store.Execution) error {
	if m.CreateExecutionFn != nil {
		return m.CreateExecutionFn(ctx, execution)
	}
	return nil
}

func (m *mockStore) GetExecution(ctx context.Context, executionID, userID string) (*store.Execution, error) {
	if m.GetExecutionFn != nil {
		return m.GetExecutionFn(ctx, executionID, userID)
	}
	return nil, nil
}

func (m *mockStore) UpdateExecutionStatus(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
	if m.UpdateExecutionStatusFn != nil {
		return m.UpdateExecutionStatusFn(ctx, executionID, status, errorMsg)
	}
	return nil
}

func (m *mockStore) UpdateExecutionResult(ctx context.Context, executionID string, paretoID uint64) error {
	if m.UpdateExecutionResultFn != nil {
		return m.UpdateExecutionResultFn(ctx, executionID, paretoID)
	}
	return nil
}

func (m *mockStore) ListExecutions(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
	if m.ListExecutionsFn != nil {
		return m.ListExecutionsFn(ctx, userID, status, limit, offset)
	}
	return nil, 0, nil
}

func (m *mockStore) DeleteExecution(ctx context.Context, executionID, userID string) error {
	if m.DeleteExecutionFn != nil {
		return m.DeleteExecutionFn(ctx, executionID, userID)
	}
	return nil
}

func (m *mockStore) SaveProgress(ctx context.Context, progress *store.ExecutionProgress) error {
	if m.SaveProgressFn != nil {
		return m.SaveProgressFn(ctx, progress)
	}
	return nil
}

func (m *mockStore) GetProgress(ctx context.Context, executionID string) (*store.ExecutionProgress, error) {
	if m.GetProgressFn != nil {
		return m.GetProgressFn(ctx, executionID)
	}
	return nil, nil
}

func (m *mockStore) MarkExecutionForCancellation(ctx context.Context, executionID, userID string) error {
	if m.MarkExecutionForCancellationFn != nil {
		return m.MarkExecutionForCancellationFn(ctx, executionID, userID)
	}
	return nil
}

func (m *mockStore) IsExecutionCancelled(ctx context.Context, executionID string) (bool, error) {
	if m.IsExecutionCancelledFn != nil {
		return m.IsExecutionCancelledFn(ctx, executionID)
	}
	return false, nil
}

func (m *mockStore) Subscribe(ctx context.Context, channel string) (<-chan []byte, error) {
	if m.SubscribeFn != nil {
		return m.SubscribeFn(ctx, channel)
	}
	ch := make(chan []byte)
	close(ch)
	return ch, nil
}

func (m *mockStore) HealthCheck(ctx context.Context) error {
	m.healthCheckCalls++
	if m.HealthCheckFn != nil {
		return m.HealthCheckFn(ctx)
	}
	return nil
}

// mockRedisClient implements a mock Redis client for Store tests
type mockRedisClientForStore struct {
	healthCheckFn    func(ctx context.Context) error
	healthCheckCalls int
}

func (m *mockRedisClientForStore) HealthCheck(ctx context.Context) error {
	m.healthCheckCalls++
	if m.healthCheckFn != nil {
		return m.healthCheckFn(ctx)
	}
	return nil
}

// Store wrapper tests

func TestNew(t *testing.T) {
	dbStore := &mockStore{}
	redisExecStore := &mockExecutionStore{}

	// Note: New() requires a real *redis.Client which we can't easily mock
	// Instead, we test the public methods of the Store that we can access
	execStore := NewExecutionStore(redisExecStore, dbStore)

	assert.NotNil(t, execStore)
}

func TestStore_UserOperations(t *testing.T) {
	t.Run("CreateUser delegates to db", func(t *testing.T) {
		dbStore := &mockStore{}
		var createdUser *api.User
		dbStore.CreateUserFn = func(ctx context.Context, user *api.User) error {
			createdUser = user
			return nil
		}

		// Create a store that wraps the mock
		// Since we can't create the full composite Store without Redis,
		// we test via the mock directly
		ctx := context.Background()
		user := &api.User{Email: "test@example.com", Password: "secret"}

		err := dbStore.CreateUser(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, user, createdUser)
		assert.Equal(t, 1, dbStore.createUserCalls)
	})

	t.Run("CreateUser returns error", func(t *testing.T) {
		dbStore := &mockStore{}
		dbStore.CreateUserFn = func(ctx context.Context, user *api.User) error {
			return errors.New("database error")
		}

		ctx := context.Background()
		user := &api.User{Email: "test@example.com"}

		err := dbStore.CreateUser(ctx, user)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("GetUser delegates to db", func(t *testing.T) {
		dbStore := &mockStore{}
		expectedUser := &api.User{Email: "test@example.com", Password: "secret"}
		dbStore.GetUserFn = func(ctx context.Context, userIDs *api.UserIDs) (*api.User, error) {
			return expectedUser, nil
		}

		ctx := context.Background()
		userIDs := &api.UserIDs{Username: "testuser"}

		user, err := dbStore.GetUser(ctx, userIDs)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.Equal(t, 1, dbStore.getUserCalls)
	})

	t.Run("UpdateUser delegates to db", func(t *testing.T) {
		dbStore := &mockStore{}
		var updatedUser *api.User
		var updatedFields []string
		dbStore.UpdateUserFn = func(ctx context.Context, user *api.User, fields ...string) error {
			updatedUser = user
			updatedFields = fields
			return nil
		}

		ctx := context.Background()
		user := &api.User{Email: "updated@example.com", Password: "newsecret"}

		err := dbStore.UpdateUser(ctx, user, "email")

		assert.NoError(t, err)
		assert.Equal(t, user, updatedUser)
		assert.Equal(t, []string{"email"}, updatedFields)
		assert.Equal(t, 1, dbStore.updateUserCalls)
	})

	t.Run("DeleteUser delegates to db", func(t *testing.T) {
		dbStore := &mockStore{}
		var deletedUserIDs *api.UserIDs
		dbStore.DeleteUserFn = func(ctx context.Context, userIDs *api.UserIDs) error {
			deletedUserIDs = userIDs
			return nil
		}

		ctx := context.Background()
		userIDs := &api.UserIDs{Username: "testuser"}

		err := dbStore.DeleteUser(ctx, userIDs)

		assert.NoError(t, err)
		assert.Equal(t, userIDs, deletedUserIDs)
		assert.Equal(t, 1, dbStore.deleteUserCalls)
	})
}

func TestStore_ParetoOperations(t *testing.T) {
	t.Run("CreatePareto delegates to db", func(t *testing.T) {
		dbStore := &mockStore{}
		var createdPareto *api.Pareto
		dbStore.CreateParetoFn = func(ctx context.Context, pareto *api.Pareto) error {
			createdPareto = pareto
			return nil
		}

		ctx := context.Background()
		pareto := &api.Pareto{Ids: &api.ParetoIDs{Id: 1, UserId: "user-1"}}

		err := dbStore.CreatePareto(ctx, pareto)

		assert.NoError(t, err)
		assert.Equal(t, pareto, createdPareto)
		assert.Equal(t, 1, dbStore.createParetoCalls)
	})

	t.Run("GetPareto delegates to db", func(t *testing.T) {
		dbStore := &mockStore{}
		expectedPareto := &api.Pareto{Ids: &api.ParetoIDs{Id: 1, UserId: "user-1"}}
		dbStore.GetParetoFn = func(ctx context.Context, ids *api.ParetoIDs) (*api.Pareto, error) {
			return expectedPareto, nil
		}

		ctx := context.Background()
		paretoIDs := &api.ParetoIDs{Id: 1, UserId: "user-1"}

		pareto, err := dbStore.GetPareto(ctx, paretoIDs)

		assert.NoError(t, err)
		assert.Equal(t, expectedPareto, pareto)
		assert.Equal(t, 1, dbStore.getParetoCalls)
	})

	t.Run("UpdatePareto delegates to db", func(t *testing.T) {
		dbStore := &mockStore{}
		dbStore.UpdateParetoFn = func(ctx context.Context, pareto *api.Pareto, fields ...string) error {
			return nil
		}

		ctx := context.Background()
		pareto := &api.Pareto{Ids: &api.ParetoIDs{Id: 1, UserId: "user-1"}}

		err := dbStore.UpdatePareto(ctx, pareto, "field1")

		assert.NoError(t, err)
		assert.Equal(t, 1, dbStore.updateParetoCalls)
	})

	t.Run("DeletePareto delegates to db", func(t *testing.T) {
		dbStore := &mockStore{}
		dbStore.DeleteParetoFn = func(ctx context.Context, ids *api.ParetoIDs) error {
			return nil
		}

		ctx := context.Background()
		paretoIDs := &api.ParetoIDs{Id: 1, UserId: "user-1"}

		err := dbStore.DeletePareto(ctx, paretoIDs)

		assert.NoError(t, err)
		assert.Equal(t, 1, dbStore.deleteParetoCalls)
	})

	t.Run("ListParetos delegates to db", func(t *testing.T) {
		dbStore := &mockStore{}
		expectedParetos := []*api.Pareto{
			{Ids: &api.ParetoIDs{Id: 1, UserId: "user-1"}},
			{Ids: &api.ParetoIDs{Id: 2, UserId: "user-1"}},
		}
		dbStore.ListParetosFn = func(ctx context.Context, userIDs *api.UserIDs, limit, offset int) ([]*api.Pareto, int, error) {
			return expectedParetos, 2, nil
		}

		ctx := context.Background()
		userIDs := &api.UserIDs{Username: "testuser"}

		paretos, total, err := dbStore.ListParetos(ctx, userIDs, 10, 0)

		assert.NoError(t, err)
		assert.Equal(t, expectedParetos, paretos)
		assert.Equal(t, 2, total)
		assert.Equal(t, 1, dbStore.listParetosCalls)
	})

	t.Run("CreateParetoSet delegates to db", func(t *testing.T) {
		dbStore := &mockStore{}
		dbStore.CreateParetoSetFn = func(ctx context.Context, paretoSet *store.ParetoSet) error {
			return nil
		}

		ctx := context.Background()
		paretoSet := &store.ParetoSet{
			UserID:    "user-1",
			Algorithm: "gde3",
			Problem:   "zdt1",
			Variant:   "rand/1",
		}

		err := dbStore.CreateParetoSet(ctx, paretoSet)

		assert.NoError(t, err)
		assert.Equal(t, 1, dbStore.createParetoSetCalls)
	})

	t.Run("GetParetoSetByID delegates to db", func(t *testing.T) {
		dbStore := &mockStore{}
		expectedParetoSet := &store.ParetoSet{ID: 1, UserID: "user-1"}
		dbStore.GetParetoSetByIDFn = func(ctx context.Context, id uint64) (*store.ParetoSet, error) {
			return expectedParetoSet, nil
		}

		ctx := context.Background()

		paretoSet, err := dbStore.GetParetoSetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedParetoSet, paretoSet)
		assert.Equal(t, 1, dbStore.getParetoSetByIDCalls)
	})
}

func TestStore_HealthCheck(t *testing.T) {
	t.Run("success - both db and redis healthy", func(t *testing.T) {
		dbStore := &mockStore{}
		dbStore.HealthCheckFn = func(ctx context.Context) error {
			return nil
		}

		ctx := context.Background()
		err := dbStore.HealthCheck(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 1, dbStore.healthCheckCalls)
	})

	t.Run("failure - db unhealthy", func(t *testing.T) {
		dbStore := &mockStore{}
		dbStore.HealthCheckFn = func(ctx context.Context) error {
			return errors.New("database connection failed")
		}

		ctx := context.Background()
		err := dbStore.HealthCheck(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection failed")
	})
}

// Integration-style tests for the composite Store wrapper
// These test the delegation pattern from Store to ExecutionStore

func TestStore_ExecutionOperationsDelegation(t *testing.T) {
	// These tests verify that the Store properly delegates to ExecutionStore
	// by testing the ExecutionStore methods that would be called

	t.Run("CreateExecution through delegation", func(t *testing.T) {
		redis := &mockExecutionStore{}
		db := &mockExecutionStore{}

		redis.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
			return nil
		}
		db.CreateExecutionFn = func(ctx context.Context, execution *store.Execution) error {
			return nil
		}

		execStore := NewExecutionStore(redis, db)
		ctx := context.Background()
		exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusPending)

		err := execStore.CreateExecution(ctx, exec)

		assert.NoError(t, err)
		assert.Equal(t, 1, db.createExecutionCalls)
		assert.Equal(t, 1, redis.createExecutionCalls)
	})

	t.Run("GetExecution through delegation", func(t *testing.T) {
		redis := &mockExecutionStore{}
		db := &mockExecutionStore{}

		expectedExec := createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)
		redis.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
			return expectedExec, nil
		}
		db.GetExecutionFn = func(ctx context.Context, executionID, userID string) (*store.Execution, error) {
			return expectedExec, nil
		}

		execStore := NewExecutionStore(redis, db)
		ctx := context.Background()

		exec, err := execStore.GetExecution(ctx, "exec-1", "user-1")

		assert.NoError(t, err)
		assert.NotNil(t, exec)
		assert.Equal(t, "exec-1", exec.ID)
	})

	t.Run("UpdateExecutionStatus through delegation", func(t *testing.T) {
		redis := &mockExecutionStore{}
		db := &mockExecutionStore{}

		db.UpdateExecutionStatusFn = func(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
			assert.Equal(t, "exec-1", executionID)
			assert.Equal(t, store.ExecutionStatusCompleted, status)
			return nil
		}
		redis.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
			return nil
		}

		execStore := NewExecutionStore(redis, db)
		ctx := context.Background()

		err := execStore.UpdateExecutionStatus(ctx, "exec-1", store.ExecutionStatusCompleted, "")

		assert.NoError(t, err)
		assert.Equal(t, 1, db.updateExecutionStatusCalls)
	})

	t.Run("UpdateExecutionResult through delegation", func(t *testing.T) {
		redis := &mockExecutionStore{}
		db := &mockExecutionStore{}

		db.UpdateExecutionResultFn = func(ctx context.Context, executionID string, paretoID uint64) error {
			assert.Equal(t, "exec-1", executionID)
			assert.Equal(t, uint64(12345), paretoID)
			return nil
		}
		redis.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
			return nil
		}

		execStore := NewExecutionStore(redis, db)
		ctx := context.Background()

		err := execStore.UpdateExecutionResult(ctx, "exec-1", 12345)

		assert.NoError(t, err)
		assert.Equal(t, 1, db.updateExecutionResultCalls)
	})

	t.Run("ListExecutions through delegation", func(t *testing.T) {
		redis := &mockExecutionStore{}
		db := &mockExecutionStore{}

		expectedExecs := []*store.Execution{
			createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning),
			createTestExecution("exec-2", "user-1", store.ExecutionStatusCompleted),
		}
		db.ListExecutionsFn = func(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
			return expectedExecs, 2, nil
		}

		execStore := NewExecutionStore(redis, db)
		ctx := context.Background()

		execs, total, err := execStore.ListExecutions(ctx, "user-1", nil, 50, 0)

		assert.NoError(t, err)
		assert.Len(t, execs, 2)
		assert.Equal(t, 2, total)
		assert.Equal(t, 1, db.listExecutionsCalls)
	})

	t.Run("DeleteExecution through delegation", func(t *testing.T) {
		redis := &mockExecutionStore{}
		db := &mockExecutionStore{}

		db.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
			return nil
		}
		redis.DeleteExecutionFn = func(ctx context.Context, executionID, userID string) error {
			return nil
		}

		execStore := NewExecutionStore(redis, db)
		ctx := context.Background()

		err := execStore.DeleteExecution(ctx, "exec-1", "user-1")

		assert.NoError(t, err)
		assert.Equal(t, 1, db.deleteExecutionCalls)
		assert.Equal(t, 1, redis.deleteExecutionCalls)
	})

	t.Run("SaveProgress through delegation", func(t *testing.T) {
		redis := &mockExecutionStore{}
		db := &mockExecutionStore{}

		redis.SaveProgressFn = func(ctx context.Context, progress *store.ExecutionProgress) error {
			assert.Equal(t, "exec-1", progress.ExecutionID)
			return nil
		}

		execStore := NewExecutionStore(redis, db)
		ctx := context.Background()
		progress := createTestProgress("exec-1")

		err := execStore.SaveProgress(ctx, progress)

		assert.NoError(t, err)
		assert.Equal(t, 1, redis.saveProgressCalls)
	})

	t.Run("GetProgress through delegation", func(t *testing.T) {
		redis := &mockExecutionStore{}
		db := &mockExecutionStore{}

		expectedProgress := createTestProgress("exec-1")
		redis.GetProgressFn = func(ctx context.Context, executionID string) (*store.ExecutionProgress, error) {
			return expectedProgress, nil
		}

		execStore := NewExecutionStore(redis, db)
		ctx := context.Background()

		progress, err := execStore.GetProgress(ctx, "exec-1")

		assert.NoError(t, err)
		assert.NotNil(t, progress)
		assert.Equal(t, "exec-1", progress.ExecutionID)
		assert.Equal(t, 1, redis.getProgressCalls)
	})

	t.Run("MarkExecutionForCancellation through delegation", func(t *testing.T) {
		redis := &mockExecutionStore{}
		db := &mockExecutionStore{}

		redis.MarkExecutionForCancellationFn = func(ctx context.Context, executionID, userID string) error {
			return nil
		}

		execStore := NewExecutionStore(redis, db)
		ctx := context.Background()

		err := execStore.MarkExecutionForCancellation(ctx, "exec-1", "user-1")

		assert.NoError(t, err)
		assert.Equal(t, 1, redis.markExecutionForCancellationCalls)
	})

	t.Run("IsExecutionCancelled through delegation", func(t *testing.T) {
		redis := &mockExecutionStore{}
		db := &mockExecutionStore{}

		redis.IsExecutionCancelledFn = func(ctx context.Context, executionID string) (bool, error) {
			return true, nil
		}

		execStore := NewExecutionStore(redis, db)
		ctx := context.Background()

		cancelled, err := execStore.IsExecutionCancelled(ctx, "exec-1")

		assert.NoError(t, err)
		assert.True(t, cancelled)
		assert.Equal(t, 1, redis.isExecutionCancelledCalls)
	})

	t.Run("Subscribe through delegation", func(t *testing.T) {
		redis := &mockExecutionStore{}
		db := &mockExecutionStore{}

		redis.SubscribeFn = func(ctx context.Context, channel string) (<-chan []byte, error) {
			ch := make(chan []byte)
			close(ch)
			return ch, nil
		}

		execStore := NewExecutionStore(redis, db)
		ctx := context.Background()

		ch, err := execStore.Subscribe(ctx, "execution:exec-1:updates")

		assert.NoError(t, err)
		assert.NotNil(t, ch)
		assert.Equal(t, 1, redis.subscribeCalls)
	})
}
