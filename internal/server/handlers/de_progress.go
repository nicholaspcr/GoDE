package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StreamProgress streams real-time progress updates for an execution.
func (deh *deHandler) StreamProgress(
	req *api.StreamProgressRequest,
	stream api.DifferentialEvolutionService_StreamProgressServer,
) error {
	tracer := otel.Tracer("handlers.de")
	ctx, span := tracer.Start(stream.Context(), "deHandler.StreamProgress")
	defer span.End()

	span.SetAttributes(attribute.String("execution_id", req.ExecutionId))

	// Verify user access
	_, err := deh.verifyExecutionAccess(ctx, req.ExecutionId, span)
	if err != nil {
		return err
	}

	// Set up progress channels
	progressCh := make(chan *store.ExecutionProgress, 10)
	errCh := make(chan error, 1)

	// Start progress subscription goroutine
	go deh.subscribeToProgress(ctx, req.ExecutionId, progressCh, errCh)

	// Stream progress to client
	return deh.streamProgressToClient(ctx, stream, progressCh, errCh, span)
}

// verifyExecutionAccess checks if the user has access to the execution.
func (deh *deHandler) verifyExecutionAccess(ctx context.Context, executionID string, span trace.Span) (string, error) {
	userID, err := usernameFromContext(ctx)
	if err != nil {
		return "", err
	}

	execution, err := deh.Store.GetExecution(ctx, executionID, userID) //nolint:staticcheck // Explicit for clarity
	if err != nil {
		if errors.Is(err, store.ErrExecutionNotFound) {
			return "", status.Error(codes.NotFound, "execution not found")
		}
		span.RecordError(err)
		return "", status.Error(codes.Internal, "failed to get execution")
	}

	if execution.UserID != userID {
		return "", status.Error(codes.PermissionDenied, "execution does not belong to user")
	}

	return userID, nil
}

// subscribeToProgress subscribes to Redis pub/sub for progress updates.
func (deh *deHandler) subscribeToProgress(
	ctx context.Context,
	executionID string,
	progressCh chan<- *store.ExecutionProgress,
	errCh chan<- error,
) {
	defer close(progressCh)
	defer close(errCh)

	channel := fmt.Sprintf("execution:%s:updates", executionID)
	progressBytes, err := deh.Store.Subscribe(ctx, channel) //nolint:staticcheck // Explicit for clarity
	if err != nil {
		errCh <- err
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-progressBytes:
			if !ok {
				return
			}
			progress := &store.ExecutionProgress{}
			if err := progress.UnmarshalJSON(data); err != nil {
				continue
			}
			progressCh <- progress
		}
	}
}

// streamProgressToClient sends progress updates to the gRPC stream.
func (deh *deHandler) streamProgressToClient(
	ctx context.Context,
	stream api.DifferentialEvolutionService_StreamProgressServer,
	progressCh <-chan *store.ExecutionProgress,
	errCh <-chan error,
	span trace.Span,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			if err != nil {
				span.RecordError(err)
				return status.Error(codes.Internal, "progress stream error")
			}
		case progress, ok := <-progressCh:
			if !ok {
				return nil
			}

			apiProgress := &api.StreamProgressResponse{
				ExecutionId:         progress.ExecutionID,
				CurrentGeneration:   progress.CurrentGeneration,
				TotalGenerations:    progress.TotalGenerations,
				CompletedExecutions: progress.CompletedExecutions,
				TotalExecutions:     progress.TotalExecutions,
				PartialPareto:       progress.PartialPareto,
			}

			if err := stream.Send(apiProgress); err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				span.RecordError(err)
				return status.Error(codes.Internal, "failed to send progress")
			}
		}
	}
}
