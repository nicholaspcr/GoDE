package handlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nicholaspcr/GoDE/internal/server/auth"
	"github.com/nicholaspcr/GoDE/internal/server/middleware"
	"github.com/nicholaspcr/GoDE/internal/store/mock"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// testRevoker is a simple in-memory TokenRevoker for handler tests.
type testRevoker struct {
	revoked   map[string]bool
	returnErr error
}

func newTestRevoker() *testRevoker { return &testRevoker{revoked: make(map[string]bool)} }

func (r *testRevoker) RevokeToken(_ context.Context, jti string, _ time.Duration) error {
	if r.returnErr != nil {
		return r.returnErr
	}
	r.revoked[jti] = true
	return nil
}

func (r *testRevoker) IsRevoked(_ context.Context, jti string) (bool, error) {
	if r.returnErr != nil {
		return false, r.returnErr
	}
	return r.revoked[jti], nil
}

func TestAuthHandler_Register(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", 24*time.Hour)

	tests := []struct {
		req        *api.AuthServiceRegisterRequest
		setupMock  func(*mock.MockStore)
		name       string
		wantCode   codes.Code
		wantErr    bool
		checkStore bool
	}{
		{
			name: "successful registration",
			req: &api.AuthServiceRegisterRequest{
				User: &api.User{
					Ids:      &api.UserIDs{Username: "testuser"},
					Email:    "test@example.com",
					Password: "validpass123",
				},
			},
			setupMock: func(m *mock.MockStore) {
				m.CreateUserFn = func(ctx context.Context, user *api.User) error {
					// Verify password was hashed
					assert.NotEqual(t, "validpass123", user.Password)
					// Verify email was normalized
					assert.Equal(t, "test@example.com", user.Email)
					return nil
				}
			},
			wantErr:    false,
			checkStore: true,
		},
		{
			name: "invalid username (too short)",
			req: &api.AuthServiceRegisterRequest{
				User: &api.User{
					Ids:      &api.UserIDs{Username: "ab"},
					Email:    "test@example.com",
					Password: "validpass123",
				},
			},
			setupMock: func(m *mock.MockStore) {
				// Should not be called
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "invalid email",
			req: &api.AuthServiceRegisterRequest{
				User: &api.User{
					Ids:      &api.UserIDs{Username: "testuser"},
					Email:    "invalid-email",
					Password: "validpass123",
				},
			},
			setupMock: func(m *mock.MockStore) {
				// Should not be called
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "password too short",
			req: &api.AuthServiceRegisterRequest{
				User: &api.User{
					Ids:      &api.UserIDs{Username: "testuser"},
					Email:    "test@example.com",
					Password: "short",
				},
			},
			setupMock: func(m *mock.MockStore) {
				// Should not be called
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "database error",
			req: &api.AuthServiceRegisterRequest{
				User: &api.User{
					Ids:      &api.UserIDs{Username: "testuser"},
					Email:    "test@example.com",
					Password: "validpass123",
				},
			},
			setupMock: func(m *mock.MockStore) {
				m.CreateUserFn = func(ctx context.Context, user *api.User) error {
					return errors.New("database error")
				}
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "email normalization",
			req: &api.AuthServiceRegisterRequest{
				User: &api.User{
					Ids:      &api.UserIDs{Username: "testuser"},
					Email:    "  TEST@EXAMPLE.COM  ",
					Password: "validpass123",
				},
			},
			setupMock: func(m *mock.MockStore) {
				m.CreateUserFn = func(ctx context.Context, user *api.User) error {
					assert.Equal(t, "test@example.com", user.Email)
					return nil
				}
			},
			wantErr:    false,
			checkStore: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &mock.MockStore{}
			tt.setupMock(mockStore)

			handler := NewAuthHandler(mockStore, jwtService, 15*time.Minute, nil)

			_, err := handler.(*authHandler).Register(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantCode != codes.OK {
					st, ok := status.FromError(err)
					require.True(t, ok)
					assert.Equal(t, tt.wantCode, st.Code())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", 24*time.Hour)

	// Create a hashed password for testing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("validpass123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	tests := []struct {
		req        *api.AuthServiceLoginRequest
		setupMock  func(*mock.MockStore)
		name       string
		wantCode   codes.Code
		wantErr    bool
		checkToken bool
	}{
		{
			name: "successful login",
			req: &api.AuthServiceLoginRequest{
				Username: "testuser",
				Password: "validpass123",
			},
			setupMock: func(m *mock.MockStore) {
				m.GetUserFn = func(ctx context.Context, userIDs *api.UserIDs) (*api.User, error) {
					return &api.User{
						Ids:      &api.UserIDs{Username: "testuser"},
						Password: string(hashedPassword),
					}, nil
				}
			},
			wantErr:    false,
			checkToken: true,
		},
		{
			name: "invalid username (too short)",
			req: &api.AuthServiceLoginRequest{
				Username: "ab",
				Password: "validpass123",
			},
			setupMock: func(m *mock.MockStore) {
				// Should not be called
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "empty password",
			req: &api.AuthServiceLoginRequest{
				Username: "testuser",
				Password: "",
			},
			setupMock: func(m *mock.MockStore) {
				// Should not be called
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "user not found",
			req: &api.AuthServiceLoginRequest{
				Username: "testuser",
				Password: "validpass123",
			},
			setupMock: func(m *mock.MockStore) {
				m.GetUserFn = func(ctx context.Context, userIDs *api.UserIDs) (*api.User, error) {
					return nil, errors.New("user not found")
				}
			},
			wantErr:  true,
			wantCode: codes.Unauthenticated, // Changed from NotFound to prevent user enumeration
		},
		{
			name: "wrong password",
			req: &api.AuthServiceLoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			setupMock: func(m *mock.MockStore) {
				m.GetUserFn = func(ctx context.Context, userIDs *api.UserIDs) (*api.User, error) {
					return &api.User{
						Ids:      &api.UserIDs{Username: "testuser"},
						Password: string(hashedPassword),
					}, nil
				}
			},
			wantErr:  true,
			wantCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &mock.MockStore{}
			tt.setupMock(mockStore)

			handler := NewAuthHandler(mockStore, jwtService, 15*time.Minute, nil)

			resp, err := handler.(*authHandler).Login(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantCode != codes.OK {
					st, ok := status.FromError(err)
					require.True(t, ok)
					assert.Equal(t, tt.wantCode, st.Code())
				}
			} else {
				assert.NoError(t, err)
				if tt.checkToken {
					require.NotNil(t, resp)
					assert.NotEmpty(t, resp.AccessToken)
					assert.NotEmpty(t, resp.RefreshToken)
					assert.Greater(t, resp.ExpiresIn, int64(0))

					// Verify access token is valid
					accessClaims, err := jwtService.ValidateToken(resp.AccessToken)
					require.NoError(t, err)
					assert.Equal(t, "testuser", accessClaims.Username)

					// Verify refresh token is valid
					refreshClaims, err := jwtService.ValidateRefreshToken(resp.RefreshToken)
					require.NoError(t, err)
					assert.Equal(t, "testuser", refreshClaims.Username)
				}
			}
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", 24*time.Hour)
	mockStore := &mock.MockStore{}

	handler := NewAuthHandler(mockStore, jwtService, 15*time.Minute, nil)

	// Logout should always succeed (JWT is stateless)
	_, err := handler.(*authHandler).Logout(context.Background(), &api.AuthServiceLogoutRequest{})
	assert.NoError(t, err)
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret-key-for-refresh-tests", 24*time.Hour)
	mockStore := &mock.MockStore{}

	handler := NewAuthHandler(mockStore, jwtService, 15*time.Minute, nil).(*authHandler)

	t.Run("empty refresh token", func(t *testing.T) {
		_, err := handler.RefreshToken(context.Background(), &api.AuthServiceRefreshTokenRequest{
			RefreshToken: "",
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		_, err := handler.RefreshToken(context.Background(), &api.AuthServiceRefreshTokenRequest{
			RefreshToken: "invalid-token-string",
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})

	t.Run("access token used as refresh token", func(t *testing.T) {
		// Generate a valid access token (not refresh token)
		accessToken, _, err := jwtService.GenerateTokenPair("testuser")
		require.NoError(t, err)

		_, err = handler.RefreshToken(context.Background(), &api.AuthServiceRefreshTokenRequest{
			RefreshToken: accessToken,
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("successful refresh", func(t *testing.T) {
		// Generate valid token pair
		_, refreshToken, err := jwtService.GenerateTokenPair("testuser")
		require.NoError(t, err)

		resp, err := handler.RefreshToken(context.Background(), &api.AuthServiceRefreshTokenRequest{
			RefreshToken: refreshToken,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.Greater(t, resp.ExpiresIn, int64(0))

		// Verify new access token is valid
		claims, err := jwtService.ValidateToken(resp.AccessToken)
		require.NoError(t, err)
		assert.Equal(t, "testuser", claims.Username)
	})
}

func TestAuthHandler_Logout_WithRevocation(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", 15*time.Minute)
	mockStore := &mock.MockStore{}
	revoker := newTestRevoker()
	handler := NewAuthHandler(mockStore, jwtService, 15*time.Minute, revoker).(*authHandler)

	// Generate a token and extract its JTI
	token, err := jwtService.GenerateToken("testuser")
	require.NoError(t, err)
	claims, err := jwtService.ValidateToken(token)
	require.NoError(t, err)

	// Build context with claims (as the auth middleware would)
	ctx := middleware.ContextWithClaims(context.Background(), claims)

	_, err = handler.Logout(ctx, &api.AuthServiceLogoutRequest{Username: "testuser"})
	require.NoError(t, err)

	// JTI should now be in the revocation list
	revoked, err := revoker.IsRevoked(context.Background(), claims.ID)
	require.NoError(t, err)
	assert.True(t, revoked, "access token JTI should be revoked after logout")
}

func TestAuthHandler_Logout_RevokerError(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", 15*time.Minute)
	mockStore := &mock.MockStore{}
	revoker := &testRevoker{revoked: make(map[string]bool), returnErr: errors.New("redis down")}
	handler := NewAuthHandler(mockStore, jwtService, 15*time.Minute, revoker).(*authHandler)

	token, err := jwtService.GenerateToken("testuser")
	require.NoError(t, err)
	claims, err := jwtService.ValidateToken(token)
	require.NoError(t, err)

	ctx := middleware.ContextWithClaims(context.Background(), claims)

	// Logout should succeed even if the revoker errors (best-effort)
	_, err = handler.Logout(ctx, &api.AuthServiceLogoutRequest{Username: "testuser"})
	assert.NoError(t, err)
}

func TestAuthHandler_RefreshToken_RevokesOldToken(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", 15*time.Minute)
	mockStore := &mock.MockStore{}
	revoker := newTestRevoker()
	handler := NewAuthHandler(mockStore, jwtService, 15*time.Minute, revoker).(*authHandler)

	_, refreshToken, err := jwtService.GenerateTokenPair("testuser")
	require.NoError(t, err)
	refreshClaims, err := jwtService.ValidateRefreshToken(refreshToken)
	require.NoError(t, err)
	oldJTI := refreshClaims.ID

	resp, err := handler.RefreshToken(context.Background(), &api.AuthServiceRefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Old refresh token JTI should be revoked (token rotation)
	revoked, err := revoker.IsRevoked(context.Background(), oldJTI)
	require.NoError(t, err)
	assert.True(t, revoked, "old refresh token should be revoked after rotation")

	// New tokens should be valid
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}

func TestAuthHandler_RefreshToken_AlreadyRevoked(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", 15*time.Minute)
	mockStore := &mock.MockStore{}
	revoker := newTestRevoker()
	handler := NewAuthHandler(mockStore, jwtService, 15*time.Minute, revoker).(*authHandler)

	_, refreshToken, err := jwtService.GenerateTokenPair("testuser")
	require.NoError(t, err)
	refreshClaims, err := jwtService.ValidateRefreshToken(refreshToken)
	require.NoError(t, err)

	// Pre-revoke the refresh token
	revoker.revoked[refreshClaims.ID] = true

	_, err = handler.RefreshToken(context.Background(), &api.AuthServiceRefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "revoked")
}
