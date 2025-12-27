package store

import (
	"encoding/json"
	"time"

	"github.com/nicholaspcr/GoDE/internal/store/errors"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
)

// Re-export errors from the errors package for backward compatibility.
// All error definitions are centralized in internal/store/errors/errors.go.
var (
	ErrExecutionNotFound = errors.ErrExecutionNotFound
	ErrParetoSetNotFound = errors.ErrParetoSetNotFound
)

// ExecutionStatus represents the state of an execution.
type ExecutionStatus string

const (
	ExecutionStatusPending   ExecutionStatus = "pending"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusFailed    ExecutionStatus = "failed"
	ExecutionStatusCancelled ExecutionStatus = "cancelled"
)

// Execution represents a DE algorithm execution.
type Execution struct {
	ID          string
	UserID      string
	Status      ExecutionStatus
	Config      *api.DEConfig
	ParetoID    *uint64
	Error       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
}

// ExecutionProgress represents the current progress of a running execution.
type ExecutionProgress struct {
	ExecutionID         string
	CurrentGeneration   int32
	TotalGenerations    int32
	CompletedExecutions int32
	TotalExecutions     int32
	PartialPareto       []*api.Vector
	UpdatedAt           time.Time
}

// MarshalJSON implements json.Marshaler for ExecutionProgress.
func (ep *ExecutionProgress) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ExecutionID         string        `json:"execution_id"`
		CurrentGeneration   int32         `json:"current_generation"`
		TotalGenerations    int32         `json:"total_generations"`
		CompletedExecutions int32         `json:"completed_executions"`
		TotalExecutions     int32         `json:"total_executions"`
		PartialPareto       []*api.Vector `json:"partial_pareto"`
		UpdatedAt           time.Time     `json:"updated_at"`
	}{
		ExecutionID:         ep.ExecutionID,
		CurrentGeneration:   ep.CurrentGeneration,
		TotalGenerations:    ep.TotalGenerations,
		CompletedExecutions: ep.CompletedExecutions,
		TotalExecutions:     ep.TotalExecutions,
		PartialPareto:       ep.PartialPareto,
		UpdatedAt:           ep.UpdatedAt,
	})
}

// UnmarshalJSON implements json.Unmarshaler for ExecutionProgress.
func (ep *ExecutionProgress) UnmarshalJSON(data []byte) error {
	aux := struct {
		ExecutionID         string        `json:"execution_id"`
		CurrentGeneration   int32         `json:"current_generation"`
		TotalGenerations    int32         `json:"total_generations"`
		CompletedExecutions int32         `json:"completed_executions"`
		TotalExecutions     int32         `json:"total_executions"`
		PartialPareto       []*api.Vector `json:"partial_pareto"`
		UpdatedAt           time.Time     `json:"updated_at"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	ep.ExecutionID = aux.ExecutionID
	ep.CurrentGeneration = aux.CurrentGeneration
	ep.TotalGenerations = aux.TotalGenerations
	ep.CompletedExecutions = aux.CompletedExecutions
	ep.TotalExecutions = aux.TotalExecutions
	ep.PartialPareto = aux.PartialPareto
	ep.UpdatedAt = aux.UpdatedAt

	return nil
}

// MaxObjectives holds the maximum objective values for a pareto set.
type MaxObjectives struct {
	Values []float64
}

// ParetoSet represents a saved pareto set result.
type ParetoSet struct {
	ID            uint64
	UserID        string
	Vectors       []*api.Vector
	MaxObjectives []*MaxObjectives
	CreatedAt     time.Time
}
