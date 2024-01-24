package handlers

import (
	"context"

	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserHandler is responsible for the user service operations.
type UserHandler struct {
	store.Store
	api.UnimplementedUserServiceServer
}

func (uh *UserHandler) Create(
	ctx context.Context, req *api.UserServiceCreateRequest,
) (*emptypb.Empty, error) {
	if err := uh.Store.Create(ctx, req.User); err != nil {
		return nil, err
	}
	return api.Empty, nil
}

func (uh *UserHandler) Get(
	ctx context.Context, req *api.UserServiceGetRequest,
) (*api.UserServiceGetResponse, error) {
	usr, err := uh.Store.Get(ctx, req.UserIds)
	if err != nil {
		return nil, err
	}
	return &api.UserServiceGetResponse{User: usr}, nil
}

func (uh *UserHandler) Update(
	ctx context.Context, req *api.UserServiceUpdateRequest,
) (*emptypb.Empty, error) {
	err := uh.Store.Update(ctx, req.User, "email", "password") // TODO: Add fieldmask to request.
	if err != nil {
		return nil, err
	}
	return api.Empty, err
}

func (uh *UserHandler) Delete(
	ctx context.Context, req *api.UserServiceDeleteRequest,
) (*emptypb.Empty, error) {
	if err := uh.Store.Delete(ctx, req.UserIds); err != nil {
		return nil, err
	}
	return api.Empty, nil
}
