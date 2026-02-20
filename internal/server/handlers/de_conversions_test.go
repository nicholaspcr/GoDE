package handlers

import (
	"testing"
	"time"

	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertExecutionStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    store.ExecutionStatus
		expected api.ExecutionStatus
	}{
		{"pending", store.ExecutionStatusPending, api.ExecutionStatus_EXECUTION_STATUS_PENDING},
		{"running", store.ExecutionStatusRunning, api.ExecutionStatus_EXECUTION_STATUS_RUNNING},
		{"completed", store.ExecutionStatusCompleted, api.ExecutionStatus_EXECUTION_STATUS_COMPLETED},
		{"failed", store.ExecutionStatusFailed, api.ExecutionStatus_EXECUTION_STATUS_FAILED},
		{"cancelled", store.ExecutionStatusCancelled, api.ExecutionStatus_EXECUTION_STATUS_CANCELLED},
		{"unknown", store.ExecutionStatus("unknown"), api.ExecutionStatus_EXECUTION_STATUS_UNSPECIFIED},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertExecutionStatus(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertAPIStatusToStore(t *testing.T) {
	tests := []struct {
		name     string
		input    api.ExecutionStatus
		expected store.ExecutionStatus
	}{
		{"pending", api.ExecutionStatus_EXECUTION_STATUS_PENDING, store.ExecutionStatusPending},
		{"running", api.ExecutionStatus_EXECUTION_STATUS_RUNNING, store.ExecutionStatusRunning},
		{"completed", api.ExecutionStatus_EXECUTION_STATUS_COMPLETED, store.ExecutionStatusCompleted},
		{"failed", api.ExecutionStatus_EXECUTION_STATUS_FAILED, store.ExecutionStatusFailed},
		{"cancelled", api.ExecutionStatus_EXECUTION_STATUS_CANCELLED, store.ExecutionStatusCancelled},
		{"unspecified defaults to pending", api.ExecutionStatus_EXECUTION_STATUS_UNSPECIFIED, store.ExecutionStatusPending},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertAPIStatusToStore(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestVectorToPB(t *testing.T) {
	t.Run("copies fields and does not share slices", func(t *testing.T) {
		v := models.Vector{
			Elements:         []float64{1.0, 2.0, 3.0},
			Objectives:       []float64{0.5, 0.6},
			CrowdingDistance: 1.23,
		}
		pb := vectorToPB(v)

		assert.Equal(t, v.Elements, pb.Elements)
		assert.Equal(t, v.Objectives, pb.Objectives)
		assert.Equal(t, v.CrowdingDistance, pb.CrowdingDistance)

		pb.Elements[0] = 99.9
		assert.Equal(t, 1.0, v.Elements[0], "vectorToPB should deep-copy elements")
	})

	t.Run("empty vector", func(t *testing.T) {
		pb := vectorToPB(models.Vector{})
		assert.Empty(t, pb.Elements)
		assert.Empty(t, pb.Objectives)
		assert.Equal(t, 0.0, pb.CrowdingDistance)
	})
}

func TestVectorFromPB(t *testing.T) {
	t.Run("copies fields and does not share slices", func(t *testing.T) {
		pb := &api.Vector{
			Elements:         []float64{1.0, 2.0},
			Objectives:       []float64{0.5},
			CrowdingDistance: 2.5,
		}
		v := vectorFromPB(pb)

		assert.Equal(t, pb.Elements, v.Elements)
		assert.Equal(t, pb.Objectives, v.Objectives)
		assert.Equal(t, pb.CrowdingDistance, v.CrowdingDistance)

		pb.Elements[0] = 99.9
		assert.Equal(t, 1.0, v.Elements[0], "vectorFromPB should deep-copy elements")
	})
}

func TestExecutionToProto(t *testing.T) {
	now := time.Now()
	completedAt := now.Add(time.Minute)
	paretoID := uint64(42)

	t.Run("full execution with all fields", func(t *testing.T) {
		exec := &store.Execution{
			ID:          "exec-123",
			UserID:      "user-1",
			Status:      store.ExecutionStatusCompleted,
			Config:      &api.DEConfig{Executions: 5},
			Algorithm:   "gde3",
			Variant:     "rand1",
			Problem:     "zdt1",
			CreatedAt:   now,
			UpdatedAt:   now,
			CompletedAt: &completedAt,
			ParetoID:    &paretoID,
			Error:       "",
		}

		proto := executionToProto(exec)
		require.NotNil(t, proto)
		assert.Equal(t, "exec-123", proto.Id)
		assert.Equal(t, "user-1", proto.UserId)
		assert.Equal(t, api.ExecutionStatus_EXECUTION_STATUS_COMPLETED, proto.Status)
		assert.Equal(t, "gde3", proto.Algorithm)
		assert.Equal(t, "rand1", proto.Variant)
		assert.Equal(t, "zdt1", proto.Problem)
		assert.NotNil(t, proto.CreatedAt)
		assert.NotNil(t, proto.UpdatedAt)
		assert.NotNil(t, proto.CompletedAt)
		assert.Equal(t, uint64(42), proto.ParetoId)
	})

	t.Run("execution without optional fields", func(t *testing.T) {
		exec := &store.Execution{
			ID:        "exec-456",
			UserID:    "user-2",
			Status:    store.ExecutionStatusRunning,
			Config:    &api.DEConfig{},
			CreatedAt: now,
			UpdatedAt: now,
		}

		proto := executionToProto(exec)
		require.NotNil(t, proto)
		assert.Equal(t, "exec-456", proto.Id)
		assert.Nil(t, proto.CompletedAt)
		assert.Equal(t, uint64(0), proto.ParetoId)
	})
}
