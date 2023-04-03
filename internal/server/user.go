package server

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TODO: Implement the user service.

type userServices struct {
	api.UserBaseServicesServer
}

func newUser() api.UserBaseServicesServer {
	return &userServices{}
}

// Create creates a new user
func (u *userServices) Create(
	context.Context, *api.User,
) (*emptypb.Empty, error) {
	return nil, nil
}

// Read reads a user
func (u *userServices) Read(
	context.Context, *api.UserID,
) (*api.User, error) {
	return &api.User{}, nil
}

// Update updates a user
func (u *userServices) Update(
	context.Context, *api.User,
) (*emptypb.Empty, error) {
	return nil, nil
}

// Delete deletes a user
func (u *userServices) Delete(
	context.Context, *api.UserID,
) (*emptypb.Empty, error) {
	return nil, nil
}
