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
)

// GetExecutionStatus returns the current status of an execution.
func (deh *deHandler) GetExecutionStatus(
	ctx context.Context, req *api.GetExecutionStatusRequest,
) (*api.GetExecutionStatusResponse, error) {
	tracer := otel.Tracer("handlers.de")
	ctx, span := tracer.Start(ctx, "deHandler.GetExecutionStatus")
	defer span.End()

	span.SetAttributes(attribute.String("execution_id", req.ExecutionId))

	userID, err := usernameFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check authorization - requires de:read scope
	if err := middleware.RequireScope(ctx, auth.ScopeDERead); err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Get execution
	execution, err := deh.Store.GetExecution(ctx, req.ExecutionId, userID) //nolint:staticcheck // Explicit for clarity
	if err != nil {
		if errors.Is(err, store.ErrExecutionNotFound) {
			return nil, status.Error(codes.NotFound, "execution not found")
		}
		span.RecordError(err)
		return nil, status.Error(codes.Internal, "failed to get execution")
	}

	// Convert to API execution
	apiExecution := executionToProto(execution)

	// Get progress if available
	var apiProgress *api.StreamProgressResponse
	if execution.Status == store.ExecutionStatusRunning {
		progress, err := deh.Store.GetProgress(ctx, req.ExecutionId) //nolint:staticcheck // Explicit for clarity
		if err == nil {
			apiProgress = &api.StreamProgressResponse{
				ExecutionId:         progress.ExecutionID,
				CurrentGeneration:   progress.CurrentGeneration,
				TotalGenerations:    progress.TotalGenerations,
				CompletedExecutions: progress.CompletedExecutions,
				TotalExecutions:     progress.TotalExecutions,
				PartialPareto:       progress.PartialPareto,
			}
		}
	}

	return &api.GetExecutionStatusResponse{
		Execution: apiExecution,
		Progress:  apiProgress,
	}, nil
}

// GetExecutionResults returns the results of a completed execution.
func (deh *deHandler) GetExecutionResults(
	ctx context.Context, req *api.GetExecutionResultsRequest,
) (*api.GetExecutionResultsResponse, error) {
	tracer := otel.Tracer("handlers.de")
	ctx, span := tracer.Start(ctx, "deHandler.GetExecutionResults")
	defer span.End()

	span.SetAttributes(attribute.String("execution_id", req.ExecutionId))

	userID, err := usernameFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check authorization - requires de:read scope
	if err := middleware.RequireScope(ctx, auth.ScopeDERead); err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Get execution
	execution, err := deh.Store.GetExecution(ctx, req.ExecutionId, userID) //nolint:staticcheck // Explicit for clarity
	if err != nil {
		if errors.Is(err, store.ErrExecutionNotFound) {
			return nil, status.Error(codes.NotFound, "execution not found")
		}
		span.RecordError(err)
		return nil, status.Error(codes.Internal, "failed to get execution")
	}

	// Check if execution is completed
	if execution.Status != store.ExecutionStatusCompleted {
		return nil, status.Error(codes.FailedPrecondition, "execution is not completed")
	}

	// Check if pareto ID exists
	if execution.ParetoID == nil {
		return nil, status.Error(codes.NotFound, "execution results not found")
	}

	// Get pareto set
	paretoSet, err := deh.Store.GetParetoSetByID(ctx, *execution.ParetoID) //nolint:staticcheck // Explicit for clarity
	if err != nil {
		if errors.Is(err, store.ErrParetoSetNotFound) {
			return nil, status.Error(codes.NotFound, "pareto set not found")
		}
		span.RecordError(err)
		return nil, status.Error(codes.Internal, "failed to get pareto set")
	}

	// Flatten max objectives from store.MaxObjectives
	flatMaxObjs := make([]float64, 0)
	for _, maxObj := range paretoSet.MaxObjectives {
		if maxObj != nil && maxObj.Values != nil {
			flatMaxObjs = append(flatMaxObjs, maxObj.Values...)
		}
	}

	return &api.GetExecutionResultsResponse{
		Pareto: &api.Pareto{
			Vectors: paretoSet.Vectors,
			MaxObjs: flatMaxObjs,
		},
	}, nil
}
