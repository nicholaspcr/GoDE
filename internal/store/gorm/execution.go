package gorm

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/nicholaspcr/GoDE/internal/store"
	storeerrors "github.com/nicholaspcr/GoDE/internal/store/errors"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/gorm"
)

// executionModel represents the database model for executions.
type executionModel struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)"`
	UserID      string    `gorm:"type:varchar(255);not null;index"`
	Status      string    `gorm:"type:varchar(20);not null;index"`
	ConfigJSON  string    `gorm:"type:text;not null"`
	Algorithm   string    `gorm:"type:varchar(50);not null;default:''"`
	Variant     string    `gorm:"type:varchar(50);not null;default:''"`
	Problem     string    `gorm:"type:varchar(50);not null;default:''"`
	ParetoID    *uint64   `gorm:"type:bigint;index"`
	Error       string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"not null;index"`
	UpdatedAt   time.Time `gorm:"not null"`
	CompletedAt *time.Time
}

func (executionModel) TableName() string {
	return "executions"
}

// executionStore implements ExecutionOperations using GORM.
type executionStore struct {
	db *gorm.DB
}

func newExecutionStore(db *gorm.DB) *executionStore {
	return &executionStore{db: db}
}

// CreateExecution creates a new execution record in the database.
func (s *executionStore) CreateExecution(ctx context.Context, execution *store.Execution) error {
	configJSON, err := protojson.Marshal(execution.Config)
	if err != nil {
		return err
	}

	model := &executionModel{
		ID:         execution.ID,
		UserID:     execution.UserID,
		Status:     string(execution.Status),
		ConfigJSON: string(configJSON),
		Algorithm:  execution.Algorithm,
		Variant:    execution.Variant,
		Problem:    execution.Problem,
		ParetoID:   execution.ParetoID,
		Error:      execution.Error,
		CreatedAt:  execution.CreatedAt,
		UpdatedAt:  execution.UpdatedAt,
	}

	return s.db.WithContext(ctx).Create(model).Error
}

// GetExecution retrieves an execution by ID and verifies ownership.
func (s *executionStore) GetExecution(ctx context.Context, executionID, userID string) (*store.Execution, error) {
	var model executionModel
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", executionID, userID).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, store.ErrExecutionNotFound
		}
		return nil, err
	}

	return s.modelToExecution(&model)
}

// UpdateExecutionStatus updates the status of an execution.
func (s *executionStore) UpdateExecutionStatus(ctx context.Context, executionID string, status store.ExecutionStatus, errorMsg string) error {
	updates := map[string]any{
		"status":     string(status),
		"updated_at": time.Now(),
	}

	if errorMsg != "" {
		updates["error"] = errorMsg
	}

	if status == store.ExecutionStatusCompleted || status == store.ExecutionStatusFailed || status == store.ExecutionStatusCancelled {
		updates["completed_at"] = time.Now()
	}

	return s.db.WithContext(ctx).Model(&executionModel{}).Where("id = ?", executionID).Updates(updates).Error
}

// UpdateExecutionResult updates the pareto ID for a completed execution.
func (s *executionStore) UpdateExecutionResult(ctx context.Context, executionID string, paretoID uint64) error {
	return s.db.WithContext(ctx).Model(&executionModel{}).Where("id = ?", executionID).Updates(map[string]any{
		"pareto_id":  paretoID,
		"updated_at": time.Now(),
	}).Error
}

// ListExecutions retrieves executions for a user with pagination, optionally filtered by status.
func (s *executionStore) ListExecutions(ctx context.Context, userID string, status *store.ExecutionStatus, limit, offset int) ([]*store.Execution, int, error) {
	// Apply defaults and max limits
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	query := s.db.WithContext(ctx).Where("user_id = ?", userID)

	if status != nil {
		query = query.Where("status = ?", string(*status))
	}

	// Get total count
	var totalCount int64
	if err := query.Model(&executionModel{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	query = query.Order("created_at DESC").Limit(limit).Offset(offset)

	var models []executionModel
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	executions := make([]*store.Execution, 0, len(models))
	for _, model := range models {
		execution, err := s.modelToExecution(&model)
		if err != nil {
			slog.Warn("skipping invalid execution record",
				slog.String("execution_id", model.ID),
				slog.String("error", err.Error()),
			)
			continue
		}
		executions = append(executions, execution)
	}

	return executions, int(totalCount), nil
}

// DeleteExecution removes an execution from the database.
func (s *executionStore) DeleteExecution(ctx context.Context, executionID, userID string) error {
	result := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", executionID, userID).Delete(&executionModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return store.ErrExecutionNotFound
	}
	return nil
}

// SaveProgress is not implemented for GORM store (handled by Redis).
func (s *executionStore) SaveProgress(ctx context.Context, progress *store.ExecutionProgress) error {
	return storeerrors.ErrProgressNotSupported
}

// GetProgress is not implemented for GORM store (handled by Redis).
func (s *executionStore) GetProgress(ctx context.Context, executionID string) (*store.ExecutionProgress, error) {
	return nil, storeerrors.ErrProgressNotSupported
}

// MarkExecutionForCancellation is not implemented for GORM store (handled by Redis).
func (s *executionStore) MarkExecutionForCancellation(ctx context.Context, executionID, userID string) error {
	return storeerrors.ErrCancellationNotSupported
}

// IsExecutionCancelled is not implemented for GORM store (handled by Redis).
func (s *executionStore) IsExecutionCancelled(ctx context.Context, executionID string) (bool, error) {
	return false, storeerrors.ErrCancellationNotSupported
}

// Subscribe is not implemented for GORM store (handled by Redis).
func (s *executionStore) Subscribe(ctx context.Context, channel string) (<-chan []byte, error) {
	return nil, storeerrors.ErrPubSubNotSupported
}

// modelToExecution converts a database model to a store.Execution.
func (s *executionStore) modelToExecution(model *executionModel) (*store.Execution, error) {
	var config api.DEConfig
	if err := protojson.Unmarshal([]byte(model.ConfigJSON), &config); err != nil {
		return nil, err
	}

	return &store.Execution{
		ID:          model.ID,
		UserID:      model.UserID,
		Status:      store.ExecutionStatus(model.Status),
		Config:      &config,
		Algorithm:   model.Algorithm,
		Variant:     model.Variant,
		Problem:     model.Problem,
		ParetoID:    model.ParetoID,
		Error:       model.Error,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		CompletedAt: model.CompletedAt,
	}, nil
}
