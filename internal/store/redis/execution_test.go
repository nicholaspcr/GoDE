package redis

import (
	"context"
	"encoding/json"
	"errors"
	"maps"
	"sync"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/nicholaspcr/GoDE/internal/cache/redis"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRedisClient implements redis.ClientInterface for unit testing
type mockRedisClient struct {
	mu       sync.RWMutex
	data     map[string]string
	hashData map[string]map[string]string
	expiries map[string]time.Time
	pubsub   map[string][]chan string

	// Error injection
	setErr     error
	getErr     error
	deleteErr  error
	hsetErr    error
	hgetAllErr error
	hscanErr   error
	hdelErr    error
	expireErr  error
	publishErr error
}

// Verify mockRedisClient implements redis.ClientInterface
var _ redis.ClientInterface = (*mockRedisClient)(nil)

func newMockRedisClient() *mockRedisClient {
	return &mockRedisClient{
		data:     make(map[string]string),
		hashData: make(map[string]map[string]string),
		expiries: make(map[string]time.Time),
		pubsub:   make(map[string][]chan string),
	}
}

func (m *mockRedisClient) Get(_ context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.getErr != nil {
		return "", m.getErr
	}

	// Check expiry
	if exp, ok := m.expiries[key]; ok && time.Now().After(exp) {
		delete(m.data, key)
		delete(m.expiries, key)
		return "", goredis.Nil
	}

	val, ok := m.data[key]
	if !ok {
		return "", goredis.Nil
	}
	return val, nil
}

func (m *mockRedisClient) Set(_ context.Context, key string, value any, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.setErr != nil {
		return m.setErr
	}

	var strVal string
	switch v := value.(type) {
	case string:
		strVal = v
	case []byte:
		strVal = string(v)
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		strVal = string(data)
	}

	m.data[key] = strVal
	if ttl > 0 {
		m.expiries[key] = time.Now().Add(ttl)
	}
	return nil
}

func (m *mockRedisClient) Delete(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.deleteErr != nil {
		return m.deleteErr
	}

	delete(m.data, key)
	delete(m.expiries, key)
	return nil
}

func (m *mockRedisClient) HSet(_ context.Context, key string, values ...any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.hsetErr != nil {
		return m.hsetErr
	}

	if m.hashData[key] == nil {
		m.hashData[key] = make(map[string]string)
	}

	// Process field-value pairs
	for i := 0; i < len(values); i += 2 {
		if i+1 >= len(values) {
			break
		}
		field, ok := values[i].(string)
		if !ok {
			continue
		}
		var val string
		switch v := values[i+1].(type) {
		case string:
			val = v
		case []byte:
			val = string(v)
		default:
			data, _ := json.Marshal(v)
			val = string(data)
		}
		m.hashData[key][field] = val
	}
	return nil
}

func (m *mockRedisClient) HGetAll(_ context.Context, key string) (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.hgetAllErr != nil {
		return nil, m.hgetAllErr
	}

	result := make(map[string]string)
	if hash, ok := m.hashData[key]; ok {
		maps.Copy(result, hash)
	}
	return result, nil
}

func (m *mockRedisClient) HGet(_ context.Context, key, field string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if hash, ok := m.hashData[key]; ok {
		if val, ok := hash[field]; ok {
			return val, nil
		}
	}
	return "", goredis.Nil
}

func (m *mockRedisClient) HScan(_ context.Context, key string, cursor uint64, _ string, _ int64) ([]string, uint64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.hscanErr != nil {
		return nil, 0, m.hscanErr
	}

	// Simple implementation: return all fields on cursor=0, then return empty with cursor=0
	if cursor != 0 {
		return []string{}, 0, nil
	}

	var pairs []string
	if hash, ok := m.hashData[key]; ok {
		for field, value := range hash {
			pairs = append(pairs, field, value)
		}
	}
	return pairs, 0, nil // cursor=0 means scan complete
}

func (m *mockRedisClient) HDel(_ context.Context, key string, fields ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.hdelErr != nil {
		return m.hdelErr
	}

	if hash, ok := m.hashData[key]; ok {
		for _, field := range fields {
			delete(hash, field)
		}
	}
	return nil
}

func (m *mockRedisClient) HLen(_ context.Context, key string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if hash, ok := m.hashData[key]; ok {
		return int64(len(hash)), nil
	}
	return 0, nil
}

func (m *mockRedisClient) Expire(_ context.Context, key string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.expireErr != nil {
		return m.expireErr
	}

	m.expiries[key] = time.Now().Add(ttl)
	return nil
}

func (m *mockRedisClient) Publish(_ context.Context, channel string, message any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.publishErr != nil {
		return m.publishErr
	}

	// Notify subscribers (non-blocking)
	if subs, ok := m.pubsub[channel]; ok {
		var msgStr string
		switch v := message.(type) {
		case string:
			msgStr = v
		case []byte:
			msgStr = string(v)
		default:
			data, _ := json.Marshal(v)
			msgStr = string(data)
		}
		for _, ch := range subs {
			select {
			case ch <- msgStr:
			default:
			}
		}
	}
	return nil
}

func (m *mockRedisClient) Subscribe(_ context.Context, _ string) *goredis.PubSub {
	// Return nil - Subscribe functionality tested separately via integration tests
	return nil
}

// Test helper functions

func createTestExecution(id, userID string, status store.ExecutionStatus) *store.Execution {
	return &store.Execution{
		ID:        id,
		UserID:    userID,
		Status:    status,
		Config:    createTestDEConfig(),
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
	}
}

// Unit Tests

func TestNewExecutionStore(t *testing.T) {
	mock := newMockRedisClient()
	executionTTL := 24 * time.Hour
	progressTTL := time.Hour

	s := NewExecutionStore(mock, executionTTL, progressTTL)

	assert.NotNil(t, s)
	assert.Equal(t, executionTTL, s.executionTTL)
	assert.Equal(t, progressTTL, s.progressTTL)
}

func TestExecutionStore_CreateExecution(t *testing.T) {
	tests := []struct {
		name      string
		execution *store.Execution
		setupMock func(*mockRedisClient)
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "create execution successfully",
			execution: createTestExecution("exec-1", "user-1", store.ExecutionStatusPending),
			setupMock: func(_ *mockRedisClient) {},
			wantErr:   false,
		},
		{
			name:      "create execution with nil config",
			execution: &store.Execution{ID: "exec-2", UserID: "user-1", Status: store.ExecutionStatusPending},
			setupMock: func(_ *mockRedisClient) {},
			wantErr:   false,
		},
		{
			name:      "create execution fails on set error",
			execution: createTestExecution("exec-3", "user-1", store.ExecutionStatusPending),
			setupMock: func(m *mockRedisClient) {
				m.setErr = errors.New("redis connection error")
			},
			wantErr: true,
			errMsg:  "failed to store execution",
		},
		{
			name:      "create execution fails on hset error",
			execution: createTestExecution("exec-4", "user-1", store.ExecutionStatusPending),
			setupMock: func(m *mockRedisClient) {
				m.hsetErr = errors.New("hset error")
			},
			wantErr: true,
			errMsg:  "failed to add execution to user set",
		},
		{
			name:      "create execution fails on expire error",
			execution: createTestExecution("exec-5", "user-1", store.ExecutionStatusPending),
			setupMock: func(m *mockRedisClient) {
				m.expireErr = errors.New("expire error")
			},
			wantErr: true,
			errMsg:  "failed to set TTL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockRedisClient()
			tt.setupMock(mock)

			s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
			ctx := context.Background()

			err := s.CreateExecution(ctx, tt.execution)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)

				// Verify execution was stored
				retrieved, err := s.GetExecution(ctx, tt.execution.ID, tt.execution.UserID)
				require.NoError(t, err)
				assert.Equal(t, tt.execution.ID, retrieved.ID)
				assert.Equal(t, tt.execution.UserID, retrieved.UserID)
				assert.Equal(t, tt.execution.Status, retrieved.Status)
			}
		})
	}
}

func TestExecutionStore_GetExecution(t *testing.T) {
	tests := []struct {
		name        string
		executionID string
		userID      string
		setup       func(*ExecutionStore)
		wantErr     bool
		errContains string
	}{
		{
			name:        "get existing execution",
			executionID: "exec-1",
			userID:      "user-1",
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)
				_ = s.CreateExecution(context.Background(), exec)
			},
			wantErr: false,
		},
		{
			name:        "get non-existent execution",
			executionID: "exec-nonexistent",
			userID:      "user-1",
			setup:       func(_ *ExecutionStore) {},
			wantErr:     true,
			errContains: "execution not found",
		},
		{
			name:        "get execution with wrong user",
			executionID: "exec-1",
			userID:      "user-2",
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)
				_ = s.CreateExecution(context.Background(), exec)
			},
			wantErr:     true,
			errContains: "does not belong to user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockRedisClient()
			s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
			ctx := context.Background()

			tt.setup(s)

			execution, err := s.GetExecution(ctx, tt.executionID, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, execution)
				assert.Equal(t, tt.executionID, execution.ID)
				assert.Equal(t, tt.userID, execution.UserID)
			}
		})
	}
}

func TestExecutionStore_GetExecution_WithGetError(t *testing.T) {
	mock := newMockRedisClient()
	mock.getErr = errors.New("redis get error")
	s := NewExecutionStore(mock, 24*time.Hour, time.Hour)

	_, err := s.GetExecution(context.Background(), "exec-1", "user-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "execution not found")
}

func TestExecutionStore_UpdateExecutionStatus(t *testing.T) {
	tests := []struct {
		name        string
		executionID string
		newStatus   store.ExecutionStatus
		errorMsg    string
		setup       func(*ExecutionStore)
		wantErr     bool
		verify      func(*testing.T, *ExecutionStore)
	}{
		{
			name:        "update status to running",
			executionID: "exec-1",
			newStatus:   store.ExecutionStatusRunning,
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusPending)
				_ = s.CreateExecution(context.Background(), exec)
			},
			wantErr: false,
			verify: func(t *testing.T, s *ExecutionStore) {
				exec, err := s.GetExecution(context.Background(), "exec-1", "user-1")
				require.NoError(t, err)
				assert.Equal(t, store.ExecutionStatusRunning, exec.Status)
				assert.Nil(t, exec.CompletedAt)
			},
		},
		{
			name:        "update status to completed sets CompletedAt",
			executionID: "exec-1",
			newStatus:   store.ExecutionStatusCompleted,
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)
				_ = s.CreateExecution(context.Background(), exec)
			},
			wantErr: false,
			verify: func(t *testing.T, s *ExecutionStore) {
				exec, err := s.GetExecution(context.Background(), "exec-1", "user-1")
				require.NoError(t, err)
				assert.Equal(t, store.ExecutionStatusCompleted, exec.Status)
				assert.NotNil(t, exec.CompletedAt)
			},
		},
		{
			name:        "update status to failed with error message",
			executionID: "exec-1",
			newStatus:   store.ExecutionStatusFailed,
			errorMsg:    "algorithm failed: out of memory",
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)
				_ = s.CreateExecution(context.Background(), exec)
			},
			wantErr: false,
			verify: func(t *testing.T, s *ExecutionStore) {
				exec, err := s.GetExecution(context.Background(), "exec-1", "user-1")
				require.NoError(t, err)
				assert.Equal(t, store.ExecutionStatusFailed, exec.Status)
				assert.Equal(t, "algorithm failed: out of memory", exec.Error)
				assert.NotNil(t, exec.CompletedAt)
			},
		},
		{
			name:        "update status to cancelled sets CompletedAt",
			executionID: "exec-1",
			newStatus:   store.ExecutionStatusCancelled,
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)
				_ = s.CreateExecution(context.Background(), exec)
			},
			wantErr: false,
			verify: func(t *testing.T, s *ExecutionStore) {
				exec, err := s.GetExecution(context.Background(), "exec-1", "user-1")
				require.NoError(t, err)
				assert.Equal(t, store.ExecutionStatusCancelled, exec.Status)
				assert.NotNil(t, exec.CompletedAt)
			},
		},
		{
			name:        "update non-existent execution fails",
			executionID: "exec-nonexistent",
			newStatus:   store.ExecutionStatusRunning,
			setup:       func(_ *ExecutionStore) {},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockRedisClient()
			s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
			ctx := context.Background()

			tt.setup(s)

			err := s.UpdateExecutionStatus(ctx, tt.executionID, tt.newStatus, tt.errorMsg)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.verify != nil {
					tt.verify(t, s)
				}
			}
		})
	}
}

func TestExecutionStore_UpdateExecutionStatus_SetError(t *testing.T) {
	mock := newMockRedisClient()
	s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
	ctx := context.Background()

	exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusPending)
	err := s.CreateExecution(ctx, exec)
	require.NoError(t, err)

	// Now inject error for the update
	mock.setErr = errors.New("redis set error")

	err = s.UpdateExecutionStatus(ctx, "exec-1", store.ExecutionStatusRunning, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update execution")
}

func TestExecutionStore_UpdateExecutionStatus_HSetError(t *testing.T) {
	mock := newMockRedisClient()
	s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
	ctx := context.Background()

	exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusPending)
	err := s.CreateExecution(ctx, exec)
	require.NoError(t, err)

	// Now inject error for the update
	mock.hsetErr = errors.New("redis hset error")

	err = s.UpdateExecutionStatus(ctx, "exec-1", store.ExecutionStatusRunning, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update execution in user set")
}

func TestExecutionStore_UpdateExecutionResult(t *testing.T) {
	tests := []struct {
		name        string
		executionID string
		paretoID    uint64
		setup       func(*ExecutionStore)
		wantErr     bool
		verify      func(*testing.T, *ExecutionStore)
	}{
		{
			name:        "update execution result successfully",
			executionID: "exec-1",
			paretoID:    12345,
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusCompleted)
				_ = s.CreateExecution(context.Background(), exec)
			},
			wantErr: false,
			verify: func(t *testing.T, s *ExecutionStore) {
				exec, err := s.GetExecution(context.Background(), "exec-1", "user-1")
				require.NoError(t, err)
				require.NotNil(t, exec.ParetoID)
				assert.Equal(t, uint64(12345), *exec.ParetoID)
			},
		},
		{
			name:        "update non-existent execution fails",
			executionID: "exec-nonexistent",
			paretoID:    12345,
			setup:       func(_ *ExecutionStore) {},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockRedisClient()
			s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
			ctx := context.Background()

			tt.setup(s)

			err := s.UpdateExecutionResult(ctx, tt.executionID, tt.paretoID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.verify != nil {
					tt.verify(t, s)
				}
			}
		})
	}
}

func TestExecutionStore_UpdateExecutionResult_SetError(t *testing.T) {
	mock := newMockRedisClient()
	s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
	ctx := context.Background()

	exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusCompleted)
	err := s.CreateExecution(ctx, exec)
	require.NoError(t, err)

	// Now inject error for the update
	mock.setErr = errors.New("redis set error")

	err = s.UpdateExecutionResult(ctx, "exec-1", 12345)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update execution")
}

func TestExecutionStore_UpdateExecutionResult_HSetError(t *testing.T) {
	mock := newMockRedisClient()
	s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
	ctx := context.Background()

	exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusCompleted)
	err := s.CreateExecution(ctx, exec)
	require.NoError(t, err)

	// Now inject error for the update
	mock.hsetErr = errors.New("redis hset error")

	err = s.UpdateExecutionResult(ctx, "exec-1", 12345)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update execution in user set")
}

func TestExecutionStore_ListExecutions(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		status    *store.ExecutionStatus
		limit     int
		offset    int
		setup     func(*ExecutionStore)
		wantCount int
		wantTotal int
		wantErr   bool
	}{
		{
			name:   "list all executions for user",
			userID: "user-1",
			limit:  50,
			offset: 0,
			setup: func(s *ExecutionStore) {
				for i := range 5 {
					exec := createTestExecution("exec-"+string(rune('1'+i)), "user-1", store.ExecutionStatusPending)
					_ = s.CreateExecution(context.Background(), exec)
				}
			},
			wantCount: 5,
			wantTotal: 5,
			wantErr:   false,
		},
		{
			name:   "list executions with status filter",
			userID: "user-1",
			status: func() *store.ExecutionStatus { s := store.ExecutionStatusRunning; return &s }(),
			limit:  50,
			offset: 0,
			setup: func(s *ExecutionStore) {
				exec1 := createTestExecution("exec-1", "user-1", store.ExecutionStatusPending)
				exec2 := createTestExecution("exec-2", "user-1", store.ExecutionStatusRunning)
				exec3 := createTestExecution("exec-3", "user-1", store.ExecutionStatusRunning)
				_ = s.CreateExecution(context.Background(), exec1)
				_ = s.CreateExecution(context.Background(), exec2)
				_ = s.CreateExecution(context.Background(), exec3)
			},
			wantCount: 2,
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name:   "list executions with pagination",
			userID: "user-1",
			limit:  2,
			offset: 1,
			setup: func(s *ExecutionStore) {
				for i := range 5 {
					exec := createTestExecution("exec-"+string(rune('1'+i)), "user-1", store.ExecutionStatusPending)
					_ = s.CreateExecution(context.Background(), exec)
				}
			},
			wantCount: 2,
			wantTotal: 5,
			wantErr:   false,
		},
		{
			name:   "list executions with invalid limit defaults to 50",
			userID: "user-1",
			limit:  -1,
			offset: 0,
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusPending)
				_ = s.CreateExecution(context.Background(), exec)
			},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name:   "list executions with limit > 100 defaults to 50",
			userID: "user-1",
			limit:  200,
			offset: 0,
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusPending)
				_ = s.CreateExecution(context.Background(), exec)
			},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name:      "list executions for user with no executions",
			userID:    "user-empty",
			limit:     50,
			offset:    0,
			setup:     func(_ *ExecutionStore) {},
			wantCount: 0,
			wantTotal: 0,
			wantErr:   false,
		},
		{
			name:   "list executions with negative offset treats as 0",
			userID: "user-1",
			limit:  50,
			offset: -5,
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusPending)
				_ = s.CreateExecution(context.Background(), exec)
			},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockRedisClient()
			s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
			ctx := context.Background()

			tt.setup(s)

			executions, total, err := s.ListExecutions(ctx, tt.userID, tt.status, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, executions, tt.wantCount)
				assert.Equal(t, tt.wantTotal, total)
			}
		})
	}
}

func TestExecutionStore_ListExecutions_HScanError(t *testing.T) {
	mock := newMockRedisClient()
	mock.hscanErr = errors.New("hscan error")
	s := NewExecutionStore(mock, 24*time.Hour, time.Hour)

	_, _, err := s.ListExecutions(context.Background(), "user-1", nil, 50, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to scan user executions")
}

func TestExecutionStore_DeleteExecution(t *testing.T) {
	tests := []struct {
		name        string
		executionID string
		userID      string
		setup       func(*ExecutionStore)
		wantErr     bool
		errContains string
		verify      func(*testing.T, *ExecutionStore)
	}{
		{
			name:        "delete existing execution",
			executionID: "exec-1",
			userID:      "user-1",
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusCompleted)
				_ = s.CreateExecution(context.Background(), exec)
			},
			wantErr: false,
			verify: func(t *testing.T, s *ExecutionStore) {
				_, err := s.GetExecution(context.Background(), "exec-1", "user-1")
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "execution not found")
			},
		},
		{
			name:        "delete execution with empty userID (cache invalidation)",
			executionID: "exec-1",
			userID:      "",
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusCompleted)
				_ = s.CreateExecution(context.Background(), exec)
			},
			wantErr: false,
		},
		{
			name:        "delete execution with wrong user fails",
			executionID: "exec-1",
			userID:      "user-2",
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusCompleted)
				_ = s.CreateExecution(context.Background(), exec)
			},
			wantErr:     true,
			errContains: "does not belong to user",
		},
		{
			name:        "delete non-existent execution fails",
			executionID: "exec-nonexistent",
			userID:      "user-1",
			setup:       func(_ *ExecutionStore) {},
			wantErr:     true,
			errContains: "execution not found",
		},
		{
			name:        "delete cleans up progress and cancel keys",
			executionID: "exec-1",
			userID:      "user-1",
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)
				_ = s.CreateExecution(context.Background(), exec)
				progress := createTestProgress("exec-1")
				_ = s.SaveProgress(context.Background(), progress)
				_ = s.MarkExecutionForCancellation(context.Background(), "exec-1", "user-1")
			},
			wantErr: false,
			verify: func(t *testing.T, s *ExecutionStore) {
				// Execution should be gone
				_, err := s.GetExecution(context.Background(), "exec-1", "user-1")
				assert.Error(t, err)

				// Progress should be gone
				_, err = s.GetProgress(context.Background(), "exec-1")
				assert.Error(t, err)

				// Cancellation flag should be gone
				cancelled, err := s.IsExecutionCancelled(context.Background(), "exec-1")
				assert.NoError(t, err)
				assert.False(t, cancelled)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockRedisClient()
			s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
			ctx := context.Background()

			tt.setup(s)

			err := s.DeleteExecution(ctx, tt.executionID, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				if tt.verify != nil {
					tt.verify(t, s)
				}
			}
		})
	}
}

func TestExecutionStore_DeleteExecution_DeleteError(t *testing.T) {
	mock := newMockRedisClient()
	s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
	ctx := context.Background()

	exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusCompleted)
	err := s.CreateExecution(ctx, exec)
	require.NoError(t, err)

	// Now inject error for the delete
	mock.deleteErr = errors.New("redis delete error")

	err = s.DeleteExecution(ctx, "exec-1", "user-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete execution")
}

func TestExecutionStore_SaveProgress(t *testing.T) {
	tests := []struct {
		name      string
		progress  *store.ExecutionProgress
		setupMock func(*mockRedisClient)
		wantErr   bool
	}{
		{
			name:      "save progress successfully",
			progress:  createTestProgress("exec-1"),
			setupMock: func(_ *mockRedisClient) {},
			wantErr:   false,
		},
		{
			name: "save progress with partial pareto",
			progress: &store.ExecutionProgress{
				ExecutionID:       "exec-1",
				CurrentGeneration: 25,
				TotalGenerations:  100,
				PartialPareto: []*api.Vector{
					{Elements: []float64{0.1, 0.2}, Objectives: []float64{0.5, 0.5}},
				},
			},
			setupMock: func(_ *mockRedisClient) {},
			wantErr:   false,
		},
		{
			name:     "save progress fails on set error",
			progress: createTestProgress("exec-1"),
			setupMock: func(m *mockRedisClient) {
				m.setErr = errors.New("redis error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockRedisClient()
			tt.setupMock(mock)

			s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
			ctx := context.Background()

			err := s.SaveProgress(ctx, tt.progress)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify progress was stored
				retrieved, err := s.GetProgress(ctx, tt.progress.ExecutionID)
				require.NoError(t, err)
				assert.Equal(t, tt.progress.ExecutionID, retrieved.ExecutionID)
				assert.Equal(t, tt.progress.CurrentGeneration, retrieved.CurrentGeneration)
			}
		})
	}
}

func TestExecutionStore_SaveProgress_PublishError(t *testing.T) {
	mock := newMockRedisClient()
	mock.publishErr = errors.New("publish error")
	s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
	ctx := context.Background()

	progress := createTestProgress("exec-1")
	err := s.SaveProgress(ctx, progress)
	// Should succeed even if publish fails (publish is non-critical)
	assert.NoError(t, err)
}

func TestExecutionStore_GetProgress(t *testing.T) {
	tests := []struct {
		name        string
		executionID string
		setup       func(*ExecutionStore)
		wantErr     bool
		errContains string
	}{
		{
			name:        "get existing progress",
			executionID: "exec-1",
			setup: func(s *ExecutionStore) {
				progress := createTestProgress("exec-1")
				_ = s.SaveProgress(context.Background(), progress)
			},
			wantErr: false,
		},
		{
			name:        "get non-existent progress",
			executionID: "exec-nonexistent",
			setup:       func(_ *ExecutionStore) {},
			wantErr:     true,
			errContains: "progress not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockRedisClient()
			s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
			ctx := context.Background()

			tt.setup(s)

			progress, err := s.GetProgress(ctx, tt.executionID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, progress)
				assert.Equal(t, tt.executionID, progress.ExecutionID)
			}
		})
	}
}

func TestExecutionStore_GetProgress_InvalidJSON(t *testing.T) {
	mock := newMockRedisClient()
	// Manually set invalid JSON in the mock
	mock.data["execution:exec-1:progress"] = "not valid json"

	s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
	ctx := context.Background()

	_, err := s.GetProgress(ctx, "exec-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal progress")
}

func TestExecutionStore_MarkExecutionForCancellation(t *testing.T) {
	tests := []struct {
		name        string
		executionID string
		userID      string
		setup       func(*ExecutionStore)
		wantErr     bool
		errContains string
	}{
		{
			name:        "mark execution for cancellation successfully",
			executionID: "exec-1",
			userID:      "user-1",
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)
				_ = s.CreateExecution(context.Background(), exec)
			},
			wantErr: false,
		},
		{
			name:        "mark non-existent execution fails",
			executionID: "exec-nonexistent",
			userID:      "user-1",
			setup:       func(_ *ExecutionStore) {},
			wantErr:     true,
			errContains: "execution not found",
		},
		{
			name:        "mark execution with wrong user fails",
			executionID: "exec-1",
			userID:      "user-2",
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)
				_ = s.CreateExecution(context.Background(), exec)
			},
			wantErr:     true,
			errContains: "does not belong to user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockRedisClient()
			s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
			ctx := context.Background()

			tt.setup(s)

			err := s.MarkExecutionForCancellation(ctx, tt.executionID, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)

				// Verify cancellation was set
				cancelled, err := s.IsExecutionCancelled(ctx, tt.executionID)
				require.NoError(t, err)
				assert.True(t, cancelled)
			}
		})
	}
}

func TestExecutionStore_MarkExecutionForCancellation_SetError(t *testing.T) {
	mock := newMockRedisClient()
	s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
	ctx := context.Background()

	exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)
	err := s.CreateExecution(ctx, exec)
	require.NoError(t, err)

	// Now inject error for the set
	mock.setErr = errors.New("redis set error")

	err = s.MarkExecutionForCancellation(ctx, "exec-1", "user-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to set cancellation flag")
}

func TestExecutionStore_IsExecutionCancelled(t *testing.T) {
	tests := []struct {
		name        string
		executionID string
		setup       func(*ExecutionStore)
		want        bool
	}{
		{
			name:        "execution is cancelled",
			executionID: "exec-1",
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)
				_ = s.CreateExecution(context.Background(), exec)
				_ = s.MarkExecutionForCancellation(context.Background(), "exec-1", "user-1")
			},
			want: true,
		},
		{
			name:        "execution is not cancelled",
			executionID: "exec-1",
			setup: func(s *ExecutionStore) {
				exec := createTestExecution("exec-1", "user-1", store.ExecutionStatusRunning)
				_ = s.CreateExecution(context.Background(), exec)
			},
			want: false,
		},
		{
			name:        "non-existent execution is not cancelled",
			executionID: "exec-nonexistent",
			setup:       func(_ *ExecutionStore) {},
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockRedisClient()
			s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
			ctx := context.Background()

			tt.setup(s)

			cancelled, err := s.IsExecutionCancelled(ctx, tt.executionID)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, cancelled)
		})
	}
}

func TestExecutionStore_Subscribe(t *testing.T) {
	mock := newMockRedisClient()
	s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Note: This test only verifies that Subscribe returns a channel
	// Full pub/sub functionality requires integration testing with real Redis
	ch, err := s.Subscribe(ctx, "execution:exec-1:updates")

	// With mock client returning nil PubSub, the implementation will handle this gracefully
	assert.NoError(t, err)
	assert.NotNil(t, ch)
}

// Tests for marshaling/unmarshaling functions

func TestMarshalUnmarshalExecution(t *testing.T) {
	tests := []struct {
		name      string
		execution *store.Execution
	}{
		{
			name:      "execution with all fields",
			execution: createTestExecution("exec-1", "user-1", store.ExecutionStatusCompleted),
		},
		{
			name: "execution with nil config",
			execution: &store.Execution{
				ID:        "exec-2",
				UserID:    "user-1",
				Status:    store.ExecutionStatusPending,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			name: "execution with pareto ID",
			execution: func() *store.Execution {
				exec := createTestExecution("exec-3", "user-1", store.ExecutionStatusCompleted)
				paretoID := uint64(12345)
				exec.ParetoID = &paretoID
				return exec
			}(),
		},
		{
			name: "execution with error",
			execution: func() *store.Execution {
				exec := createTestExecution("exec-4", "user-1", store.ExecutionStatusFailed)
				exec.Error = "algorithm failed"
				return exec
			}(),
		},
		{
			name: "execution with completed at",
			execution: func() *store.Execution {
				exec := createTestExecution("exec-5", "user-1", store.ExecutionStatusCompleted)
				now := time.Now()
				exec.CompletedAt = &now
				return exec
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := marshalExecution(tt.execution)
			require.NoError(t, err)
			assert.NotEmpty(t, data)

			unmarshaled, err := unmarshalExecution(data)
			require.NoError(t, err)

			assert.Equal(t, tt.execution.ID, unmarshaled.ID)
			assert.Equal(t, tt.execution.UserID, unmarshaled.UserID)
			assert.Equal(t, tt.execution.Status, unmarshaled.Status)
			assert.Equal(t, tt.execution.Error, unmarshaled.Error)

			// Compare pointers
			if tt.execution.ParetoID != nil {
				require.NotNil(t, unmarshaled.ParetoID)
				assert.Equal(t, *tt.execution.ParetoID, *unmarshaled.ParetoID)
			} else {
				assert.Nil(t, unmarshaled.ParetoID)
			}

			// Compare config
			if tt.execution.Config != nil {
				require.NotNil(t, unmarshaled.Config)
				assert.Equal(t, tt.execution.Config.Generations, unmarshaled.Config.Generations)
			} else {
				assert.Nil(t, unmarshaled.Config)
			}
		})
	}
}

func TestUnmarshalExecution_InvalidJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:    "invalid json",
			data:    []byte("not json"),
			wantErr: true,
		},
		{
			name:    "valid json but wrong structure",
			data:    []byte(`{"foo": "bar"}`),
			wantErr: false, // Will unmarshal with zero values
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := unmarshalExecution(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUnmarshalExecution_InvalidConfigJSON(t *testing.T) {
	// Create valid wrapper but with invalid config JSON
	data := []byte(`{"id":"exec-1","user_id":"user-1","status":"pending","config_json":"not valid protojson"}`)
	_, err := unmarshalExecution(data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal config")
}

// Test ConfigToProto function

func TestConfigToProto(t *testing.T) {
	tests := []struct {
		name   string
		config *api.DEConfig
		want   map[string]any
	}{
		{
			name: "config with GDE3",
			config: &api.DEConfig{
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
			},
			want: map[string]any{
				"executions":      int64(1),
				"generations":     int64(100),
				"population_size": int64(50),
				"dimensions_size": int64(10),
				"objectives_size": int64(2),
				"floor_limiter":   float32(0.0),
				"ceil_limiter":    float32(1.0),
				"gde3_config": map[string]any{
					"cr": float32(0.9),
					"f":  float32(0.5),
					"p":  float32(0.1),
				},
			},
		},
		{
			name: "config without algorithm config",
			config: &api.DEConfig{
				Executions:     2,
				Generations:    50,
				PopulationSize: 30,
				DimensionsSize: 5,
				ObjectivesSize: 3,
				FloorLimiter:   -1.0,
				CeilLimiter:    1.0,
			},
			want: map[string]any{
				"executions":      int64(2),
				"generations":     int64(50),
				"population_size": int64(30),
				"dimensions_size": int64(5),
				"objectives_size": int64(3),
				"floor_limiter":   float32(-1.0),
				"ceil_limiter":    float32(1.0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConfigToProto(tt.config)

			assert.Equal(t, tt.want["executions"], result["executions"])
			assert.Equal(t, tt.want["generations"], result["generations"])
			assert.Equal(t, tt.want["population_size"], result["population_size"])
			assert.Equal(t, tt.want["dimensions_size"], result["dimensions_size"])
			assert.Equal(t, tt.want["objectives_size"], result["objectives_size"])
			assert.Equal(t, tt.want["floor_limiter"], result["floor_limiter"])
			assert.Equal(t, tt.want["ceil_limiter"], result["ceil_limiter"])

			if tt.want["gde3_config"] != nil {
				wantGde3 := tt.want["gde3_config"].(map[string]any)
				resultGde3, ok := result["gde3_config"].(map[string]any)
				require.True(t, ok)
				assert.Equal(t, wantGde3["cr"], resultGde3["cr"])
				assert.Equal(t, wantGde3["f"], resultGde3["f"])
				assert.Equal(t, wantGde3["p"], resultGde3["p"])
			}
		})
	}
}

// Test key generation functions

func TestKeyGeneration(t *testing.T) {
	mock := newMockRedisClient()
	s := NewExecutionStore(mock, 24*time.Hour, time.Hour)

	t.Run("executionKey", func(t *testing.T) {
		key := s.executionKey("exec-123")
		assert.Equal(t, "execution:exec-123", key)
	})

	t.Run("progressKey", func(t *testing.T) {
		key := s.progressKey("exec-123")
		assert.Equal(t, "execution:exec-123:progress", key)
	})

	t.Run("cancelKey", func(t *testing.T) {
		key := s.cancelKey("exec-123")
		assert.Equal(t, "execution:exec-123:cancel", key)
	})

	t.Run("userExecutionsKey", func(t *testing.T) {
		key := s.userExecutionsKey("user-456")
		assert.Equal(t, "user:user-456:executions", key)
	})
}

// Test edge cases

func TestExecutionStore_EdgeCases(t *testing.T) {
	t.Run("execution with very long ID", func(t *testing.T) {
		mock := newMockRedisClient()
		s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
		ctx := context.Background()

		longID := string(make([]byte, 1000))
		for i := range len(longID) {
			longID = longID[:i] + "a" + longID[i+1:]
		}

		exec := &store.Execution{
			ID:        longID,
			UserID:    "user-1",
			Status:    store.ExecutionStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := s.CreateExecution(ctx, exec)
		assert.NoError(t, err)

		retrieved, err := s.GetExecution(ctx, longID, "user-1")
		require.NoError(t, err)
		assert.Equal(t, longID, retrieved.ID)
	})

	t.Run("execution with special characters in ID", func(t *testing.T) {
		mock := newMockRedisClient()
		s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
		ctx := context.Background()

		specialID := "exec:with:colons:and-dashes_underscores"

		exec := &store.Execution{
			ID:        specialID,
			UserID:    "user-1",
			Status:    store.ExecutionStatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := s.CreateExecution(ctx, exec)
		assert.NoError(t, err)

		retrieved, err := s.GetExecution(ctx, specialID, "user-1")
		require.NoError(t, err)
		assert.Equal(t, specialID, retrieved.ID)
	})

	t.Run("concurrent execution creation", func(t *testing.T) {
		mock := newMockRedisClient()
		s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
		ctx := context.Background()

		var wg sync.WaitGroup
		numGoroutines := 10

		for i := range numGoroutines {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				exec := createTestExecution("exec-concurrent-"+string(rune('0'+idx)), "user-1", store.ExecutionStatusPending)
				err := s.CreateExecution(ctx, exec)
				assert.NoError(t, err)
			}(i)
		}

		wg.Wait()

		// Verify all executions were created
		executions, total, err := s.ListExecutions(ctx, "user-1", nil, 50, 0)
		require.NoError(t, err)
		assert.Equal(t, numGoroutines, total)
		assert.Len(t, executions, numGoroutines)
	})

	t.Run("rapid status updates", func(t *testing.T) {
		mock := newMockRedisClient()
		s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
		ctx := context.Background()

		exec := createTestExecution("exec-rapid", "user-1", store.ExecutionStatusPending)
		err := s.CreateExecution(ctx, exec)
		require.NoError(t, err)

		statuses := []store.ExecutionStatus{
			store.ExecutionStatusRunning,
			store.ExecutionStatusCompleted,
		}

		for _, status := range statuses {
			err := s.UpdateExecutionStatus(ctx, "exec-rapid", status, "")
			assert.NoError(t, err)
		}

		retrieved, err := s.GetExecution(ctx, "exec-rapid", "user-1")
		require.NoError(t, err)
		assert.Equal(t, store.ExecutionStatusCompleted, retrieved.Status)
	})
}

// Benchmark tests

func BenchmarkExecutionStore_CreateExecution(b *testing.B) {
	mock := newMockRedisClient()
	s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		exec := createTestExecution("exec-bench-"+string(rune(i)), "user-1", store.ExecutionStatusPending)
		_ = s.CreateExecution(ctx, exec)
	}
}

func BenchmarkExecutionStore_GetExecution(b *testing.B) {
	mock := newMockRedisClient()
	s := NewExecutionStore(mock, 24*time.Hour, time.Hour)
	ctx := context.Background()

	exec := createTestExecution("exec-bench", "user-1", store.ExecutionStatusPending)
	_ = s.CreateExecution(ctx, exec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.GetExecution(ctx, "exec-bench", "user-1")
	}
}

func BenchmarkMarshalExecution(b *testing.B) {
	exec := createTestExecution("exec-bench", "user-1", store.ExecutionStatusPending)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = marshalExecution(exec)
	}
}

func BenchmarkUnmarshalExecution(b *testing.B) {
	exec := createTestExecution("exec-bench", "user-1", store.ExecutionStatusPending)
	data, _ := marshalExecution(exec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = unmarshalExecution(data)
	}
}
