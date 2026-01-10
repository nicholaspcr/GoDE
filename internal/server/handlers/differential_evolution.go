package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/internal/executor"
	"github.com/nicholaspcr/GoDE/internal/server/middleware"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/de"
	_ "github.com/nicholaspcr/GoDE/pkg/de/gde3"               // Register GDE3 algorithm
	"github.com/nicholaspcr/GoDE/pkg/problems"
	_ "github.com/nicholaspcr/GoDE/pkg/problems/many/dtlz" // Register DTLZ problems
	_ "github.com/nicholaspcr/GoDE/pkg/problems/many/wfg"  // Register WFG problems
	_ "github.com/nicholaspcr/GoDE/pkg/problems/multi"     // Register multi-objective problems
	"github.com/nicholaspcr/GoDE/pkg/validation"
	"github.com/nicholaspcr/GoDE/pkg/variants"
	_ "github.com/nicholaspcr/GoDE/pkg/variants/best"            // Register best variants
	_ "github.com/nicholaspcr/GoDE/pkg/variants/current-to-best" // Register current-to-best variant
	_ "github.com/nicholaspcr/GoDE/pkg/variants/pbest"           // Register pbest variant
	_ "github.com/nicholaspcr/GoDE/pkg/variants/rand"            // Register rand variants
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// deHandler is responsible for the de service operations.
type deHandler struct {
	api.UnimplementedDifferentialEvolutionServiceServer
	store.Store
	executor *executor.Executor
}

// NewDEHandler returns a handler that implements
// DifferentialEvolutionServiceServer.
func NewDEHandler(st store.Store, exec *executor.Executor) Handler {
	return &deHandler{Store: st, executor: exec}
}

// RegisterService adds DifferentialEvolutionService to the RPC server.
func (deh *deHandler) RegisterService(srv *grpc.Server) {
	api.RegisterDifferentialEvolutionServiceServer(srv, deh)
}

// RegisterHTTPHandler adds DifferentialEvolutionService to the grpc-gateway.
func (deh *deHandler) RegisterHTTPHandler(
	ctx context.Context,
	mux *runtime.ServeMux,
	lisAddr string,
	dialOpts []grpc.DialOption,
) error {
	return api.RegisterDifferentialEvolutionServiceHandlerFromEndpoint(
		ctx, mux, lisAddr, dialOpts,
	)
}

func (deh *deHandler) ListSupportedAlgorithms(
	ctx context.Context, _ *emptypb.Empty,
) (*api.ListSupportedAlgorithmsResponse, error) {
	return &api.ListSupportedAlgorithmsResponse{
		Algorithms: de.DefaultRegistry.List(),
	}, nil
}

func (deh *deHandler) ListSupportedVariants(
	ctx context.Context, _ *emptypb.Empty,
) (*api.ListSupportedVariantsResponse, error) {
	metas := variants.DefaultRegistry.ListMetadata()
	apiVariants := make([]*api.Variant, len(metas))
	for i, meta := range metas {
		apiVariants[i] = &api.Variant{
			Name:        meta.Name,
			Description: meta.Description,
		}
	}
	return &api.ListSupportedVariantsResponse{Variants: apiVariants}, nil
}

func (deh *deHandler) ListSupportedProblems(
	ctx context.Context, _ *emptypb.Empty,
) (*api.ListSupportedProblemsResponse, error) {
	metas := problems.DefaultRegistry.ListMetadata()
	apiProblems := make([]*api.Problem, len(metas))
	for i, meta := range metas {
		apiProblems[i] = &api.Problem{
			Name:        meta.Name,
			Description: meta.Description,
		}
	}
	return &api.ListSupportedProblemsResponse{Problems: apiProblems}, nil
}

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

	// Validate DE configuration
	if err := validation.ValidateDEConfig(req.DeConfig); err != nil {
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
	userID := middleware.UsernameFromContext(ctx)
	if userID == "" {
		return "", status.Error(codes.Unauthenticated, "user not authenticated")
	}

	execution, err := deh.Store.GetExecution(ctx, executionID, userID) //nolint:staticcheck // Explicit for clarity
	if err != nil {
		if errors.Is(err, store.ErrExecutionNotFound) {
			return "", status.Error(codes.NotFound, "execution not found")
		}
		span.RecordError(err)
		return "", status.Errorf(codes.Internal, "failed to get execution: %v", err)
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
				return status.Errorf(codes.Internal, "progress stream error: %v", err)
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
				return status.Errorf(codes.Internal, "failed to send progress: %v", err)
			}
		}
	}
}

// GetExecutionStatus returns the current status of an execution.
func (deh *deHandler) GetExecutionStatus(
	ctx context.Context, req *api.GetExecutionStatusRequest,
) (*api.GetExecutionStatusResponse, error) {
	tracer := otel.Tracer("handlers.de")
	ctx, span := tracer.Start(ctx, "deHandler.GetExecutionStatus")
	defer span.End()

	span.SetAttributes(attribute.String("execution_id", req.ExecutionId))

	// Get user ID from context
	userID := middleware.UsernameFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Get execution
	execution, err := deh.Store.GetExecution(ctx, req.ExecutionId, userID) //nolint:staticcheck // Explicit for clarity
	if err != nil {
		if errors.Is(err, store.ErrExecutionNotFound) {
			return nil, status.Error(codes.NotFound, "execution not found")
		}
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to get execution: %v", err)
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

	// Get user ID from context
	userID := middleware.UsernameFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Get execution
	execution, err := deh.Store.GetExecution(ctx, req.ExecutionId, userID) //nolint:staticcheck // Explicit for clarity
	if err != nil {
		if errors.Is(err, store.ErrExecutionNotFound) {
			return nil, status.Error(codes.NotFound, "execution not found")
		}
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to get execution: %v", err)
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
		return nil, status.Errorf(codes.Internal, "failed to get pareto set: %v", err)
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

// ListExecutions returns a list of executions for the current user.
func (deh *deHandler) ListExecutions(
	ctx context.Context, req *api.ListExecutionsRequest,
) (*api.ListExecutionsResponse, error) {
	tracer := otel.Tracer("handlers.de")
	ctx, span := tracer.Start(ctx, "deHandler.ListExecutions")
	defer span.End()

	// Get user ID from context
	userID := middleware.UsernameFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Convert API status filter to store status
	var statusFilter *store.ExecutionStatus
	if req.Status != api.ExecutionStatus_EXECUTION_STATUS_UNSPECIFIED {
		storeStatus := convertAPIStatusToStore(req.Status)
		statusFilter = &storeStatus
	}

	// Extract pagination parameters
	limit := int(req.Limit)
	offset := int(req.Offset)

	// List executions with pagination
	executions, totalCount, err := deh.Store.ListExecutions(ctx, userID, statusFilter, limit, offset)
	if err != nil {
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to list executions: %v", err)
	}

	// Apply defaults for response
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
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

	// Get user ID from context
	userID := middleware.UsernameFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Cancel execution
	if err := deh.executor.CancelExecution(ctx, req.ExecutionId, userID); err != nil {
		if errors.Is(err, store.ErrExecutionNotFound) {
			return nil, status.Error(codes.NotFound, "execution not found")
		}
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to cancel execution: %v", err)
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

	// Get user ID from context
	userID := middleware.UsernameFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	// Delete execution
	if err := deh.Store.DeleteExecution(ctx, req.ExecutionId, userID); err != nil {
		if errors.Is(err, store.ErrExecutionNotFound) {
			return nil, status.Error(codes.NotFound, "execution not found")
		}
		span.RecordError(err)
		return nil, status.Errorf(codes.Internal, "failed to delete execution: %v", err)
	}

	return &emptypb.Empty{}, nil
}

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
