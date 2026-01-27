package handlers

import (
	"context"
	"fmt"

	"github.com/nicholaspcr/GoDE/internal/server/auth"
	"github.com/nicholaspcr/GoDE/internal/server/middleware"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/validation"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RunAsync submits a DE execution to run in the background.
func (deh *deHandler) RunAsync(
	ctx context.Context, req *api.RunAsyncRequest,
) (*api.RunAsyncResponse, error) {
	tracer := otel.Tracer("handlers.de")
	ctx, span := tracer.Start(ctx, "deHandler.RunAsync")
	defer span.End()

	span.SetAttributes(
		attribute.String("algorithm", req.Algorithm),
		attribute.String("problem", req.Problem),
		attribute.String("variant", req.Variant),
	)

	// Get user ID from context
	userID := middleware.UsernameFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Check authorization - requires de:run scope
	if err := middleware.RequireScope(ctx, auth.ScopeDERun); err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Validate DE configuration and variant-specific constraints
	if err := validation.ValidateRunAsyncRequest(req.Algorithm, req.Variant, req.Problem, req.DeConfig); err != nil {
		span.RecordError(err)
		return nil, ValidationErrorToStatus(err)
	}

	span.SetAttributes(
		attribute.Int64("executions", req.DeConfig.Executions),
		attribute.Int64("generations", req.DeConfig.Generations),
		attribute.Int64("population_size", req.DeConfig.PopulationSize),
	)

	// Validate algorithm is supported
	if !de.DefaultRegistry.IsSupported(req.Algorithm) {
		err := fmt.Errorf("unsupported algorithm: %s", req.Algorithm)
		span.RecordError(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Submit execution with algorithm, problem, and variant names
	executionID, err := deh.executor.SubmitExecution(ctx, userID, req.Algorithm, req.Problem, req.Variant, req.DeConfig)
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to submit execution: %v", err)
	}

	span.SetAttributes(attribute.String("execution_id", executionID))

	return &api.RunAsyncResponse{
		ExecutionId: executionID,
	}, nil
}
