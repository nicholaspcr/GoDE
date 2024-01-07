package handlers

import (
	"context"

	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserHandler is responsible for the user service operations.
type UserHandler struct {
	store.Store
	api.UnimplementedUserServiceServer
}

func (*UserHandler) Create(
	ctx context.Context, req *api.UserServiceCreateRequest,
) (*emptypb.Empty, error) {
	err := status.Errorf(codes.Unimplemented, "method Create not implemented")
	return api.Empty, err
}

func (*UserHandler) Get(
	ctx context.Context, req *api.UserServiceGetRequest,
) (*api.UserServiceGetResponse, error) {
	// err := status.Errorf(codes.Unimplemented, "method Read not implemented")
	return &api.UserServiceGetResponse{User: &api.User{
		Ids:      &api.UserIDs{UserId: "nicholaspcr@gmail.com"},
		Password: "123456",
	}}, nil
}

func (*UserHandler) Update(
	ctx context.Context, usr *api.UserServiceUpdateRequest,
) (*emptypb.Empty, error) {
	err := status.Errorf(codes.Unimplemented, "method Update not implemented")
	return api.Empty, err
}

func (*UserHandler) Delete(
	ctx context.Context, usrIDs *api.UserServiceDeleteRequest,
) (*emptypb.Empty, error) {
	err := status.Errorf(codes.Unimplemented, "method Delete not implemented")
	return api.Empty, err
}
