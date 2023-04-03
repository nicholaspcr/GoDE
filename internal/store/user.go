package store

import (
	"context"

	"github.com/nicholaspcr/GoDE/pkg/api"
)

// User is the interface for the user store.
type User interface {
	Create(context.Context, *api.User) error
	Read(context.Context, *api.UserID) error
	Update(context.Context, *api.User) error
	Delete(context.Context, *api.UserID) error
}

type user struct{}

func (u *user) Create(ctx context.Context, user *api.User) error {
	return nil
}

func (u *user) Read(ctx context.Context, userID *api.UserID) error {
	return nil
}

func (u *user) Update(ctx context.Context, user *api.User) error {
	return nil
}

func (u *user) Delete(ctx context.Context, userID *api.UserID) error {
	return nil
}
