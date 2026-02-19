package handlers

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/internal/server/auth"
	"github.com/nicholaspcr/GoDE/internal/server/middleware"
	"github.com/nicholaspcr/GoDE/internal/store"
	storerrors "github.com/nicholaspcr/GoDE/internal/store/errors"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/validation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// userHandler is responsible for the user service operations.
type userHandler struct {
	api.UnimplementedUserServiceServer
	store.Store
}

// NewUserHandler returns a handle that implements api's UserServiceServer.
func NewUserHandler(st store.Store) Handler {
	return &userHandler{Store: st}
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
	if err := middleware.RequireScope(ctx, auth.ScopeUserWrite); err != nil {
		return nil, err
	}

	// Validate user
	if err := validation.ValidateUser(req.User); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := uh.CreateUser(ctx, req.User); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}
	return api.Empty, nil
}

func (uh *userHandler) Get(
	ctx context.Context, req *api.UserServiceGetRequest,
) (*api.UserServiceGetResponse, error) {
	if err := middleware.RequireScope(ctx, auth.ScopeUserRead); err != nil {
		return nil, err
	}

	// Users can only read their own data unless admin
	callerUsername, err := usernameFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if req.UserIds != nil && req.UserIds.Username != callerUsername && !middleware.HasScope(ctx, auth.ScopeAdmin) {
		return nil, status.Error(codes.PermissionDenied, "cannot access other users' data")
	}

	usr, err := uh.GetUser(ctx, req.UserIds)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user")
	}
	// Convert User to UserResponse (exclude password)
	userResp := &api.UserResponse{
		Ids:   usr.Ids,
		Email: usr.Email,
	}
	return &api.UserServiceGetResponse{User: userResp}, nil
}

func (uh *userHandler) Update(
	ctx context.Context, req *api.UserServiceUpdateRequest,
) (*emptypb.Empty, error) {
	if err := middleware.RequireScope(ctx, auth.ScopeUserWrite); err != nil {
		return nil, err
	}

	// Users can only update their own data unless admin
	callerUsername, err := usernameFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if req.User != nil && req.User.Ids != nil && req.User.Ids.Username != callerUsername && !middleware.HasScope(ctx, auth.ScopeAdmin) {
		return nil, status.Error(codes.PermissionDenied, "cannot modify other users' data")
	}

	if err = uh.UpdateUser(ctx, req.User, req.FieldMask.GetPaths()...); err != nil {
		if errors.Is(err, storerrors.ErrUnsupportedFieldMask) {
			return nil, status.Error(codes.InvalidArgument, "unsupported field mask")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user")
	}
	return api.Empty, nil
}

func (uh *userHandler) Delete(
	ctx context.Context, req *api.UserServiceDeleteRequest,
) (*emptypb.Empty, error) {
	if err := middleware.RequireScope(ctx, auth.ScopeUserWrite); err != nil {
		return nil, err
	}

	// Users can only delete their own account unless admin
	callerUsername, err := usernameFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if req.UserIds != nil && req.UserIds.Username != callerUsername && !middleware.HasScope(ctx, auth.ScopeAdmin) {
		return nil, status.Error(codes.PermissionDenied, "cannot delete other users' accounts")
	}

	if err := uh.DeleteUser(ctx, req.UserIds); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user")
	}
	return api.Empty, nil
}
