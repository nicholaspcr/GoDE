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

// dummyHash is a pre-computed bcrypt hash used to mitigate timing attacks
// when a user doesn't exist. This ensures constant-time comparison even
// for non-existent users by performing bcrypt verification against this hash.
// It's initialized during handler creation to avoid panics at init time.
var dummyHash = mustGenerateDummyHash()

// mustGenerateDummyHash generates a bcrypt hash for timing attack mitigation.
// This uses a pre-computed hash as fallback if generation fails.
func mustGenerateDummyHash() []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte("dummy-password-for-timing-mitigation"), bcrypt.DefaultCost)
	if err != nil {
		// Fallback to a valid pre-computed hash (bcrypt of "dummy-password")
		// This should never happen, but provides a safe fallback instead of panicking
		// Hash generated with: bcrypt.GenerateFromPassword([]byte("dummy-password"), bcrypt.DefaultCost)
		return []byte("$2a$10$N9qo8uLOickgx2ZMRZoMye3IVI564L9ILxI6Jj4Yq1SQXhWKMNXKu")
	}
	return hash
}

// authHandler is responsible for the auth service operations.
type authHandler struct {
	api.UnimplementedAuthServiceServer
	db            store.Store
	jwtService    auth.JWTService
	accessExpiry  time.Duration
}

// NewAuthHandler returns a handle that implements api's authServiceServer.
func NewAuthHandler(st store.Store, jwtService auth.JWTService, accessExpiry time.Duration) Handler {
	return &authHandler{db: st, jwtService: jwtService, accessExpiry: accessExpiry}
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

	// Create a new user object for storage to avoid modifying the request
	// and to prevent potential password exposure in error logs.
	userToCreate := &api.User{
		Ids:      req.User.Ids,
		Email:    req.User.Email,
		Password: string(hashedPassword),
	}

	// Create user
	if err := ah.db.CreateUser(ctx, userToCreate); err != nil {
		// Return a generic error message to avoid leaking internal details
		return nil, status.Error(codes.Internal, "failed to create user")
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
		// Always perform bcrypt comparison to prevent timing-based user enumeration
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(req.Password))
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
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
		ExpiresIn:    int64(ah.accessExpiry.Seconds()), // Access token expiry in seconds
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
		ExpiresIn:    int64(ah.accessExpiry.Seconds()), // Access token expiry in seconds
	}, nil
}
