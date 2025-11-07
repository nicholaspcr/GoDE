package handlers

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// paretoHandler is responsible for the pareto service operations.
type paretoHandler struct {
	api.UnimplementedParetoServiceServer
	store.Store
}

// NewParetoHandler returns a handle that implements api's ParetoServiceServer.
func NewParetoHandler() Handler { return &paretoHandler{} }

func (ph *paretoHandler) SetStore(st store.Store) {
	ph.Store = st
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
	if req.ParetoIds == nil || req.ParetoIds.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "pareto_ids.id is required")
	}

	pareto, err := ph.GetPareto(ctx, req.ParetoIds)
	if err != nil {
		return nil, status.Error(codes.NotFound, "pareto set not found")
	}

	return &api.ParetoServiceGetResponse{Pareto: pareto}, nil
}

// Delete removes a pareto set.
func (ph *paretoHandler) Delete(
	ctx context.Context, req *api.ParetoServiceDeleteRequest,
) (*emptypb.Empty, error) {
	if req.ParetoIds == nil || req.ParetoIds.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "pareto_ids.id is required")
	}

	if err := ph.DeletePareto(ctx, req.ParetoIds); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete pareto set")
	}

	return &emptypb.Empty{}, nil
}

// ListByUser streams all pareto sets for a given user.
func (ph *paretoHandler) ListByUser(
	req *api.ParetoServiceListByUserRequest,
	stream api.ParetoService_ListByUserServer,
) error {
	if req.UserIds == nil || req.UserIds.Username == "" {
		return status.Error(codes.InvalidArgument, "user_ids.username is required")
	}

	paretos, err := ph.ListParetos(stream.Context(), req.UserIds)
	if err != nil {
		return status.Error(codes.Internal, "failed to list pareto sets")
	}

	// Stream each pareto set to the client
	for _, pareto := range paretos {
		if err := stream.Send(&api.ParetoServiceListByUserResponse{
			Pareto: pareto,
		}); err != nil {
			return status.Error(codes.Internal, "failed to send pareto set")
		}
	}

	return nil
}
