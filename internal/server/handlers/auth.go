package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/internal/server/auth"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/validation"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// authHandler is responsible for the auth service operations.
type authHandler struct {
	api.UnimplementedAuthServiceServer
	db         store.Store
	jwtService auth.JWTService
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
	// Validate user input
	if err := validation.ValidateUser(req.User); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Normalize email
	req.User.Email = strings.ToLower(strings.TrimSpace(req.User.Email))

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}
	req.User.Password = string(hashedPassword)

	// Create user
	if err := ah.db.CreateUser(ctx, req.User); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return api.Empty, nil
}

func (ah authHandler) Login(
	ctx context.Context, req *api.AuthServiceLoginRequest,
) (*api.AuthServiceLoginResponse, error) {
	// Validate username
	if err := validation.ValidateUsername(req.Username); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Validate password is not empty (don't validate full requirements for login)
	if err := validation.ValidateNonEmpty(req.Password, "password"); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usr, err := ah.db.GetUser(ctx, &api.UserIDs{Username: req.Username})
	if err != nil {
		return nil, status.Error(codes.NotFound, "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// Generate JWT token pair
	accessToken, refreshToken, err := ah.jwtService.GenerateTokenPair(usr.Ids.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate tokens")
	}

	return &api.AuthServiceLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64((24 * time.Hour).Seconds()), // Access token expiry in seconds
	}, nil
}

func (ah authHandler) Logout(
	ctx context.Context, req *api.AuthServiceLogoutRequest,
) (*emptypb.Empty, error) {
	// JWT is stateless, so logout is handled client-side by discarding the token
	// This endpoint exists for API compatibility and future extensions
	// (e.g., token blacklisting, revocation lists)
	return api.Empty, nil
}

func (ah authHandler) RefreshToken(
	ctx context.Context, req *api.AuthServiceRefreshTokenRequest,
) (*api.AuthServiceRefreshTokenResponse, error) {
	// Validate refresh token is provided
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	// Use the JWT service to refresh the access token
	accessToken, newRefreshToken, err := ah.jwtService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		// Map JWT errors to appropriate gRPC status codes
		switch err {
		case auth.ErrExpiredToken:
			return nil, status.Error(codes.Unauthenticated, "refresh token has expired")
		case auth.ErrInvalidToken:
			return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
		case auth.ErrInvalidTokenType:
			return nil, status.Error(codes.InvalidArgument, "token is not a refresh token")
		default:
			return nil, status.Error(codes.Internal, "failed to refresh token")
		}
	}

	return &api.AuthServiceRefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64((24 * time.Hour).Seconds()), // Access token expiry in seconds
	}, nil
}
