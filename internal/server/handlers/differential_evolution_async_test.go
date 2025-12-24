package handlers

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	"github.com/nicholaspcr/GoDE/internal/executor"
	"github.com/nicholaspcr/GoDE/internal/server/middleware"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	_ "github.com/nicholaspcr/GoDE/pkg/problems/many/dtlz" // Register DTLZ problems
	_ "github.com/nicholaspcr/GoDE/pkg/problems/many/wfg"  // Register WFG problems
	_ "github.com/nicholaspcr/GoDE/pkg/problems/multi"     // Register ZDT and VNT problems
	"github.com/nicholaspcr/GoDE/pkg/variants"
	_ "github.com/nicholaspcr/GoDE/pkg/variants/best"            // Register best/* variants
	_ "github.com/nicholaspcr/GoDE/pkg/variants/current-to-best" // Register current-to-best/* variants
	_ "github.com/nicholaspcr/GoDE/pkg/variants/pbest"           // Register pbest/* variants
	_ "github.com/nicholaspcr/GoDE/pkg/variants/rand"            // Register rand/* variants
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// testStore is a minimal in-memory store for testing
type testStore struct {
	executions     map[string]*store.Execution
	progress       map[string]*store.ExecutionProgress
	paretoSets     map[uint64]*store.ParetoSet
	cancelledExecs map[string]bool // Track cancellation requests
	nextID         uint64
}

func newTestStore() *testStore {
	return &testStore{
		executions:     make(map[string]*store.Execution),
		progress:       make(map[string]*store.ExecutionProgress),
		paretoSets:     make(map[uint64]*store.ParetoSet),
		cancelledExecs: make(map[string]bool),
		nextID:         1,
	}
}

func (ts *testStore) CreateExecution(ctx context.Context, execution *store.Execution) error {
	ts.executions[execution.ID] = execution
	return nil
}

func (ts *testStore) GetExecution(ctx context.Context, executionID, userID string) (*store.Execution, error) {
	exec, exists := ts.executions[executionID]
	if !exists {
		return nil, store.ErrExecutionNotFound
	}
	if exec.UserID != userID {
		return nil, store.ErrExecutionNotFound
	}
	return exec, nil
}

func (ts *testStore) UpdateExecutionStatus(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
	exec, exists := ts.executions[executionID]
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

func (ts *testStore) UpdateExecutionResult(ctx context.Context, executionID string, paretoID uint64) error {
	exec, exists := ts.executions[executionID]
	if !exists {
		return store.ErrExecutionNotFound
	}
	exec.ParetoID = &paretoID
	return nil
}

func (ts *testStore) ListExecutions(ctx context.Context, userID string, statusFilter *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
	var allMatching []*store.Execution
	for _, exec := range ts.executions {
		if exec.UserID == userID {
			if statusFilter == nil || exec.Status == *statusFilter {
				allMatching = append(allMatching, exec)
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
	end := start + limit
	if end > totalCount {
		end = totalCount
	}

	return allMatching[start:end], totalCount, nil
}

func (ts *testStore) DeleteExecution(ctx context.Context, executionID, userID string) error {
	exec, exists := ts.executions[executionID]
	if !exists || exec.UserID != userID {
		return store.ErrExecutionNotFound
	}
	delete(ts.executions, executionID)
	return nil
}

func (ts *testStore) SaveProgress(ctx context.Context, progress *store.ExecutionProgress) error {
	ts.progress[progress.ExecutionID] = progress
	return nil
}

func (ts *testStore) GetProgress(ctx context.Context, executionID string) (*store.ExecutionProgress, error) {
	prog, exists := ts.progress[executionID]
	if !exists {
		return nil, store.ErrExecutionNotFound
	}
	return prog, nil
}

func (ts *testStore) MarkExecutionForCancellation(ctx context.Context, executionID, userID string) error {
	exec, exists := ts.executions[executionID]
	if !exists || exec.UserID != userID {
		return store.ErrExecutionNotFound
	}
	ts.cancelledExecs[executionID] = true
	return nil
}

func (ts *testStore) IsExecutionCancelled(ctx context.Context, executionID string) (bool, error) {
	cancelled, exists := ts.cancelledExecs[executionID]
	return exists && cancelled, nil
}

func (ts *testStore) Subscribe(ctx context.Context, channel string) (<-chan []byte, error) {
	// For testing, return a channel that can be populated externally
	ch := make(chan []byte, 10)

	// For most tests, close immediately (no progress updates)
	// Tests that need progress updates will need to implement their own Subscribe
	go func() {
		<-ctx.Done()
		close(ch)
	}()

	return ch, nil
}

func (ts *testStore) CreateParetoSet(ctx context.Context, paretoSet *store.ParetoSet) error {
	paretoSet.ID = ts.nextID
	ts.paretoSets[ts.nextID] = paretoSet
	ts.nextID++
	return nil
}

func (ts *testStore) GetParetoSetByID(ctx context.Context, id uint64) (*store.ParetoSet, error) {
	ps, exists := ts.paretoSets[id]
	if !exists {
		return nil, store.ErrParetoSetNotFound
	}
	return ps, nil
}

// Stub implementations for required interface methods
func (ts *testStore) CreateUser(ctx context.Context, user *api.User) error                          { return nil }
func (ts *testStore) GetUser(ctx context.Context, ids *api.UserIDs) (*api.User, error)             { return nil, nil }
func (ts *testStore) UpdateUser(ctx context.Context, user *api.User, fields ...string) error       { return nil }
func (ts *testStore) DeleteUser(ctx context.Context, ids *api.UserIDs) error                       { return nil }
func (ts *testStore) CreatePareto(ctx context.Context, pareto *api.Pareto) error                   { return nil }
func (ts *testStore) GetPareto(ctx context.Context, ids *api.ParetoIDs) (*api.Pareto, error)       { return nil, nil }
func (ts *testStore) UpdatePareto(ctx context.Context, pareto *api.Pareto, fields ...string) error { return nil }
func (ts *testStore) DeletePareto(ctx context.Context, ids *api.ParetoIDs) error                   { return nil }
func (ts *testStore) ListParetos(ctx context.Context, ids *api.UserIDs) ([]*api.Pareto, error)     { return nil, nil }
func (ts *testStore) HealthCheck(ctx context.Context) error                                        { return nil }

func setupTestHandler() (*deHandler, *testStore) {
	ts := newTestStore()

	exec := executor.New(executor.Config{
		Store:        ts,
		MaxWorkers:   2,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})

	// Register test problem and variant
	prob, _ := problems.DefaultRegistry.Create("zdt1", 10, 2)
	exec.RegisterProblem("zdt1", prob)

	variant, _ := variants.DefaultRegistry.Create("rand1")
	exec.RegisterVariant("rand1", variant)

	handler := NewDEHandler(exec).(*deHandler)
	handler.SetStore(ts)

	return handler, ts
}

func TestRunAsync_Success(t *testing.T) {
	handler, _ := setupTestHandler()

	// Create context with authenticated user
	ctx := middleware.ContextWithUsername(context.Background(), "testuser")

	req := &api.RunAsyncRequest{
		Algorithm: "gde3",
		Problem:   "zdt1",
		Variant:   "rand1",
		DeConfig: &api.DEConfig{
			Executions:     1,
			Generations:    2,
			PopulationSize: 10,
			DimensionsSize: 10,
			ObjectivesSize:  2,
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
	}

	resp, err := handler.RunAsync(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.ExecutionId)
}

func TestRunAsync_UnauthenticatedUser(t *testing.T) {
	handler, _ := setupTestHandler()

	req := &api.RunAsyncRequest{
		Algorithm: "gde3",
		Problem:   "zdt1",
		Variant:   "rand1",
		DeConfig:  &api.DEConfig{},
	}

	_, err := handler.RunAsync(context.Background(), req)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestRunAsync_ValidationError(t *testing.T) {
	handler, _ := setupTestHandler()

	ctx := middleware.ContextWithUsername(context.Background(), "testuser")

	req := &api.RunAsyncRequest{
		Algorithm: "gde3",
		Problem:   "zdt1",
		Variant:   "rand1",
		DeConfig: &api.DEConfig{
			// Invalid: missing required fields
			Executions: -1,
		},
	}

	_, err := handler.RunAsync(ctx, req)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestGetExecutionStatus_Success(t *testing.T) {
	handler, ts := setupTestHandler()

	ctx := middleware.ContextWithUsername(context.Background(), "testuser")

	// Create a test execution
	executionID := "test-exec-123"
	_ = ts.CreateExecution(ctx, &store.Execution{
		ID:        executionID,
		UserID:    "testuser",
		Status:    store.ExecutionStatusRunning,
		Config:    &api.DEConfig{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	req := &api.GetExecutionStatusRequest{
		ExecutionId: executionID,
	}

	resp, err := handler.GetExecutionStatus(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Execution)
	assert.Equal(t, executionID, resp.Execution.Id)
	assert.Equal(t, api.ExecutionStatus_EXECUTION_STATUS_RUNNING, resp.Execution.Status)
}

func TestGetExecutionStatus_NotFound(t *testing.T) {
	handler, _ := setupTestHandler()

	ctx := middleware.ContextWithUsername(context.Background(), "testuser")

	req := &api.GetExecutionStatusRequest{
		ExecutionId: "non-existent",
	}

	_, err := handler.GetExecutionStatus(ctx, req)
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestListExecutions_Success(t *testing.T) {
	handler, ts := setupTestHandler()

	ctx := middleware.ContextWithUsername(context.Background(), "testuser")

	// Create test executions
	_ = ts.CreateExecution(ctx, &store.Execution{
		ID:        "exec-1",
		UserID:    "testuser",
		Status:    store.ExecutionStatusCompleted,
		Config:    &api.DEConfig{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	_ = ts.CreateExecution(ctx, &store.Execution{
		ID:        "exec-2",
		UserID:    "testuser",
		Status:    store.ExecutionStatusRunning,
		Config:    &api.DEConfig{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	// Create execution for different user (should not be included)
	_ = ts.CreateExecution(ctx, &store.Execution{
		ID:        "exec-3",
		UserID:    "otheruser",
		Status:    store.ExecutionStatusCompleted,
		Config:    &api.DEConfig{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	req := &api.ListExecutionsRequest{}

	resp, err := handler.ListExecutions(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Executions, 2)
	assert.Equal(t, int32(2), resp.TotalCount)
	assert.Equal(t, int32(50), resp.Limit)  // Default limit
	assert.Equal(t, int32(0), resp.Offset)  // Default offset
	assert.False(t, resp.HasMore)
}

func TestListExecutions_WithStatusFilter(t *testing.T) {
	handler, ts := setupTestHandler()

	ctx := middleware.ContextWithUsername(context.Background(), "testuser")

	// Create test executions with different statuses
	_ = ts.CreateExecution(ctx, &store.Execution{
		ID:        "exec-1",
		UserID:    "testuser",
		Status:    store.ExecutionStatusCompleted,
		Config:    &api.DEConfig{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	_ = ts.CreateExecution(ctx, &store.Execution{
		ID:        "exec-2",
		UserID:    "testuser",
		Status:    store.ExecutionStatusRunning,
		Config:    &api.DEConfig{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	req := &api.ListExecutionsRequest{
		Status: api.ExecutionStatus_EXECUTION_STATUS_COMPLETED,
	}

	resp, err := handler.ListExecutions(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Executions, 1)
	assert.Equal(t, "exec-1", resp.Executions[0].Id)
	assert.Equal(t, int32(1), resp.TotalCount)
	assert.Equal(t, int32(50), resp.Limit)
	assert.Equal(t, int32(0), resp.Offset)
	assert.False(t, resp.HasMore)
}

func TestCancelExecution_Success(t *testing.T) {
	handler, ts := setupTestHandler()

	ctx := middleware.ContextWithUsername(context.Background(), "testuser")

	// Create a running execution
	executionID := "test-exec-123"
	_ = ts.CreateExecution(ctx, &store.Execution{
		ID:        executionID,
		UserID:    "testuser",
		Status:    store.ExecutionStatusRunning,
		Config:    &api.DEConfig{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	req := &api.CancelExecutionRequest{
		ExecutionId: executionID,
	}

	_, err := handler.CancelExecution(ctx, req)
	require.NoError(t, err)
}

func TestDeleteExecution_Success(t *testing.T) {
	handler, ts := setupTestHandler()

	ctx := middleware.ContextWithUsername(context.Background(), "testuser")

	// Create an execution
	executionID := "test-exec-123"
	_ = ts.CreateExecution(ctx, &store.Execution{
		ID:        executionID,
		UserID:    "testuser",
		Status:    store.ExecutionStatusCompleted,
		Config:    &api.DEConfig{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	req := &api.DeleteExecutionRequest{
		ExecutionId: executionID,
	}

	_, err := handler.DeleteExecution(ctx, req)
	require.NoError(t, err)

	// Verify execution was deleted
	_, getErr := ts.GetExecution(ctx, executionID, "testuser")
	assert.Error(t, getErr)
}

func TestGetExecutionResults_Success(t *testing.T) {
	handler, ts := setupTestHandler()

	ctx := middleware.ContextWithUsername(context.Background(), "testuser")

	// Create a completed execution with Pareto results
	executionID := "test-exec-123"
	paretoID := uint64(1)

	// Create Pareto set
	ts.paretoSets[paretoID] = &store.ParetoSet{
		ID: paretoID,
		Vectors: []*api.Vector{
			{
				Elements:         []float64{0.1, 0.2, 0.3},
				Objectives:       []float64{0.5, 0.6},
				CrowdingDistance: 1.0,
			},
			{
				Elements:         []float64{0.4, 0.5, 0.6},
				Objectives:       []float64{0.7, 0.8},
				CrowdingDistance: 2.0,
			},
		},
		MaxObjectives: []*store.MaxObjectives{
			{Values: []float64{1.0, 2.0}},
		},
	}

	_ = ts.CreateExecution(ctx, &store.Execution{
		ID:        executionID,
		UserID:    "testuser",
		Status:    store.ExecutionStatusCompleted,
		Config:    &api.DEConfig{},
		ParetoID:  &paretoID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	req := &api.GetExecutionResultsRequest{
		ExecutionId: executionID,
	}

	resp, err := handler.GetExecutionResults(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Pareto)
	assert.Len(t, resp.Pareto.Vectors, 2)
	assert.Len(t, resp.Pareto.MaxObjs, 2)
	assert.Equal(t, 1.0, resp.Pareto.MaxObjs[0])
	assert.Equal(t, 2.0, resp.Pareto.MaxObjs[1])
}

func TestGetExecutionResults_NotAuthenticated(t *testing.T) {
	handler, _ := setupTestHandler()

	ctx := context.Background() // No authentication

	req := &api.GetExecutionResultsRequest{
		ExecutionId: "test-exec-123",
	}

	_, err := handler.GetExecutionResults(ctx, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not authenticated")
}

func TestGetExecutionResults_NotFound(t *testing.T) {
	handler, _ := setupTestHandler()

	ctx := middleware.ContextWithUsername(context.Background(), "testuser")

	req := &api.GetExecutionResultsRequest{
		ExecutionId: "nonexistent",
	}

	_, err := handler.GetExecutionResults(ctx, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestGetExecutionResults_NotCompleted(t *testing.T) {
	handler, ts := setupTestHandler()

	ctx := middleware.ContextWithUsername(context.Background(), "testuser")

	executionID := "test-exec-123"
	_ = ts.CreateExecution(ctx, &store.Execution{
		ID:        executionID,
		UserID:    "testuser",
		Status:    store.ExecutionStatusRunning, // Not completed
		Config:    &api.DEConfig{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	req := &api.GetExecutionResultsRequest{
		ExecutionId: executionID,
	}

	_, err := handler.GetExecutionResults(ctx, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not completed")
}

func TestGetExecutionResults_NoParetoID(t *testing.T) {
	handler, ts := setupTestHandler()

	ctx := middleware.ContextWithUsername(context.Background(), "testuser")

	executionID := "test-exec-123"
	_ = ts.CreateExecution(ctx, &store.Execution{
		ID:        executionID,
		UserID:    "testuser",
		Status:    store.ExecutionStatusCompleted,
		Config:    &api.DEConfig{},
		ParetoID:  nil, // No pareto results
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	req := &api.GetExecutionResultsRequest{
		ExecutionId: executionID,
	}

	_, err := handler.GetExecutionResults(ctx, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "results not found")
}

func TestCancellationIntegration(t *testing.T) {
	handler, ts := setupTestHandler()

	ctx := middleware.ContextWithUsername(context.Background(), "testuser")

	// Submit an execution with large parameters to ensure it runs long enough
	req := &api.RunAsyncRequest{
		Algorithm: "gde3",
		Problem:   "zdt1",
		Variant:   "rand1",
		DeConfig: &api.DEConfig{
			Executions:     5,     // Multiple executions
			Generations:    1000,  // Many generations
			PopulationSize: 100,   // Large population
			DimensionsSize: 30,    // More dimensions
			ObjectivesSize:  2,
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
	}

	resp, err := handler.RunAsync(ctx, req)
	require.NoError(t, err)
	executionID := resp.ExecutionId

	// Immediately cancel (before it gets too far)
	cancelReq := &api.CancelExecutionRequest{
		ExecutionId: executionID,
	}
	_, err = handler.CancelExecution(ctx, cancelReq)
	require.NoError(t, err)

	// Verify cancellation was marked in the store
	cancelled, err := ts.IsExecutionCancelled(ctx, executionID)
	require.NoError(t, err)
	assert.True(t, cancelled, "execution should be marked for cancellation")

	// Note: The actual status may still be any value because:
	// 1. Cancellation is async - execution may complete before cancellation takes effect
	// 2. Execution may not have started yet
	// 3. Execution may be running and will be cancelled in next iteration check
	// The important part is that the cancellation flag was set successfully
}

// mockStreamServer is a mock implementation of the gRPC stream server
type mockStreamServer struct {
	ctx      context.Context
	sent     []*api.StreamProgressResponse
	sendErr  error
	mu       sync.Mutex
}

func newMockStreamServer(ctx context.Context) *mockStreamServer {
	return &mockStreamServer{
		ctx:  ctx,
		sent: make([]*api.StreamProgressResponse, 0),
	}
}

func (m *mockStreamServer) Send(resp *api.StreamProgressResponse) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.sendErr != nil {
		return m.sendErr
	}
	m.sent = append(m.sent, resp)
	return nil
}

func (m *mockStreamServer) SetHeader(metadata.MD) error  { return nil }
func (m *mockStreamServer) SendHeader(metadata.MD) error { return nil }
func (m *mockStreamServer) SetTrailer(metadata.MD)       {}
func (m *mockStreamServer) Context() context.Context     { return m.ctx }
func (m *mockStreamServer) SendMsg(any) error    { return nil }
func (m *mockStreamServer) RecvMsg(any) error    { return nil }

func (m *mockStreamServer) getSentMessages() []*api.StreamProgressResponse {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]*api.StreamProgressResponse{}, m.sent...)
}

// testStoreWithStreaming is a testStore that supports streaming progress updates
type testStoreWithStreaming struct {
	*testStore
	progressChannel chan []byte
}

func newTestStoreWithStreaming() *testStoreWithStreaming {
	return &testStoreWithStreaming{
		testStore:       newTestStore(),
		progressChannel: make(chan []byte, 10),
	}
}

func (ts *testStoreWithStreaming) Subscribe(ctx context.Context, channel string) (<-chan []byte, error) {
	// Return the progress channel that we can control in tests
	outCh := make(chan []byte, 10)

	go func() {
		defer close(outCh)
		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-ts.progressChannel:
				if !ok {
					return
				}
				outCh <- data
			}
		}
	}()

	return outCh, nil
}

func (ts *testStoreWithStreaming) sendProgress(progress *store.ExecutionProgress) error {
	data, err := json.Marshal(progress)
	if err != nil {
		return err
	}
	ts.progressChannel <- data
	return nil
}

func (ts *testStoreWithStreaming) close() {
	close(ts.progressChannel)
}

func TestStreamProgress_Success(t *testing.T) {
	// Create store with streaming support
	ts := newTestStoreWithStreaming()
	defer ts.close()

	// Create executor and handler
	exec := executor.New(executor.Config{
		Store:        ts,
		MaxWorkers:   2,
		ExecutionTTL: time.Hour,
		ResultTTL:    time.Hour,
		ProgressTTL:  time.Minute,
	})
	handler := NewDEHandler(exec).(*deHandler)
	handler.SetStore(ts)

	// Create test context with authentication
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ctx = middleware.ContextWithUsername(ctx, "testuser")

	// Create execution
	executionID := "test-exec-123"
	_ = ts.CreateExecution(ctx, &store.Execution{
		ID:        executionID,
		UserID:    "testuser",
		Status:    store.ExecutionStatusRunning,
		Config:    &api.DEConfig{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	// Create mock stream
	mockStream := newMockStreamServer(ctx)

	// Start streaming in a goroutine
	streamErr := make(chan error, 1)
	go func() {
		req := &api.StreamProgressRequest{
			ExecutionId: executionID,
		}
		streamErr <- handler.StreamProgress(req, mockStream)
	}()

	// Send some progress updates
	time.Sleep(50 * time.Millisecond)

	progress1 := &store.ExecutionProgress{
		ExecutionID:         executionID,
		CurrentGeneration:   10,
		TotalGenerations:    100,
		CompletedExecutions: 1,
		TotalExecutions:     5,
		UpdatedAt:           time.Now(),
	}
	err := ts.sendProgress(progress1)
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	progress2 := &store.ExecutionProgress{
		ExecutionID:         executionID,
		CurrentGeneration:   50,
		TotalGenerations:    100,
		CompletedExecutions: 3,
		TotalExecutions:     5,
		UpdatedAt:           time.Now(),
	}
	err = ts.sendProgress(progress2)
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	// Cancel context to end streaming
	cancel()

	// Wait for stream to complete
	select {
	case err := <-streamErr:
		// Context cancelled is expected
		if err != nil && err != context.Canceled {
			t.Fatalf("unexpected stream error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("stream did not complete in time")
	}

	// Verify messages were sent
	sent := mockStream.getSentMessages()
	assert.GreaterOrEqual(t, len(sent), 2, "should have received at least 2 progress updates")

	if len(sent) >= 2 {
		// Check first progress
		assert.Equal(t, int32(10), sent[0].CurrentGeneration)
		assert.Equal(t, int32(100), sent[0].TotalGenerations)
		assert.Equal(t, int32(1), sent[0].CompletedExecutions)
		assert.Equal(t, int32(5), sent[0].TotalExecutions)

		// Check second progress
		assert.Equal(t, int32(50), sent[1].CurrentGeneration)
		assert.Equal(t, int32(3), sent[1].CompletedExecutions)
	}
}

func TestStreamProgress_NotAuthenticated(t *testing.T) {
	handler, _ := setupTestHandler()

	ctx := context.Background() // No authentication
	mockStream := newMockStreamServer(ctx)

	req := &api.StreamProgressRequest{
		ExecutionId: "test-exec-123",
	}

	err := handler.StreamProgress(req, mockStream)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not authenticated")
}

func TestStreamProgress_ExecutionNotFound(t *testing.T) {
	handler, _ := setupTestHandler()

	ctx := middleware.ContextWithUsername(context.Background(), "testuser")
	mockStream := newMockStreamServer(ctx)

	req := &api.StreamProgressRequest{
		ExecutionId: "nonexistent",
	}

	err := handler.StreamProgress(req, mockStream)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestStreamProgress_WrongUser(t *testing.T) {
	handler, ts := setupTestHandler()

	// Create execution for different user
	executionID := "test-exec-123"
	_ = ts.CreateExecution(context.Background(), &store.Execution{
		ID:        executionID,
		UserID:    "otheruser",
		Status:    store.ExecutionStatusRunning,
		Config:    &api.DEConfig{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	// Try to stream as different user
	ctx := middleware.ContextWithUsername(context.Background(), "testuser")
	mockStream := newMockStreamServer(ctx)

	req := &api.StreamProgressRequest{
		ExecutionId: executionID,
	}

	err := handler.StreamProgress(req, mockStream)
	require.Error(t, err)
	// GetExecution returns "not found" for wrong user (security: don't reveal existence)
	// or "does not belong to user" depending on implementation
	assert.True(t,
		status.Code(err) == codes.NotFound || status.Code(err) == codes.PermissionDenied,
		"should return NotFound or PermissionDenied for wrong user")
}
