package handlers

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/internal/server/auth"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// authHandler is responsible for the auth service operations.
type authHandler struct {
	db         store.Store
	jwtService auth.JWTService
	api.UnimplementedAuthServiceServer
}

// NewAuthHandler returns a handle that implements api's authServiceServer.
func NewAuthHandler(jwtService auth.JWTService) Handler {
	return &authHandler{jwtService: jwtService}
}

// setStore assings the implementation of the store to the auth handler.
func (ah *authHandler) SetStore(st store.Store) {
	ah.db = st
}

// RegisterService adds authService to the RPC server.
func (ah *authHandler) RegisterService(srv *grpc.Server) {
	api.RegisterAuthServiceServer(srv, ah)
}

// RegisterHTTPHandler adds AuthService to the grpc-gateway.
func (ah *authHandler) RegisterHTTPHandler(
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	req.User.Password = string(hashedPassword)
	if err := ah.db.CreateUser(ctx, req.User); err != nil {
		return nil, err
	}
	return api.Empty, nil
}

func (ah authHandler) Login(
	ctx context.Context, req *api.AuthServiceLoginRequest,
) (*api.AuthServiceLoginResponse, error) {
	usr, err := ah.db.GetUser(ctx, &api.UserIDs{Username: req.Username})
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := ah.jwtService.GenerateToken(usr.Ids.Username)
	if err != nil {
		return nil, err
	}

	return &api.AuthServiceLoginResponse{Token: token}, nil
}

func (ah authHandler) Logout(
	ctx context.Context, req *api.AuthServiceLogoutRequest,
) (*emptypb.Empty, error) {
	// JWT is stateless, so logout is handled client-side by discarding the token
	// This endpoint exists for API compatibility and future extensions
	// (e.g., token blacklisting, revocation lists)
	return api.Empty, nil
}
