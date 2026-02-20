package executor

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/models"
	_ "github.com/nicholaspcr/GoDE/pkg/de/gde3"            // Register GDE3 algorithm factory
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
	mu         sync.RWMutex
}

// Deep copy helpers for mockStore
//
// These functions create deep copies of store types to prevent data races.
// This mimics the behavior of real stores (Redis, GORM) which serialize/
// deserialize data, creating implicit copies. Without deep copying, the mock
// would share pointers between algorithm goroutines and test goroutines,
// causing data races when both access the same struct/slice concurrently.

// deepCopyVector creates a deep copy of api.Vector
func deepCopyVector(src *api.Vector) *api.Vector {
	if src == nil {
		return nil
	}

	dst := &api.Vector{
		CrowdingDistance: src.CrowdingDistance,
	}

	// Copy Elements slice
	if src.Elements != nil {
		dst.Elements = make([]float64, len(src.Elements))
		copy(dst.Elements, src.Elements)
	}

	// Copy Objectives slice
	if src.Objectives != nil {
		dst.Objectives = make([]float64, len(src.Objectives))
		copy(dst.Objectives, src.Objectives)
	}

	// Copy Ids if present (may be nil)
	if src.Ids != nil {
		dst.Ids = &api.VectorIDs{
			Id: src.Ids.Id,
		}
	}

	return dst
}

// deepCopyProgress creates a deep copy of ExecutionProgress to prevent data races.
// This mimics the behavior of real stores (Redis, GORM) which serialize/deserialize.
func deepCopyProgress(src *store.ExecutionProgress) *store.ExecutionProgress {
	if src == nil {
		return nil
	}

	dst := &store.ExecutionProgress{
		ExecutionID:         src.ExecutionID,
		CurrentGeneration:   src.CurrentGeneration,
		TotalGenerations:    src.TotalGenerations,
		CompletedExecutions: src.CompletedExecutions,
		TotalExecutions:     src.TotalExecutions,
		UpdatedAt:           src.UpdatedAt,
	}

	// Deep copy PartialPareto slice
	if src.PartialPareto != nil {
		dst.PartialPareto = make([]*api.Vector, len(src.PartialPareto))
		for i, vec := range src.PartialPareto {
			dst.PartialPareto[i] = deepCopyVector(vec)
		}
	}

	return dst
}

// deepCopyExecution creates a deep copy of Execution
func deepCopyExecution(src *store.Execution) *store.Execution {
	if src == nil {
		return nil
	}

	dst := &store.Execution{
		ID:                  src.ID,
		UserID:              src.UserID,
		Status:              src.Status,
		Error:               src.Error,
		Algorithm:           src.Algorithm,
		Variant:             src.Variant,
		Problem:             src.Problem,
		IdempotencyKey:      src.IdempotencyKey,
		MaxExecutionSeconds: src.MaxExecutionSeconds,
		CreatedAt:           src.CreatedAt,
		UpdatedAt:           src.UpdatedAt,
	}

	// Copy pointer fields
	if src.ParetoID != nil {
		paretoID := *src.ParetoID
		dst.ParetoID = &paretoID
	}

	if src.CompletedAt != nil {
		completedAt := *src.CompletedAt
		dst.CompletedAt = &completedAt
	}

	// Config is typically read-only after creation, share the pointer
	dst.Config = src.Config

	return dst
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
	m.mu.Lock()
	defer m.mu.Unlock()
	m.executions[execution.ID] = deepCopyExecution(execution)
	return nil
}

func (m *mockStore) GetExecution(ctx context.Context, executionID, userID string) (*store.Execution, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	exec, exists := m.executions[executionID]
	if !exists || exec.UserID != userID {
		return nil, store.ErrExecutionNotFound
	}
	return deepCopyExecution(exec), nil
}

func (m *mockStore) UpdateExecutionStatus(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	exec, exists := m.executions[executionID]
	if !exists {
		return store.ErrExecutionNotFound
	}
	// Create a copy, modify it, store the copy
	updated := deepCopyExecution(exec)
	updated.Status = status
	updated.Error = errorMsg
	updated.UpdatedAt = time.Now()
	if status == store.ExecutionStatusCompleted || status == store.ExecutionStatusFailed || status == store.ExecutionStatusCancelled {
		now := time.Now()
		updated.CompletedAt = &now
	}
	m.executions[executionID] = updated
	return nil
}

func (m *mockStore) UpdateExecutionResult(ctx context.Context, executionID string, paretoID uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	exec, exists := m.executions[executionID]
	if !exists {
		return store.ErrExecutionNotFound
	}
	updated := deepCopyExecution(exec)
	updated.ParetoID = &paretoID
	m.executions[executionID] = updated
	return nil
}

func (m *mockStore) ListExecutions(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var allMatching []*store.Execution
	for _, exec := range m.executions {
		if exec.UserID == userID {
			if status == nil || exec.Status == *status {
				// Deep copy each matching execution
				allMatching = append(allMatching, deepCopyExecution(exec))
			}
		}
	}

	// Apply defaults
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	totalCount := len(allMatching)

	// Apply pagination
	start := offset
	if start > totalCount {
		return []*store.Execution{}, totalCount, nil
	}
	end := min(start+limit, totalCount)

	return allMatching[start:end], totalCount, nil
}

func (m *mockStore) DeleteExecution(ctx context.Context, executionID, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	exec, exists := m.executions[executionID]
	if !exists || exec.UserID != userID {
		return store.ErrExecutionNotFound
	}
	delete(m.executions, executionID)
	return nil
}

func (m *mockStore) SaveProgress(ctx context.Context, progress *store.ExecutionProgress) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Deep copy to prevent data races - mimics Redis/GORM serialization behavior
	m.progress[progress.ExecutionID] = deepCopyProgress(progress)
	return nil
}

func (m *mockStore) GetProgress(ctx context.Context, executionID string) (*store.ExecutionProgress, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	prog, exists := m.progress[executionID]
	if !exists {
		return nil, store.ErrExecutionNotFound
	}
	// Deep copy to prevent data races - mimics Redis/GORM serialization behavior
	return deepCopyProgress(prog), nil
}

func (m *mockStore) MarkExecutionForCancellation(ctx context.Context, executionID, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	exec, exists := m.executions[executionID]
	if !exists || exec.UserID != userID {
		return store.ErrExecutionNotFound
	}
	return nil
}

func (m *mockStore) IsExecutionCancelled(ctx context.Context, executionID string) (bool, error) {
	return false, nil
}

func (m *mockStore) GetExecutionByIdempotencyKey(_ context.Context, _, _ string) (string, error) {
	return "", nil
}

func (m *mockStore) Subscribe(ctx context.Context, channel string) (<-chan []byte, error) {
	ch := make(chan []byte)
	close(ch)
	return ch, nil
}

func (m *mockStore) CreateParetoSet(ctx context.Context, paretoSet *store.ParetoSet) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	paretoSet.ID = m.nextID
	m.paretoSets[m.nextID] = paretoSet
	m.nextID++
	return nil
}

func (m *mockStore) GetParetoSetByID(ctx context.Context, id uint64) (*store.ParetoSet, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ps, exists := m.paretoSets[id]
	if !exists {
		return nil, store.ErrParetoSetNotFound
	}
	return ps, nil
}

// Stub implementations for other store methods
func (m *mockStore) CreateUser(ctx context.Context, user *api.User) error { return nil }
func (m *mockStore) GetUser(ctx context.Context, ids *api.UserIDs) (*api.User, error) {
	return nil, nil
}
func (m *mockStore) UpdateUser(ctx context.Context, user *api.User, fields ...string) error {
	return nil
}
func (m *mockStore) DeleteUser(ctx context.Context, ids *api.UserIDs) error     { return nil }
func (m *mockStore) CreatePareto(ctx context.Context, pareto *api.Pareto) error { return nil }
func (m *mockStore) GetPareto(ctx context.Context, ids *api.ParetoIDs) (*api.Pareto, error) {
	return nil, nil
}
func (m *mockStore) UpdatePareto(ctx context.Context, pareto *api.Pareto, fields ...string) error {
	return nil
}
func (m *mockStore) DeletePareto(ctx context.Context, ids *api.ParetoIDs) error { return nil }
func (m *mockStore) ListParetos(ctx context.Context, ids *api.UserIDs, limit, offset int) ([]*api.Pareto, int, error) {
	return nil, 0, nil
}
func (m *mockStore) HealthCheck(ctx context.Context) error { return nil }

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

	executionID, err := exec.SubmitExecution(ctx, userID, algorithm, problem, variantName, config, "", 0)
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

// TestExecutor_WorkerPoolExhaustion tests that jobs queue when worker pool is full
func TestExecutor_WorkerPoolExhaustion(t *testing.T) {
	mockSt := newMockStore()
	exec := New(Config{
		Store:        mockSt,
		MaxWorkers:   2,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	// Register problem and variant
	prob, err := problems.DefaultRegistry.Create("zdt1", 10, 2)
	require.NoError(t, err)
	exec.RegisterProblem("zdt1", prob)

	variant, err := variants.DefaultRegistry.Create("rand1")
	require.NoError(t, err)
	exec.RegisterVariant("rand1", variant)

	ctx := context.Background()
	userID := "test-user"

	config := &api.DEConfig{
		Executions:     1,
		Generations:    5,
		PopulationSize: 10,
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

	// Submit 5 jobs with MaxWorkers=2
	// Only 2 should run concurrently, others should queue
	var executionIDs []string
	for range 5 {
		execID, err := exec.SubmitExecution(ctx, userID, "gde3", "zdt1", "rand1", config, "", 0)
		require.NoError(t, err)
		executionIDs = append(executionIDs, execID)
	}

	// Wait briefly to let workers start
	time.Sleep(100 * time.Millisecond)

	// Check that at most 2 are running concurrently
	exec.activeExecsMu.RLock()
	activeCount := len(exec.activeExecs)
	exec.activeExecsMu.RUnlock()

	assert.LessOrEqual(t, activeCount, 2, "Should not exceed MaxWorkers")

	// Wait for all to complete
	assert.Eventually(t, func() bool {
		for _, id := range executionIDs {
			execution, err := mockSt.GetExecution(ctx, id, userID)
			if err != nil || (execution.Status != store.ExecutionStatusCompleted && execution.Status != store.ExecutionStatusFailed) {
				return false
			}
		}
		return true
	}, 10*time.Second, 100*time.Millisecond, "All executions should eventually complete")
}

// TestExecutor_ConcurrentSubmissions tests race conditions in concurrent submissions
func TestExecutor_ConcurrentSubmissions(t *testing.T) {
	mockSt := newMockStore()
	exec := New(Config{
		Store:        mockSt,
		MaxWorkers:   5,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	// Register problem and variant
	prob, err := problems.DefaultRegistry.Create("zdt1", 10, 2)
	require.NoError(t, err)
	exec.RegisterProblem("zdt1", prob)

	variant, err := variants.DefaultRegistry.Create("rand1")
	require.NoError(t, err)
	exec.RegisterVariant("rand1", variant)

	ctx := context.Background()
	userID := "test-user"

	config := &api.DEConfig{
		Executions:     1,
		Generations:    2,
		PopulationSize: 10,
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

	// Submit 10 jobs concurrently from different goroutines
	var wg sync.WaitGroup
	var mu sync.Mutex
	var executionIDs []string
	var errors []error

	for range 10 {
		wg.Go(func() {
			execID, err := exec.SubmitExecution(ctx, userID, "gde3", "zdt1", "rand1", config, "", 0)
			mu.Lock()
			if err != nil {
				errors = append(errors, err)
			} else {
				executionIDs = append(executionIDs, execID)
			}
			mu.Unlock()
		})
	}

	wg.Wait()

	// All submissions should succeed
	assert.Empty(t, errors, "No submission errors expected")
	assert.Len(t, executionIDs, 10, "All 10 submissions should succeed")

	// Wait for all to complete
	assert.Eventually(t, func() bool {
		for _, id := range executionIDs {
			execution, err := mockSt.GetExecution(ctx, id, userID)
			if err != nil || (execution.Status != store.ExecutionStatusCompleted && execution.Status != store.ExecutionStatusFailed) {
				return false
			}
		}
		return true
	}, 10*time.Second, 100*time.Millisecond)
}

// TestExecutor_CompletionTracking tests that CompletedExecutions is correctly tracked
func TestExecutor_CompletionTracking(t *testing.T) {
	mockSt := newMockStore()
	exec := New(Config{
		Store:        mockSt,
		MaxWorkers:   2,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	// Register problem and variant
	prob, err := problems.DefaultRegistry.Create("zdt1", 10, 2)
	require.NoError(t, err)
	exec.RegisterProblem("zdt1", prob)

	variant, err := variants.DefaultRegistry.Create("rand1")
	require.NoError(t, err)
	exec.RegisterVariant("rand1", variant)

	ctx := context.Background()
	userID := "test-user"

	config := &api.DEConfig{
		Executions:     3,
		Generations:    2,
		PopulationSize: 10,
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

	executionID, err := exec.SubmitExecution(ctx, userID, "gde3", "zdt1", "rand1", config, "", 0)
	require.NoError(t, err)

	// Poll progress until complete
	assert.Eventually(t, func() bool {
		progress, err := mockSt.GetProgress(ctx, executionID)
		if err != nil {
			return false
		}
		// Check that CompletedExecutions increases over time
		if progress.CompletedExecutions > 0 {
			t.Logf("CompletedExecutions: %d/%d", progress.CompletedExecutions, progress.TotalExecutions)
			return progress.CompletedExecutions == progress.TotalExecutions
		}
		return false
	}, 10*time.Second, 100*time.Millisecond, "CompletedExecutions should reach TotalExecutions")
}

// TestExecutor_ActiveExecutionTracking tests that activeExecs map is correctly managed
func TestExecutor_ActiveExecutionTracking(t *testing.T) {
	mockSt := newMockStore()
	exec := New(Config{
		Store:        mockSt,
		MaxWorkers:   2,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	// Register a slow problem to ensure execution is observable
	exec.RegisterProblem("slow-problem", &slowProblem{duration: 10 * time.Millisecond})

	variant, err := variants.DefaultRegistry.Create("rand1")
	require.NoError(t, err)
	exec.RegisterVariant("rand1", variant)

	ctx := context.Background()
	userID := "test-user"

	config := &api.DEConfig{
		Executions:     1,
		Generations:    3,
		PopulationSize: 5,
		DimensionsSize: 5,
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

	executionID, err := exec.SubmitExecution(ctx, userID, "gde3", "slow-problem", "rand1", config, "", 0)
	require.NoError(t, err)

	// Should appear in activeExecs
	assert.Eventually(t, func() bool {
		exec.activeExecsMu.RLock()
		defer exec.activeExecsMu.RUnlock()
		_, exists := exec.activeExecs[executionID]
		return exists
	}, 2*time.Second, 10*time.Millisecond, "Execution should appear in activeExecs")

	// Wait for completion
	assert.Eventually(t, func() bool {
		execution, err := mockSt.GetExecution(ctx, executionID, userID)
		return err == nil && (execution.Status == store.ExecutionStatusCompleted || execution.Status == store.ExecutionStatusFailed)
	}, 10*time.Second, 100*time.Millisecond)

	// Should be removed from activeExecs
	exec.activeExecsMu.RLock()
	_, exists := exec.activeExecs[executionID]
	exec.activeExecsMu.RUnlock()
	assert.False(t, exists, "Execution should be removed from activeExecs after completion")
}

// panicProblem is a test problem that panics during evaluation
type panicProblem struct{}

func (p *panicProblem) Name() string {
	return "panic-problem"
}

func (p *panicProblem) Evaluate(vector *models.Vector, executionNumber int) error {
	panic("test panic in problem evaluation")
}

// TestExecutor_PanicRecovery tests that panics are recovered and executions marked as failed
func TestExecutor_PanicRecovery(t *testing.T) {
	mockSt := newMockStore()
	exec := New(Config{
		Store:        mockSt,
		MaxWorkers:   2,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	// Register panic problem and variant
	exec.RegisterProblem("panic-problem", &panicProblem{})

	variant, err := variants.DefaultRegistry.Create("rand1")
	require.NoError(t, err)
	exec.RegisterVariant("rand1", variant)

	ctx := context.Background()
	userID := "test-user"

	config := &api.DEConfig{
		Executions:     1,
		Generations:    2,
		PopulationSize: 10,
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

	executionID, err := exec.SubmitExecution(ctx, userID, "gde3", "panic-problem", "rand1", config, "", 0)
	require.NoError(t, err)

	// Wait for execution to fail due to panic
	assert.Eventually(t, func() bool {
		execution, err := mockSt.GetExecution(ctx, executionID, userID)
		if err != nil {
			return false
		}
		return execution.Status == store.ExecutionStatusFailed
	}, 5*time.Second, 100*time.Millisecond, "Execution should be marked as failed after panic")

	// Check error message contains panic info
	execution, err := mockSt.GetExecution(ctx, executionID, userID)
	require.NoError(t, err)
	assert.Contains(t, execution.Error, "panic", "Error should mention panic")

	// Verify execution removed from activeExecs (worker slot released)
	exec.activeExecsMu.RLock()
	_, exists := exec.activeExecs[executionID]
	exec.activeExecsMu.RUnlock()
	assert.False(t, exists, "Execution should be removed from activeExecs after panic")
}

// cancellingStore is a mock store that marks executions as cancelled immediately
type cancellingStore struct {
	*mockStore
	cancelAfter time.Duration
}

func (c *cancellingStore) IsExecutionCancelled(ctx context.Context, executionID string) (bool, error) {
	// Simulate cancellation after a delay
	time.Sleep(c.cancelAfter)
	return true, nil
}

// TestExecutor_CancellationDuringExecution tests cancelling a running execution
func TestExecutor_CancellationDuringExecution(t *testing.T) {
	mockSt := &cancellingStore{
		mockStore:   newMockStore(),
		cancelAfter: 200 * time.Millisecond,
	}

	exec := New(Config{
		Store:        mockSt,
		MaxWorkers:   2,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	// Register problem and variant
	prob, err := problems.DefaultRegistry.Create("zdt1", 10, 2)
	require.NoError(t, err)
	exec.RegisterProblem("zdt1", prob)

	variant, err := variants.DefaultRegistry.Create("rand1")
	require.NoError(t, err)
	exec.RegisterVariant("rand1", variant)

	ctx := context.Background()
	userID := "test-user"

	config := &api.DEConfig{
		Executions:     1,
		Generations:    500, // Long-running to allow cancellation
		PopulationSize: 50,
		DimensionsSize: 30,
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

	executionID, err := exec.SubmitExecution(ctx, userID, "gde3", "zdt1", "rand1", config, "", 0)
	require.NoError(t, err)

	// Wait for execution to start
	assert.Eventually(t, func() bool {
		exec.activeExecsMu.RLock()
		defer exec.activeExecsMu.RUnlock()
		_, exists := exec.activeExecs[executionID]
		return exists
	}, 2*time.Second, 10*time.Millisecond, "Execution should start")

	// Now cancel it
	err = exec.CancelExecution(ctx, executionID, userID)
	require.NoError(t, err)

	// Execution should eventually be marked as cancelled
	assert.Eventually(t, func() bool {
		execution, err := mockSt.GetExecution(ctx, executionID, userID)
		if err != nil {
			return false
		}
		return execution.Status == store.ExecutionStatusCancelled
	}, 5*time.Second, 100*time.Millisecond, "Execution should be cancelled")

	// Worker slot should be released
	exec.activeExecsMu.RLock()
	_, exists := exec.activeExecs[executionID]
	exec.activeExecsMu.RUnlock()
	assert.False(t, exists, "Execution should be removed from activeExecs after cancellation")
}

// errorProblem is a test problem that returns an error during evaluation
type errorProblem struct{}

func (e *errorProblem) Name() string {
	return "error-problem"
}

func (e *errorProblem) Evaluate(vector *models.Vector, executionNumber int) error {
	return errors.New("test evaluation error")
}

// TestExecutor_WorkerSlotReleaseOnError tests that worker slots are released on error
func TestExecutor_WorkerSlotReleaseOnError(t *testing.T) {
	mockSt := newMockStore()
	exec := New(Config{
		Store:        mockSt,
		MaxWorkers:   1,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	// Register error-prone problem
	exec.RegisterProblem("error-problem", &errorProblem{})

	variant, err := variants.DefaultRegistry.Create("rand1")
	require.NoError(t, err)
	exec.RegisterVariant("rand1", variant)

	ctx := context.Background()
	userID := "test-user"

	config := &api.DEConfig{
		Executions:     1,
		Generations:    2,
		PopulationSize: 10,
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

	// Submit execution that will fail
	executionID, err := exec.SubmitExecution(ctx, userID, "gde3", "error-problem", "rand1", config, "", 0)
	require.NoError(t, err)

	// Wait for failure
	assert.Eventually(t, func() bool {
		execution, err := mockSt.GetExecution(ctx, executionID, userID)
		if err != nil {
			return false
		}
		return execution.Status == store.ExecutionStatusFailed
	}, 5*time.Second, 100*time.Millisecond)

	// Worker slot should be released, allowing new submission
	_, err = exec.SubmitExecution(ctx, userID, "gde3", "error-problem", "rand1", config, "", 0)
	require.NoError(t, err, "Should be able to submit new execution after worker slot released")
}

// slowProblem is a test problem that sleeps for a configurable duration
type slowProblem struct {
	duration time.Duration
}

func (s *slowProblem) Name() string {
	return "slow-problem"
}

func (s *slowProblem) Evaluate(vector *models.Vector, executionNumber int) error {
	time.Sleep(s.duration)
	// Set dummy objectives
	for i := range vector.Objectives {
		vector.Objectives[i] = float64(i)
	}
	return nil
}

// TestExecutor_Shutdown_HappyPath tests successful shutdown when workers finish quickly
func TestExecutor_Shutdown_HappyPath(t *testing.T) {
	mockSt := newMockStore()
	exec := New(Config{
		Store:        mockSt,
		MaxWorkers:   2,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	// Register a fast problem
	exec.RegisterProblem("slow-problem", &slowProblem{duration: 100 * time.Millisecond})

	variant, err := variants.DefaultRegistry.Create("rand1")
	require.NoError(t, err)
	exec.RegisterVariant("rand1", variant)

	ctx := context.Background()
	userID := "test-user"

	config := &api.DEConfig{
		Executions:     2,
		Generations:    5,
		PopulationSize: 10,
		DimensionsSize: 5,
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

	// Submit multiple executions
	executionID1, err := exec.SubmitExecution(ctx, userID, "gde3", "slow-problem", "rand1", config, "", 0)
	require.NoError(t, err)
	executionID2, err := exec.SubmitExecution(ctx, userID, "gde3", "slow-problem", "rand1", config, "", 0)
	require.NoError(t, err)

	// Wait a bit for executions to start
	time.Sleep(200 * time.Millisecond)

	// Shutdown should succeed
	shutdownCtx := context.Background()
	err = exec.Shutdown(shutdownCtx)
	assert.NoError(t, err, "Shutdown should succeed when workers finish quickly")

	// Verify executions were cancelled
	for _, execID := range []string{executionID1, executionID2} {
		execution, err := mockSt.GetExecution(ctx, execID, userID)
		if err == nil {
			assert.Equal(t, store.ExecutionStatusCancelled, execution.Status,
				"Execution should be cancelled after shutdown")
		}
	}
}

// TestExecutor_Shutdown_Timeout tests that shutdown returns error when workers don't finish
func TestExecutor_Shutdown_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timeout test in short mode")
	}

	mockSt := newMockStore()
	exec := New(Config{
		Store:        mockSt,
		MaxWorkers:   1,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	// Register a very slow problem that will outlast shutdown timeout (30s)
	exec.RegisterProblem("slow-problem", &slowProblem{duration: 60 * time.Second})

	variant, err := variants.DefaultRegistry.Create("rand1")
	require.NoError(t, err)
	exec.RegisterVariant("rand1", variant)

	ctx := context.Background()
	userID := "test-user"

	config := &api.DEConfig{
		Executions:     1,
		Generations:    100, // Many generations to keep it running
		PopulationSize: 10,
		DimensionsSize: 5,
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

	// Submit execution
	_, err = exec.SubmitExecution(ctx, userID, "gde3", "slow-problem", "rand1", config, "", 0)
	require.NoError(t, err)

	// Wait for execution to start
	time.Sleep(200 * time.Millisecond)

	// Shutdown with short timeout context (5 seconds) to avoid waiting 30s
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = exec.Shutdown(shutdownCtx)
	assert.Error(t, err, "Shutdown should fail when workers don't finish in time")
	assert.ErrorIs(t, err, context.DeadlineExceeded, "Should return context deadline exceeded error")
}

// TestExecutor_Shutdown_ContextCancellation tests that shutdown aborts when context is cancelled
func TestExecutor_Shutdown_ContextCancellation(t *testing.T) {
	mockSt := newMockStore()
	exec := New(Config{
		Store:        mockSt,
		MaxWorkers:   2,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	// Register a slow problem
	exec.RegisterProblem("slow-problem", &slowProblem{duration: 5 * time.Second})

	variant, err := variants.DefaultRegistry.Create("rand1")
	require.NoError(t, err)
	exec.RegisterVariant("rand1", variant)

	ctx := context.Background()
	userID := "test-user"

	config := &api.DEConfig{
		Executions:     2,
		Generations:    20,
		PopulationSize: 10,
		DimensionsSize: 5,
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

	// Submit executions
	_, err = exec.SubmitExecution(ctx, userID, "gde3", "slow-problem", "rand1", config, "", 0)
	require.NoError(t, err)
	_, err = exec.SubmitExecution(ctx, userID, "gde3", "slow-problem", "rand1", config, "", 0)
	require.NoError(t, err)

	// Wait for executions to start
	time.Sleep(200 * time.Millisecond)

	// Create cancellable context for shutdown
	shutdownCtx, cancel := context.WithCancel(context.Background())

	// Cancel the shutdown context after a short delay
	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	// Shutdown should abort when context is cancelled
	err = exec.Shutdown(shutdownCtx)
	assert.Error(t, err, "Shutdown should fail when context is cancelled")
	assert.ErrorIs(t, err, context.Canceled, "Should return context cancelled error")
}

// TestExecutor_ExecutionTimeout tests that an execution is marked as failed when it exceeds its timeout.
// The GDE3 algorithm checks ctx.Err() at each generation boundary, so we use a short per-evaluation
// sleep combined with a tight timeout that triggers after at least one full generation completes.
func TestExecutor_ExecutionTimeout(t *testing.T) {
	mockSt := newMockStore()

	// Each evaluation sleeps 10ms; with population=5 a generation takes ~50ms.
	// The default timeout of 300ms allows a few generations then stops.
	exec := New(Config{
		Store:               mockSt,
		MaxWorkers:          2,
		ExecutionTTL:        time.Hour,
		ResultTTL:           time.Hour,
		ProgressTTL:         time.Minute,
		DefaultMaxExecution: 300 * time.Millisecond,
	})

	exec.RegisterProblem("slow-problem", &slowProblem{duration: 10 * time.Millisecond})
	variant, err := variants.DefaultRegistry.Create("rand1")
	require.NoError(t, err)
	exec.RegisterVariant("rand1", variant)

	ctx := context.Background()
	userID := "test-user"

	config := &api.DEConfig{
		Executions:     1,
		Generations:    1000, // many generations to ensure it won't finish naturally
		PopulationSize: 5,
		DimensionsSize: 5,
		ObjectivesSize: 2,
		FloorLimiter:   0.0,
		CeilLimiter:    1.0,
		AlgorithmConfig: &api.DEConfig_Gde3{
			Gde3: &api.GDE3Config{Cr: 0.9, F: 0.5, P: 0.1},
		},
	}

	executionID, err := exec.SubmitExecution(ctx, userID, "gde3", "slow-problem", "rand1", config, "", 0)
	require.NoError(t, err)

	// Wait for the execution to time out and be marked failed
	assert.Eventually(t, func() bool {
		execution, err := mockSt.GetExecution(ctx, executionID, userID)
		if err != nil {
			return false
		}
		return execution.Status == store.ExecutionStatusFailed
	}, 10*time.Second, 100*time.Millisecond, "Execution should be marked as failed after timeout")
}

// TestExecutor_ExecutionTimeout_PerRequestOverride tests that per-request max_execution_seconds overrides the default.
func TestExecutor_ExecutionTimeout_PerRequestOverride(t *testing.T) {
	mockSt := newMockStore()

	// Default is generous; the per-request override is tight (1 second).
	exec := New(Config{
		Store:               mockSt,
		MaxWorkers:          4,
		ExecutionTTL:        time.Hour,
		ResultTTL:           time.Hour,
		ProgressTTL:         time.Minute,
		DefaultMaxExecution: 60 * time.Second,
	})

	exec.RegisterProblem("slow-problem", &slowProblem{duration: 10 * time.Millisecond})
	variant, err := variants.DefaultRegistry.Create("rand1")
	require.NoError(t, err)
	exec.RegisterVariant("rand1", variant)

	ctx := context.Background()
	userID := "test-user"

	config := &api.DEConfig{
		Executions:     1,
		Generations:    1000,
		PopulationSize: 5,
		DimensionsSize: 5,
		ObjectivesSize: 2,
		FloorLimiter:   0.0,
		CeilLimiter:    1.0,
		AlgorithmConfig: &api.DEConfig_Gde3{
			Gde3: &api.GDE3Config{Cr: 0.9, F: 0.5, P: 0.1},
		},
	}

	// Submit with 1-second per-request override (MaxExecutionSeconds=1)
	executionID, err := exec.SubmitExecution(ctx, userID, "gde3", "slow-problem", "rand1", config, "", 1)
	require.NoError(t, err)

	// Should fail within a few seconds due to the 1-second per-request timeout
	assert.Eventually(t, func() bool {
		execution, err := mockSt.GetExecution(ctx, executionID, userID)
		if err != nil {
			return false
		}
		return execution.Status == store.ExecutionStatusFailed
	}, 10*time.Second, 100*time.Millisecond, "Execution with per-request timeout should fail quickly")
}

// TestExecutor_IdempotencyKey tests that idempotency key is stored with the execution.
func TestExecutor_IdempotencyKey(t *testing.T) {
	mockSt := newMockStore()
	exec := New(Config{
		Store:        mockSt,
		MaxWorkers:   2,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	prob, err := problems.DefaultRegistry.Create("zdt1", 10, 2)
	require.NoError(t, err)
	exec.RegisterProblem("zdt1", prob)

	variant, err := variants.DefaultRegistry.Create("rand1")
	require.NoError(t, err)
	exec.RegisterVariant("rand1", variant)

	ctx := context.Background()
	userID := "test-user"

	config := &api.DEConfig{
		Executions:     1,
		Generations:    2,
		PopulationSize: 10,
		DimensionsSize: 10,
		ObjectivesSize: 2,
		FloorLimiter:   0.0,
		CeilLimiter:    1.0,
		AlgorithmConfig: &api.DEConfig_Gde3{
			Gde3: &api.GDE3Config{Cr: 0.9, F: 0.5, P: 0.1},
		},
	}

	iKey := "test-idem-key-abc"
	executionID, err := exec.SubmitExecution(ctx, userID, "gde3", "zdt1", "rand1", config, iKey, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, executionID)

	// Verify the idempotency key is persisted in the execution record
	execution, err := mockSt.GetExecution(ctx, executionID, userID)
	require.NoError(t, err)
	assert.Equal(t, iKey, execution.IdempotencyKey)
}
