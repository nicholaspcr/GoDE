package store

import (
	"time"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
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
