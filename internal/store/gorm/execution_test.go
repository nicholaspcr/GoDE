package gorm

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/nicholaspcr/GoDE/internal/store"
	storeerrors "github.com/nicholaspcr/GoDE/internal/store/errors"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// setupExecutionTestDB creates an in-memory DB with all models migrated.
func setupExecutionTestDB(t *testing.T) *executionStore {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"))
	require.NoError(t, err)
	err = db.AutoMigrate(&executionModel{})
	require.NoError(t, err)
	return newExecutionStore(db)
}

func newTestExecution(id, userID string) *store.Execution {
	return &store.Execution{
		ID:        id,
		UserID:    userID,
		Status:    store.ExecutionStatusPending,
		Config:    &api.DEConfig{},
		Algorithm: "gde3",
		Variant:   "rand/1",
		Problem:   "zdt1",
		CreatedAt: time.Now().Truncate(time.Second),
		UpdatedAt: time.Now().Truncate(time.Second),
	}
}

func TestExecutionStore_Create_Get(t *testing.T) {
	s := setupExecutionTestDB(t)
	ctx := context.Background()

	exec := newTestExecution("exec-1", "user1")
	require.NoError(t, s.CreateExecution(ctx, exec))

	got, err := s.GetExecution(ctx, "exec-1", "user1")
	require.NoError(t, err)
	assert.Equal(t, exec.ID, got.ID)
	assert.Equal(t, exec.UserID, got.UserID)
	assert.Equal(t, exec.Algorithm, got.Algorithm)
	assert.Equal(t, exec.Variant, got.Variant)
	assert.Equal(t, exec.Problem, got.Problem)
	assert.Equal(t, exec.Status, got.Status)
}

func TestExecutionStore_GetExecution_NotFound(t *testing.T) {
	s := setupExecutionTestDB(t)
	ctx := context.Background()

	_, err := s.GetExecution(ctx, "nonexistent", "user1")
	require.Error(t, err)
	assert.Equal(t, store.ErrExecutionNotFound, err)
}

func TestExecutionStore_GetExecution_WrongUser(t *testing.T) {
	s := setupExecutionTestDB(t)
	ctx := context.Background()

	exec := newTestExecution("exec-2", "user1")
	require.NoError(t, s.CreateExecution(ctx, exec))

	// Querying with different user should return not found (ownership check)
	_, err := s.GetExecution(ctx, "exec-2", "other-user")
	require.Error(t, err)
	assert.Equal(t, store.ErrExecutionNotFound, err)
}

func TestExecutionStore_UpdateExecutionStatus(t *testing.T) {
	s := setupExecutionTestDB(t)
	ctx := context.Background()

	exec := newTestExecution("exec-3", "user1")
	require.NoError(t, s.CreateExecution(ctx, exec))

	require.NoError(t, s.UpdateExecutionStatus(ctx, "exec-3", store.ExecutionStatusRunning, ""))

	got, err := s.GetExecution(ctx, "exec-3", "user1")
	require.NoError(t, err)
	assert.Equal(t, store.ExecutionStatusRunning, got.Status)
	assert.Empty(t, got.Error)
}

func TestExecutionStore_UpdateExecutionStatus_WithError(t *testing.T) {
	s := setupExecutionTestDB(t)
	ctx := context.Background()

	exec := newTestExecution("exec-4", "user1")
	require.NoError(t, s.CreateExecution(ctx, exec))

	require.NoError(t, s.UpdateExecutionStatus(ctx, "exec-4", store.ExecutionStatusFailed, "something went wrong"))

	got, err := s.GetExecution(ctx, "exec-4", "user1")
	require.NoError(t, err)
	assert.Equal(t, store.ExecutionStatusFailed, got.Status)
	assert.Equal(t, "something went wrong", got.Error)
	assert.NotNil(t, got.CompletedAt, "CompletedAt should be set for terminal states")
}

func TestExecutionStore_UpdateExecutionStatus_CompletedAt(t *testing.T) {
	s := setupExecutionTestDB(t)
	ctx := context.Background()

	terminalStatuses := []store.ExecutionStatus{
		store.ExecutionStatusCompleted,
		store.ExecutionStatusFailed,
		store.ExecutionStatusCancelled,
	}

	for _, status := range terminalStatuses {
		id := "exec-terminal-" + string(status)
		exec := newTestExecution(id, "user1")
		require.NoError(t, s.CreateExecution(ctx, exec))
		require.NoError(t, s.UpdateExecutionStatus(ctx, id, status, ""))

		got, err := s.GetExecution(ctx, id, "user1")
		require.NoError(t, err)
		assert.NotNil(t, got.CompletedAt, "CompletedAt should be set for status %s", status)
	}
}

func TestExecutionStore_UpdateExecutionResult(t *testing.T) {
	s := setupExecutionTestDB(t)
	ctx := context.Background()

	exec := newTestExecution("exec-5", "user1")
	require.NoError(t, s.CreateExecution(ctx, exec))

	paretoID := uint64(42)
	require.NoError(t, s.UpdateExecutionResult(ctx, "exec-5", paretoID))

	got, err := s.GetExecution(ctx, "exec-5", "user1")
	require.NoError(t, err)
	require.NotNil(t, got.ParetoID)
	assert.Equal(t, paretoID, *got.ParetoID)
}

func TestExecutionStore_ListExecutions(t *testing.T) {
	s := setupExecutionTestDB(t)
	ctx := context.Background()

	// Create executions for user1 and user2
	for i := 0; i < 3; i++ {
		exec := newTestExecution("u1-exec-"+string(rune('a'+i)), "user1")
		exec.Status = store.ExecutionStatusCompleted
		require.NoError(t, s.CreateExecution(ctx, exec))
	}
	exec := newTestExecution("u2-exec-1", "user2")
	require.NoError(t, s.CreateExecution(ctx, exec))

	t.Run("list all for user1", func(t *testing.T) {
		execs, total, err := s.ListExecutions(ctx, "user1", nil, 50, 0)
		require.NoError(t, err)
		assert.Equal(t, 3, total)
		assert.Len(t, execs, 3)
	})

	t.Run("list with status filter", func(t *testing.T) {
		status := store.ExecutionStatusCompleted
		execs, total, err := s.ListExecutions(ctx, "user1", &status, 50, 0)
		require.NoError(t, err)
		assert.Equal(t, 3, total)
		assert.Len(t, execs, 3)
	})

	t.Run("list with non-matching status filter", func(t *testing.T) {
		status := store.ExecutionStatusFailed
		execs, total, err := s.ListExecutions(ctx, "user1", &status, 50, 0)
		require.NoError(t, err)
		assert.Equal(t, 0, total)
		assert.Empty(t, execs)
	})

	t.Run("list with pagination", func(t *testing.T) {
		execs, total, err := s.ListExecutions(ctx, "user1", nil, 2, 0)
		require.NoError(t, err)
		assert.Equal(t, 3, total)
		assert.Len(t, execs, 2)
	})

	t.Run("list with offset", func(t *testing.T) {
		execs, total, err := s.ListExecutions(ctx, "user1", nil, 50, 2)
		require.NoError(t, err)
		assert.Equal(t, 3, total)
		assert.Len(t, execs, 1)
	})

	t.Run("user isolation", func(t *testing.T) {
		execs, total, err := s.ListExecutions(ctx, "user2", nil, 50, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Len(t, execs, 1)
	})
}

func TestExecutionStore_ListExecutions_DefaultsAndMaxLimit(t *testing.T) {
	s := setupExecutionTestDB(t)
	ctx := context.Background()

	// Should not panic or fail with negative/zero limits
	execs, total, err := s.ListExecutions(ctx, "user1", nil, 0, -1)
	require.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Empty(t, execs)
}

func TestExecutionStore_DeleteExecution(t *testing.T) {
	s := setupExecutionTestDB(t)
	ctx := context.Background()

	exec := newTestExecution("exec-del", "user1")
	require.NoError(t, s.CreateExecution(ctx, exec))

	require.NoError(t, s.DeleteExecution(ctx, "exec-del", "user1"))

	_, err := s.GetExecution(ctx, "exec-del", "user1")
	assert.Equal(t, store.ErrExecutionNotFound, err)
}

func TestExecutionStore_DeleteExecution_NotFound(t *testing.T) {
	s := setupExecutionTestDB(t)
	ctx := context.Background()

	err := s.DeleteExecution(ctx, "nonexistent", "user1")
	assert.Equal(t, store.ErrExecutionNotFound, err)
}

func TestExecutionStore_DeleteExecution_WrongUser(t *testing.T) {
	s := setupExecutionTestDB(t)
	ctx := context.Background()

	exec := newTestExecution("exec-6", "user1")
	require.NoError(t, s.CreateExecution(ctx, exec))

	// Deleting with wrong user should fail
	err := s.DeleteExecution(ctx, "exec-6", "other-user")
	assert.Equal(t, store.ErrExecutionNotFound, err)

	// Execution should still exist for original user
	_, err = s.GetExecution(ctx, "exec-6", "user1")
	require.NoError(t, err)
}

func TestExecutionStore_UnsupportedOperations(t *testing.T) {
	s := setupExecutionTestDB(t)
	ctx := context.Background()

	t.Run("SaveProgress returns ErrProgressNotSupported", func(t *testing.T) {
		err := s.SaveProgress(ctx, &store.ExecutionProgress{})
		assert.ErrorIs(t, err, storeerrors.ErrProgressNotSupported)
	})

	t.Run("GetProgress returns ErrProgressNotSupported", func(t *testing.T) {
		_, err := s.GetProgress(ctx, "exec-id")
		assert.ErrorIs(t, err, storeerrors.ErrProgressNotSupported)
	})

	t.Run("MarkExecutionForCancellation returns ErrCancellationNotSupported", func(t *testing.T) {
		err := s.MarkExecutionForCancellation(ctx, "exec-id", "user1")
		assert.ErrorIs(t, err, storeerrors.ErrCancellationNotSupported)
	})

	t.Run("IsExecutionCancelled returns ErrCancellationNotSupported", func(t *testing.T) {
		_, err := s.IsExecutionCancelled(ctx, "exec-id")
		assert.ErrorIs(t, err, storeerrors.ErrCancellationNotSupported)
	})

	t.Run("Subscribe returns ErrPubSubNotSupported", func(t *testing.T) {
		_, err := s.Subscribe(ctx, "channel")
		assert.ErrorIs(t, err, storeerrors.ErrPubSubNotSupported)
	})
}
