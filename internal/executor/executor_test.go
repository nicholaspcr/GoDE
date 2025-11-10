package executor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	_ "github.com/nicholaspcr/GoDE/pkg/problems/many/dtlz" // Register DTLZ problems
	_ "github.com/nicholaspcr/GoDE/pkg/problems/many/wfg"  // Register WFG problems
	_ "github.com/nicholaspcr/GoDE/pkg/problems/multi"     // Register ZDT and VNT problems
	"github.com/nicholaspcr/GoDE/pkg/variants"
	_ "github.com/nicholaspcr/GoDE/pkg/variants/best"            // Register best/* variants
	_ "github.com/nicholaspcr/GoDE/pkg/variants/current-to-best" // Register current-to-best/* variants
	_ "github.com/nicholaspcr/GoDE/pkg/variants/pbest"           // Register pbest/* variants
	_ "github.com/nicholaspcr/GoDE/pkg/variants/rand"            // Register rand/* variants

	"github.com/nicholaspcr/GoDE/internal/store"
)

// mockStore implements a minimal store.Store for testing
type mockStore struct {
	executions map[string]*store.Execution
	progress   map[string]*store.ExecutionProgress
	paretoSets map[uint64]*store.ParetoSet
	nextID     uint64
}

func newMockStore() *mockStore {
	return &mockStore{
		executions: make(map[string]*store.Execution),
		progress:   make(map[string]*store.ExecutionProgress),
		paretoSets: make(map[uint64]*store.ParetoSet),
		nextID:     1,
	}
}

func (m *mockStore) CreateExecution(ctx context.Context, execution *store.Execution) error {
	m.executions[execution.ID] = execution
	return nil
}

func (m *mockStore) GetExecution(ctx context.Context, executionID, userID string) (*store.Execution, error) {
	exec, exists := m.executions[executionID]
	if !exists || exec.UserID != userID {
		return nil, store.ErrExecutionNotFound
	}
	return exec, nil
}

func (m *mockStore) UpdateExecutionStatus(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
	exec, exists := m.executions[executionID]
	if !exists {
		return store.ErrExecutionNotFound
	}
	exec.Status = status
	exec.Error = errorMsg
	exec.UpdatedAt = time.Now()
	if status == store.ExecutionStatusCompleted || status == store.ExecutionStatusFailed || status == store.ExecutionStatusCancelled {
		now := time.Now()
		exec.CompletedAt = &now
	}
	return nil
}

func (m *mockStore) UpdateExecutionResult(ctx context.Context, executionID string, paretoID uint64) error {
	exec, exists := m.executions[executionID]
	if !exists {
		return store.ErrExecutionNotFound
	}
	exec.ParetoID = &paretoID
	return nil
}

func (m *mockStore) ListExecutions(ctx context.Context, userID string, status *store.ExecutionStatus) ([]*store.Execution, error) {
	var result []*store.Execution
	for _, exec := range m.executions {
		if exec.UserID == userID {
			if status == nil || exec.Status == *status {
				result = append(result, exec)
			}
		}
	}
	return result, nil
}

func (m *mockStore) DeleteExecution(ctx context.Context, executionID, userID string) error {
	exec, exists := m.executions[executionID]
	if !exists || exec.UserID != userID {
		return store.ErrExecutionNotFound
	}
	delete(m.executions, executionID)
	return nil
}

func (m *mockStore) SaveProgress(ctx context.Context, progress *store.ExecutionProgress) error {
	m.progress[progress.ExecutionID] = progress
	return nil
}

func (m *mockStore) GetProgress(ctx context.Context, executionID string) (*store.ExecutionProgress, error) {
	prog, exists := m.progress[executionID]
	if !exists {
		return nil, store.ErrExecutionNotFound
	}
	return prog, nil
}

func (m *mockStore) MarkExecutionForCancellation(ctx context.Context, executionID, userID string) error {
	exec, exists := m.executions[executionID]
	if !exists || exec.UserID != userID {
		return store.ErrExecutionNotFound
	}
	return nil
}

func (m *mockStore) IsExecutionCancelled(ctx context.Context, executionID string) (bool, error) {
	return false, nil
}

func (m *mockStore) Subscribe(ctx context.Context, channel string) (<-chan []byte, error) {
	ch := make(chan []byte)
	close(ch)
	return ch, nil
}

func (m *mockStore) CreateParetoSet(ctx context.Context, paretoSet *store.ParetoSet) error {
	paretoSet.ID = m.nextID
	m.paretoSets[m.nextID] = paretoSet
	m.nextID++
	return nil
}

func (m *mockStore) GetParetoSetByID(ctx context.Context, id uint64) (*store.ParetoSet, error) {
	ps, exists := m.paretoSets[id]
	if !exists {
		return nil, store.ErrParetoSetNotFound
	}
	return ps, nil
}

// Stub implementations for other store methods
func (m *mockStore) CreateUser(ctx context.Context, user *api.User) error                         { return nil }
func (m *mockStore) GetUser(ctx context.Context, ids *api.UserIDs) (*api.User, error)            { return nil, nil }
func (m *mockStore) UpdateUser(ctx context.Context, user *api.User, fields ...string) error      { return nil }
func (m *mockStore) DeleteUser(ctx context.Context, ids *api.UserIDs) error                      { return nil }
func (m *mockStore) CreatePareto(ctx context.Context, pareto *api.Pareto) error                  { return nil }
func (m *mockStore) GetPareto(ctx context.Context, ids *api.ParetoIDs) (*api.Pareto, error)      { return nil, nil }
func (m *mockStore) UpdatePareto(ctx context.Context, pareto *api.Pareto, fields ...string) error { return nil }
func (m *mockStore) DeletePareto(ctx context.Context, ids *api.ParetoIDs) error                  { return nil }
func (m *mockStore) ListParetos(ctx context.Context, ids *api.UserIDs) ([]*api.Pareto, error)    { return nil, nil }
func (m *mockStore) HealthCheck(ctx context.Context) error                                       { return nil }

func TestExecutor_SubmitExecution(t *testing.T) {
	mockSt := newMockStore()
	exec := New(Config{
		Store:        mockSt,
		MaxWorkers:   2,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	// Register a test problem and variant
	prob, err := problems.DefaultRegistry.Create("zdt1", 10, 2)
	require.NoError(t, err)
	exec.RegisterProblem("zdt1", prob)

	variant, err := variants.DefaultRegistry.Create("rand1")
	require.NoError(t, err)
	exec.RegisterVariant("rand1", variant)

	ctx := context.Background()
	userID := "test-user"
	algorithm := "gde3"
	problem := "zdt1"
	variantName := "rand1"

	config := &api.DEConfig{
		Executions:     1,
		Generations:    2,
		PopulationSize: 10,
		DimensionsSize: 10,
		ObjetivesSize:  2,
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

	executionID, err := exec.SubmitExecution(ctx, userID, algorithm, problem, variantName, config)
	require.NoError(t, err)
	assert.NotEmpty(t, executionID)

	// Verify execution was created
	execution, err := mockSt.GetExecution(ctx, executionID, userID)
	require.NoError(t, err)
	assert.Equal(t, userID, execution.UserID)
	assert.Equal(t, store.ExecutionStatusPending, execution.Status)
}

func TestExecutor_CancelExecution(t *testing.T) {
	mockSt := newMockStore()
	exec := New(Config{
		Store:        mockSt,
		MaxWorkers:   2,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	ctx := context.Background()
	userID := "test-user"
	executionID := "test-exec-123"

	// Create a mock execution
	_ = mockSt.CreateExecution(ctx, &store.Execution{
		ID:        executionID,
		UserID:    userID,
		Status:    store.ExecutionStatusRunning,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	// Cancel the execution
	err := exec.CancelExecution(ctx, executionID, userID)
	require.NoError(t, err)
}

func TestExecutor_CancelExecution_NotFound(t *testing.T) {
	mockSt := newMockStore()
	exec := New(Config{
		Store:        mockSt,
		MaxWorkers:   2,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	ctx := context.Background()
	err := exec.CancelExecution(ctx, "non-existent", "test-user")
	assert.Error(t, err)
}

func TestExecutor_RegisterProblemAndVariant(t *testing.T) {
	mockSt := newMockStore()
	exec := New(Config{
		Store:        mockSt,
		MaxWorkers:   2,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	// Register problem
	prob, err := problems.DefaultRegistry.Create("zdt1", 10, 2)
	require.NoError(t, err)
	exec.RegisterProblem("zdt1", prob)

	// Register variant
	variant, err := variants.DefaultRegistry.Create("rand1")
	require.NoError(t, err)
	exec.RegisterVariant("rand1", variant)

	// Verify they're registered by checking the internal maps
	assert.NotNil(t, exec.problemRegistry["zdt1"])
	assert.NotNil(t, exec.variantRegistry["rand1"])
}
