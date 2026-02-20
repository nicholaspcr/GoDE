package handlers

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/internal/server/auth"
	"github.com/nicholaspcr/GoDE/internal/server/middleware"
	storerrors "github.com/nicholaspcr/GoDE/internal/store/errors"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// paretoHandler is responsible for the pareto service operations.
type paretoHandler struct {
	api.UnimplementedParetoServiceServer
	db paretoDB
}

// NewParetoHandler returns a handle that implements api's ParetoServiceServer.
func NewParetoHandler(st paretoDB) Handler {
	return &paretoHandler{db: st}
}

// RegisterService adds ParetoService to the RPC server.
func (ph *paretoHandler) RegisterService(srv *grpc.Server) {
	api.RegisterParetoServiceServer(srv, ph)
}

// RegisterHTTPHandler adds ParetoService to the grpc-gateway.
func (ph *paretoHandler) RegisterHTTPHandler(
	ctx context.Context,
	mux *runtime.ServeMux,
	lisAddr string,
	dialOpts []grpc.DialOption,
) error {
	return api.RegisterParetoServiceHandlerFromEndpoint(
		ctx, mux, lisAddr, dialOpts,
	)
}

// Get retrieves a pareto set by ID.
func (ph *paretoHandler) Get(
	ctx context.Context, req *api.ParetoServiceGetRequest,
) (*api.ParetoServiceGetResponse, error) {
	if err := middleware.RequireScope(ctx, auth.ScopeParetoRead); err != nil {
		return nil, err
	}

	if _, err := usernameFromContext(ctx); err != nil {
		return nil, err
	}

	if req.ParetoIds == nil || req.ParetoIds.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "pareto_ids.id is required")
	}

	pareto, err := ph.db.GetPareto(ctx, req.ParetoIds)
	if err != nil {
		if errors.Is(err, storerrors.ErrParetoSetNotFound) {
			return nil, status.Error(codes.NotFound, "pareto set not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get pareto set")
	}

	return &api.ParetoServiceGetResponse{Pareto: pareto}, nil
}

// Delete removes a pareto set.
func (ph *paretoHandler) Delete(
	ctx context.Context, req *api.ParetoServiceDeleteRequest,
) (*emptypb.Empty, error) {
	if err := middleware.RequireScope(ctx, auth.ScopeParetoWrite); err != nil {
		return nil, err
	}

	if _, err := usernameFromContext(ctx); err != nil {
		return nil, err
	}

	if req.ParetoIds == nil || req.ParetoIds.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "pareto_ids.id is required")
	}

	if err := ph.db.DeletePareto(ctx, req.ParetoIds); err != nil {
		if errors.Is(err, storerrors.ErrParetoSetNotFound) {
			return nil, status.Error(codes.NotFound, "pareto set not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete pareto set")
	}

	return &emptypb.Empty{}, nil
}

// ListByUser streams pareto sets for a given user with pagination.
func (ph *paretoHandler) ListByUser(
	req *api.ParetoServiceListByUserRequest,
	stream api.ParetoService_ListByUserServer,
) error {
	ctx := stream.Context()

	if err := middleware.RequireScope(ctx, auth.ScopeParetoRead); err != nil {
		return err
	}

	// Enforce that users can only list their own pareto sets unless admin
	callerUsername, err := usernameFromContext(ctx)
	if err != nil {
		return err
	}

	if req.UserIds == nil || req.UserIds.Username == "" {
		return status.Error(codes.InvalidArgument, "user_ids.username is required")
	}

	if req.UserIds.Username != callerUsername && !middleware.HasScope(ctx, auth.ScopeAdmin) {
		return status.Error(codes.PermissionDenied, "cannot list other users' pareto sets")
	}

	// Apply defaults if not provided
	limit := int(req.Limit)
	offset := int(req.Offset)
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	paretos, totalCount, err := ph.db.ListParetos(ctx, req.UserIds, limit, offset)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to list pareto sets")
	}

	hasMore := offset+len(paretos) < totalCount

	// Stream each pareto set to the client
	for i, pareto := range paretos {
		resp := &api.ParetoServiceListByUserResponse{
			Pareto: pareto,
		}
		// Include pagination metadata in first response
		if i == 0 {
			resp.TotalCount = int32(totalCount)
			resp.Limit = int32(limit)
			resp.Offset = int32(offset)
			resp.HasMore = hasMore
		}
		if err := stream.Send(resp); err != nil {
			return status.Error(codes.Internal, "failed to send pareto set")
		}
	}

	return nil
}
