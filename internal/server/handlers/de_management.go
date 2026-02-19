package handlers

import (
	"context"
	"errors"

	"github.com/nicholaspcr/GoDE/internal/server/auth"
	"github.com/nicholaspcr/GoDE/internal/server/middleware"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ListExecutions returns a list of executions for the current user.
func (deh *deHandler) ListExecutions(
	ctx context.Context, req *api.ListExecutionsRequest,
) (*api.ListExecutionsResponse, error) {
	tracer := otel.Tracer("handlers.de")
	ctx, span := tracer.Start(ctx, "deHandler.ListExecutions")
	defer span.End()

	userID, err := usernameFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Convert API status filter to store status
	var statusFilter *store.ExecutionStatus
	if req.Status != api.ExecutionStatus_EXECUTION_STATUS_UNSPECIFIED {
		storeStatus := convertAPIStatusToStore(req.Status)
		statusFilter = &storeStatus
	}

	// Extract and normalize pagination parameters before querying
	limit := int(req.Limit)
	offset := int(req.Offset)
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	// List executions with pagination
	executions, totalCount, err := deh.Store.ListExecutions(ctx, userID, statusFilter, limit, offset)
	if err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.Internal, "failed to list executions")
	}

	// Convert to API format
	apiExecutions := make([]*api.Execution, len(executions))
	for i, exec := range executions {
		apiExecutions[i] = executionToProto(exec)
	}

	// Calculate if there are more results
	hasMore := (offset + limit) < totalCount

	return &api.ListExecutionsResponse{
		Executions: apiExecutions,
		TotalCount: int32(totalCount),
		Limit:      int32(limit),
		Offset:     int32(offset),
		HasMore:    hasMore,
	}, nil
}

// CancelExecution cancels a running execution.
func (deh *deHandler) CancelExecution(
	ctx context.Context, req *api.CancelExecutionRequest,
) (*emptypb.Empty, error) {
	tracer := otel.Tracer("handlers.de")
	ctx, span := tracer.Start(ctx, "deHandler.CancelExecution")
	defer span.End()

	span.SetAttributes(attribute.String("execution_id", req.ExecutionId))

	userID, err := usernameFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check authorization - requires de:run scope to cancel executions
	if err := middleware.RequireScope(ctx, auth.ScopeDERun); err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Cancel execution
	if err := deh.executor.CancelExecution(ctx, req.ExecutionId, userID); err != nil {
		if errors.Is(err, store.ErrExecutionNotFound) {
			return nil, status.Error(codes.NotFound, "execution not found")
		}
		span.RecordError(err)
		return nil, status.Error(codes.Internal, "failed to cancel execution")
	}

	return &emptypb.Empty{}, nil
}

// DeleteExecution deletes an execution and its results.
func (deh *deHandler) DeleteExecution(
	ctx context.Context, req *api.DeleteExecutionRequest,
) (*emptypb.Empty, error) {
	tracer := otel.Tracer("handlers.de")
	ctx, span := tracer.Start(ctx, "deHandler.DeleteExecution")
	defer span.End()

	span.SetAttributes(attribute.String("execution_id", req.ExecutionId))

	userID, err := usernameFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check authorization - requires de:run scope to delete executions
	if err := middleware.RequireScope(ctx, auth.ScopeDERun); err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Delete execution
	if err := deh.Store.DeleteExecution(ctx, req.ExecutionId, userID); err != nil {
		if errors.Is(err, store.ErrExecutionNotFound) {
			return nil, status.Error(codes.NotFound, "execution not found")
		}
		span.RecordError(err)
		return nil, status.Error(codes.Internal, "failed to delete execution")
	}

	return &emptypb.Empty{}, nil
}
