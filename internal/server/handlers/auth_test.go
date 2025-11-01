package handlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nicholaspcr/GoDE/internal/server/auth"
	"github.com/nicholaspcr/GoDE/internal/store/mock"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

			handler := NewAuthHandler(jwtService)
			handler.SetStore(mockStore)

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
			wantCode: codes.NotFound,
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

			handler := NewAuthHandler(jwtService)
			handler.SetStore(mockStore)

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
					assert.NotEmpty(t, resp.Token)

					// Verify token is valid
					claims, err := jwtService.ValidateToken(resp.Token)
					require.NoError(t, err)
					assert.Equal(t, "testuser", claims.Username)
				}
			}
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", 24*time.Hour)
	mockStore := &mock.MockStore{}

	handler := NewAuthHandler(jwtService)
	handler.SetStore(mockStore)

	// Logout should always succeed (JWT is stateless)
	_, err := handler.(*authHandler).Logout(context.Background(), &api.AuthServiceLogoutRequest{})
	assert.NoError(t, err)
}
