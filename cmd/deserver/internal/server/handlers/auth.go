package handlers

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/cmd/deserver/internal/server/auth"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// authHandler is responsible for the auth service operations.
type authHandler struct {
	store.Store
	auth.SessionStore
	api.UnimplementedAuthServiceServer
}

// NewAuthHandler returns a handle that implements api's authServiceServer.
func NewAuthHandler(sessionStore auth.SessionStore) Handler {
	return &authHandler{SessionStore: sessionStore}
}

// setStore assings the implementation of the store to the auth handler.
func (uh *authHandler) SetStore(st store.Store) {
	uh.Store = st
}

// SetSessionStore assings the implementation of the session store to the auth
// handler.
func (uh *authHandler) SetSessionStore(st auth.SessionStore) {
	uh.SessionStore = st
}

// RegisterService adds authService to the RPC server.
func (uh *authHandler) RegisterService(srv *grpc.Server) {
	api.RegisterAuthServiceServer(srv, uh)
}

// RegisterHTTPHandler adds AuthService to the grpc-gateway.
func (uh *authHandler) RegisterHTTPHandler(
	ctx context.Context,
	mux *runtime.ServeMux,
	lisAddr string,
	dialOpts []grpc.DialOption,
) error {
	return api.RegisterAuthServiceHandlerFromEndpoint(
		ctx, mux, lisAddr, dialOpts,
	)
}

// Register creates an user into the database.
func (ah authHandler) Register(
	ctx context.Context, req *api.AuthServiceRegisterRequest,
) (*emptypb.Empty, error) {
	if err := ah.Store.CreateUser(ctx, req.User); err != nil {
		return nil, err
	}
	return api.Empty, nil
}

func (ah authHandler) Login(
	ctx context.Context, req *api.AuthServiceLoginRequest,
) (*api.AuthServiceLoginResponse, error) {
	usr, err := ah.Store.GetUser(ctx, &api.UserIDs{Email: req.Email})
	if err != nil {
		return nil, err
	}

	if usr.Password != req.Password {
		return nil, errors.New("invalid credentials")
	}

	authToken := base64.StdEncoding.EncodeToString([]byte(usr.Ids.Email))
	ah.SessionStore.Add(authToken)

	return &api.AuthServiceLoginResponse{Token: authToken}, nil
}

func (ah authHandler) Logout(
	ctx context.Context, req *api.AuthServiceLogoutRequest,
) (*emptypb.Empty, error) {
	authToken := base64.StdEncoding.EncodeToString([]byte(req.Email))
	ah.SessionStore.Remove(authToken)

	return api.Empty, nil
}
