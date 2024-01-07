package handlers

import (
	"context"

	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userServer struct {
	store.Store
	api.UnimplementedUserServiceServer
}

func (*userServer) Create(
	ctx context.Context, usr *api.User,
) (*emptypb.Empty, error) {
	err := status.Errorf(codes.Unimplemented, "method Create not implemented")
	return api.Empty, err
}

func (*userServer) Get(
	ctx context.Context, usrIDs *api.UserIDs,
) (*api.User, error) {
	// err := status.Errorf(codes.Unimplemented, "method Read not implemented")
	return &api.User{
		Ids:      &api.UserIDs{UserId: "nicholaspcr@gmail.com"},
		Password: "123456",
	}, nil
}

func (*userServer) Update(
	ctx context.Context, usr *api.User,
) (*emptypb.Empty, error) {
	err := status.Errorf(codes.Unimplemented, "method Update not implemented")
	return api.Empty, err
}

func (*userServer) Delete(
	ctx context.Context, usrIDs *api.UserIDs,
) (*emptypb.Empty, error) {
	err := status.Errorf(codes.Unimplemented, "method Delete not implemented")
	return api.Empty, err
}
