package handlers

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// userHandler is responsible for the user service operations.
type userHandler struct {
	store.Store
	api.UnimplementedUserServiceServer
}

// NewUserHandler returns a handle that implements api's UserServiceServer.
func NewUserHandler() Handler { return &userHandler{} }

func (uh *userHandler) SetStore(st store.Store) {
	uh.Store = st
}

// RegisterService adds UserService to the RPC server.
func (uh *userHandler) RegisterService(srv *grpc.Server) {
	api.RegisterUserServiceServer(srv, uh)
}

// RegisterHTTPHandler adds UserService to the grpc-gateway.
func (uh *userHandler) RegisterHTTPHandler(
	ctx context.Context,
	mux *runtime.ServeMux,
	lisAddr string,
	dialOpts []grpc.DialOption,
) error {
	return api.RegisterUserServiceHandlerFromEndpoint(
		ctx, mux, lisAddr, dialOpts,
	)
}

func (uh *userHandler) Create(
	ctx context.Context, req *api.UserServiceCreateRequest,
) (*emptypb.Empty, error) {
	if err := uh.CreateUser(ctx, req.User); err != nil {
		return nil, err
	}
	return api.Empty, nil
}

func (uh *userHandler) Get(
	ctx context.Context, req *api.UserServiceGetRequest,
) (*api.UserServiceGetResponse, error) {
	usr, err := uh.GetUser(ctx, req.UserIds)
	if err != nil {
		return nil, err
	}
	return &api.UserServiceGetResponse{User: usr}, nil
}

func (uh *userHandler) Update(
	ctx context.Context, req *api.UserServiceUpdateRequest,
) (*emptypb.Empty, error) {
	err := uh.UpdateUser(ctx, req.User, req.FieldMask.GetPaths()...)
	if err != nil {
		return nil, err
	}
	return api.Empty, err
}

func (uh *userHandler) Delete(
	ctx context.Context, req *api.UserServiceDeleteRequest,
) (*emptypb.Empty, error) {
	if err := uh.DeleteUser(ctx, req.UserIds); err != nil {
		return nil, err
	}
	return api.Empty, nil
}
