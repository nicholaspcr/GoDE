package handlers

import (
	"time"

	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// convertExecutionStatus converts store.ExecutionStatus to api.ExecutionStatus.
func convertExecutionStatus(status store.ExecutionStatus) api.ExecutionStatus {
	switch status {
	case store.ExecutionStatusPending:
		return api.ExecutionStatus_EXECUTION_STATUS_PENDING
	case store.ExecutionStatusRunning:
		return api.ExecutionStatus_EXECUTION_STATUS_RUNNING
	case store.ExecutionStatusCompleted:
		return api.ExecutionStatus_EXECUTION_STATUS_COMPLETED
	case store.ExecutionStatusFailed:
		return api.ExecutionStatus_EXECUTION_STATUS_FAILED
	case store.ExecutionStatusCancelled:
		return api.ExecutionStatus_EXECUTION_STATUS_CANCELLED
	default:
		return api.ExecutionStatus_EXECUTION_STATUS_UNSPECIFIED
	}
}

// convertAPIStatusToStore converts api.ExecutionStatus to store.ExecutionStatus.
func convertAPIStatusToStore(status api.ExecutionStatus) store.ExecutionStatus {
	switch status {
	case api.ExecutionStatus_EXECUTION_STATUS_PENDING:
		return store.ExecutionStatusPending
	case api.ExecutionStatus_EXECUTION_STATUS_RUNNING:
		return store.ExecutionStatusRunning
	case api.ExecutionStatus_EXECUTION_STATUS_COMPLETED:
		return store.ExecutionStatusCompleted
	case api.ExecutionStatus_EXECUTION_STATUS_FAILED:
		return store.ExecutionStatusFailed
	case api.ExecutionStatus_EXECUTION_STATUS_CANCELLED:
		return store.ExecutionStatusCancelled
	default:
		return store.ExecutionStatusPending
	}
}

// executionToProto converts store.Execution to api.Execution.
func executionToProto(exec *store.Execution) *api.Execution {
	apiExec := &api.Execution{
		Id:        exec.ID,
		UserId:    exec.UserID,
		Status:    convertExecutionStatus(exec.Status),
		Config:    exec.Config,
		Algorithm: exec.Algorithm,
		Variant:   exec.Variant,
		Problem:   exec.Problem,
		CreatedAt: timestampProto(exec.CreatedAt),
		UpdatedAt: timestampProto(exec.UpdatedAt),
		Error:     exec.Error,
	}

	if exec.CompletedAt != nil {
		apiExec.CompletedAt = timestampProto(*exec.CompletedAt)
	}

	if exec.ParetoID != nil {
		apiExec.ParetoId = *exec.ParetoID
	}

	return apiExec
}

// timestampProto converts time.Time to timestamppb.Timestamp.
func timestampProto(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
