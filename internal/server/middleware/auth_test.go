package middleware

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nicholaspcr/GoDE/internal/server/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestUnaryAuthMiddleware_IgnoredMethods(t *testing.T) {
	// Note: Testing ignored methods is complex because grpc.Method(ctx) uses internal
	// context keys that can't be easily mocked. The auth bypass logic in the middleware
	// checks hardcoded method names ("/api.v1.AuthService/Login", "/api.v1.AuthService/Register").
	//
	// In a real gRPC server, the method is automatically added to the context by the gRPC
	// framework. Since we can't easily replicate this in unit tests without using the full
	// gRPC stack, we verify the other critical paths instead.
	//
	// This functionality is better tested in integration tests where the full gRPC server
	// is running and can properly populate the context.
	t.Skip("Skipping ignored methods test - requires full gRPC server context")
}

func TestUnaryAuthMiddleware_MissingMetadata(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", 15*time.Minute)
	middleware := UnaryAuthMiddleware(jwtService, nil)

	mockHandler := func(ctx context.Context, req any) (any, error) {
		t.Fatal("handler should not be called")
		return nil, nil
	}

	ctx := context.Background()
	info := &grpc.UnaryServerInfo{
		FullMethod: "/api.v1.SomeService/ProtectedMethod",
	}

	_, err := middleware(ctx, nil, info, mockHandler)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "metadata is not provided")
}

func TestUnaryAuthMiddleware_MissingToken(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", 15*time.Minute)
	middleware := UnaryAuthMiddleware(jwtService, nil)

	mockHandler := func(ctx context.Context, req any) (any, error) {
		t.Fatal("handler should not be called")
		return nil, nil
	}

	// Create context with metadata but no authorization header
	md := metadata.New(map[string]string{
		"other-header": "value",
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	info := &grpc.UnaryServerInfo{
		FullMethod: "/api.v1.SomeService/ProtectedMethod",
	}

	_, err := middleware(ctx, nil, info, mockHandler)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "authorization token is not provided")
}

func TestUnaryAuthMiddleware_InvalidToken(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", 15*time.Minute)
	middleware := UnaryAuthMiddleware(jwtService, nil)

	mockHandler := func(ctx context.Context, req any) (any, error) {
		t.Fatal("handler should not be called")
		return nil, nil
	}

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "malformed token",
			token: "invalid.token.here",
		},
		{
			name:  "empty token",
			token: "",
		},
		{
			name:  "random string",
			token: "Bearer randomstring",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := metadata.New(map[string]string{
				"authorization": tt.token,
			})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			info := &grpc.UnaryServerInfo{
				FullMethod: "/api.v1.SomeService/ProtectedMethod",
			}

			_, err := middleware(ctx, nil, info, mockHandler)

			assert.Error(t, err)
			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, codes.Unauthenticated, st.Code())
			assert.Contains(t, st.Message(), "authorization token is invalid")
		})
	}
}

func TestUnaryAuthMiddleware_ExpiredToken(t *testing.T) {
	// Create JWT service with very short expiration
	jwtService := auth.NewJWTService("test-secret", 100*time.Millisecond)
	middleware := UnaryAuthMiddleware(jwtService, nil)

	// Generate token
	token, err := jwtService.GenerateToken("testuser")
	assert.NoError(t, err)

	// Wait for token to expire
	time.Sleep(200 * time.Millisecond)

	mockHandler := func(ctx context.Context, req any) (any, error) {
		t.Fatal("handler should not be called")
		return nil, nil
	}

	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	info := &grpc.UnaryServerInfo{
		FullMethod: "/api.v1.SomeService/ProtectedMethod",
	}

	_, err = middleware(ctx, nil, info, mockHandler)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "authorization token is invalid")
}

func TestUnaryAuthMiddleware_ValidToken(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", 15*time.Minute)
	middleware := UnaryAuthMiddleware(jwtService, nil)

	// Generate valid token
	token, err := jwtService.GenerateToken("testuser")
	assert.NoError(t, err)

	tests := []struct {
		name       string
		authHeader string
	}{
		{
			name:       "token with Bearer prefix",
			authHeader: "Bearer " + token,
		},
		{
			name:       "token without Bearer prefix",
			authHeader: token,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerCalled := false
			var receivedCtx context.Context
			mockHandler := func(ctx context.Context, req any) (any, error) {
				handlerCalled = true
				receivedCtx = ctx
				return "success", nil
			}

			md := metadata.New(map[string]string{
				"authorization": tt.authHeader,
			})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			info := &grpc.UnaryServerInfo{
				FullMethod: "/api.v1.SomeService/ProtectedMethod",
			}

			resp, err := middleware(ctx, nil, info, mockHandler)

			assert.NoError(t, err)
			assert.Equal(t, "success", resp)
			assert.True(t, handlerCalled, "handler should have been called")

			// Verify username was added to context
			username := receivedCtx.Value(usernameCtxKey)
			assert.Equal(t, "testuser", username)
		})
	}
}

func TestUnaryAuthMiddleware_WrongSecret(t *testing.T) {
	// Create token with one secret
	jwtService1 := auth.NewJWTService("secret1", 15*time.Minute)
	token, err := jwtService1.GenerateToken("testuser")
	assert.NoError(t, err)

	// Try to validate with different secret
	jwtService2 := auth.NewJWTService("secret2", 15*time.Minute)
	middleware := UnaryAuthMiddleware(jwtService2, nil)

	mockHandler := func(ctx context.Context, req any) (any, error) {
		t.Fatal("handler should not be called")
		return nil, nil
	}

	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	info := &grpc.UnaryServerInfo{
		FullMethod: "/api.v1.SomeService/ProtectedMethod",
	}

	_, err = middleware(ctx, nil, info, mockHandler)

	assert.Error(t, err)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
	assert.Contains(t, st.Message(), "authorization token is invalid")
}

func TestUnaryAuthMiddleware_MultipleAuthorizationHeaders(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", 15*time.Minute)
	middleware := UnaryAuthMiddleware(jwtService, nil)

	// Generate valid token
	token, err := jwtService.GenerateToken("testuser")
	assert.NoError(t, err)

	handlerCalled := false
	mockHandler := func(ctx context.Context, req any) (any, error) {
		handlerCalled = true
		return "success", nil
	}

	// Create metadata with multiple authorization headers (gRPC allows this)
	md := metadata.MD{
		"authorization": []string{"Bearer " + token, "Bearer invalid"},
	}
	ctx := metadata.NewIncomingContext(context.Background(), md)
	info := &grpc.UnaryServerInfo{
		FullMethod: "/api.v1.SomeService/ProtectedMethod",
	}

	resp, err := middleware(ctx, nil, info, mockHandler)

	// Should use the first token
	assert.NoError(t, err)
	assert.Equal(t, "success", resp)
	assert.True(t, handlerCalled)
}
// stubRevoker is a simple TokenRevoker for testing.
type stubRevoker struct {
	revokedJTIs map[string]bool
	revokeErr   error
}

func (s *stubRevoker) RevokeToken(_ context.Context, jti string, _ time.Duration) error {
	if s.revokedJTIs == nil {
		s.revokedJTIs = make(map[string]bool)
	}
	s.revokedJTIs[jti] = true
	return s.revokeErr
}

func (s *stubRevoker) IsRevoked(_ context.Context, jti string) (bool, error) {
	if s.revokeErr != nil {
		return false, s.revokeErr
	}
	return s.revokedJTIs[jti], nil
}

func TestUnaryAuthMiddleware_RevokedToken(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", 15*time.Minute)
	token, err := jwtService.GenerateToken("testuser")
	require.NoError(t, err)

	// Validate once to extract the JTI
	claims, err := jwtService.ValidateToken(token)
	require.NoError(t, err)

	revoker := &stubRevoker{revokedJTIs: map[string]bool{claims.ID: true}}
	mw := UnaryAuthMiddleware(jwtService, revoker)

	mockHandler := func(ctx context.Context, req any) (any, error) {
		t.Fatal("handler should not be called for a revoked token")
		return nil, nil
	}

	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	info := &grpc.UnaryServerInfo{FullMethod: "/api.v1.SomeService/ProtectedMethod"}

	_, reqErr := mw(ctx, nil, info, mockHandler)

	assert.Error(t, reqErr)
	st, ok := status.FromError(reqErr)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestUnaryAuthMiddleware_RevokerError(t *testing.T) {
	jwtService := auth.NewJWTService("test-secret", 15*time.Minute)
	token, err := jwtService.GenerateToken("testuser")
	require.NoError(t, err)

	// Revoker returns an error — middleware should allow through (fail open)
	revoker := &stubRevoker{revokeErr: errors.New("redis unavailable")}
	mw := UnaryAuthMiddleware(jwtService, revoker)

	handlerCalled := false
	mockHandler := func(ctx context.Context, req any) (any, error) {
		handlerCalled = true
		return "ok", nil
	}

	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	info := &grpc.UnaryServerInfo{FullMethod: "/api.v1.SomeService/ProtectedMethod"}

	resp, reqErr := mw(ctx, nil, info, mockHandler)

	assert.NoError(t, reqErr)
	assert.Equal(t, "ok", resp)
	assert.True(t, handlerCalled, "should allow through when revoker errors")
}

func TestClaimsFromContext(t *testing.T) {
	t.Run("returns nil when no claims in context", func(t *testing.T) {
		ctx := context.Background()
		claims := ClaimsFromContext(ctx)
		assert.Nil(t, claims)
	})

	t.Run("returns claims when present in context", func(t *testing.T) {
		ctx := context.Background()
		expectedClaims := &auth.Claims{
			Username: "testuser",
			Scopes:   []auth.Scope{auth.ScopeDERun, auth.ScopeDERead},
		}
		ctx = ContextWithClaims(ctx, expectedClaims)

		claims := ClaimsFromContext(ctx)
		require.NotNil(t, claims)
		assert.Equal(t, "testuser", claims.Username)
		assert.Contains(t, claims.Scopes, auth.ScopeDERun)
		assert.Contains(t, claims.Scopes, auth.ScopeDERead)
	})
}

func TestContextWithClaims(t *testing.T) {
	t.Run("stores claims in context", func(t *testing.T) {
		ctx := context.Background()
		claims := &auth.Claims{
			Username: "testuser",
			Scopes:   []auth.Scope{auth.ScopeDERun},
		}

		ctx = ContextWithClaims(ctx, claims)

		retrievedClaims := ClaimsFromContext(ctx)
		require.NotNil(t, retrievedClaims)
		assert.Equal(t, claims.Username, retrievedClaims.Username)
		assert.Equal(t, claims.Scopes, retrievedClaims.Scopes)

		// Should also set username in context
		username := UsernameFromContext(ctx)
		assert.Equal(t, "testuser", username)
	})

	t.Run("handles nil claims", func(t *testing.T) {
		ctx := context.Background()
		ctx = ContextWithClaims(ctx, nil)

		claims := ClaimsFromContext(ctx)
		assert.Nil(t, claims)
	})
}

func TestRequireScope(t *testing.T) {
	tests := []struct {
		name        string
		scopes      []auth.Scope
		required    auth.Scope
		expectError bool
		errorCode   codes.Code
	}{
		{
			name:        "user has required scope",
			scopes:      []auth.Scope{auth.ScopeDERun, auth.ScopeDERead},
			required:    auth.ScopeDERun,
			expectError: false,
		},
		{
			name:        "user has admin scope (grants all)",
			scopes:      []auth.Scope{auth.ScopeAdmin},
			required:    auth.ScopeDERun,
			expectError: false,
		},
		{
			name:        "user missing required scope",
			scopes:      []auth.Scope{auth.ScopeDERead},
			required:    auth.ScopeDERun,
			expectError: true,
			errorCode:   codes.PermissionDenied,
		},
		{
			name:        "user has no scopes",
			scopes:      []auth.Scope{},
			required:    auth.ScopeDERun,
			expectError: true,
			errorCode:   codes.PermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			claims := &auth.Claims{
				Username: "testuser",
				Scopes:   tt.scopes,
			}
			ctx = ContextWithClaims(ctx, claims)

			err := RequireScope(ctx, tt.required)

			if tt.expectError {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.errorCode, st.Code())
			} else {
				assert.NoError(t, err)
			}
		})
	}

	t.Run("not authenticated", func(t *testing.T) {
		ctx := context.Background() // No claims

		err := RequireScope(ctx, auth.ScopeDERun)

		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestRequireAnyScope(t *testing.T) {
	tests := []struct {
		name        string
		scopes      []auth.Scope
		required    []auth.Scope
		expectError bool
		errorCode   codes.Code
	}{
		{
			name:        "user has one of required scopes",
			scopes:      []auth.Scope{auth.ScopeDERun},
			required:    []auth.Scope{auth.ScopeDERun, auth.ScopeDERead},
			expectError: false,
		},
		{
			name:        "user has all required scopes",
			scopes:      []auth.Scope{auth.ScopeDERun, auth.ScopeDERead},
			required:    []auth.Scope{auth.ScopeDERun, auth.ScopeDERead},
			expectError: false,
		},
		{
			name:        "user has admin scope",
			scopes:      []auth.Scope{auth.ScopeAdmin},
			required:    []auth.Scope{auth.ScopeDERun, auth.ScopeDERead},
			expectError: false,
		},
		{
			name:        "user has none of required scopes",
			scopes:      []auth.Scope{auth.ScopeUserRead},
			required:    []auth.Scope{auth.ScopeDERun, auth.ScopeDERead},
			expectError: true,
			errorCode:   codes.PermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			claims := &auth.Claims{
				Username: "testuser",
				Scopes:   tt.scopes,
			}
			ctx = ContextWithClaims(ctx, claims)

			err := RequireAnyScope(ctx, tt.required...)

			if tt.expectError {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.errorCode, st.Code())
			} else {
				assert.NoError(t, err)
			}
		})
	}

	t.Run("not authenticated", func(t *testing.T) {
		ctx := context.Background() // No claims

		err := RequireAnyScope(ctx, auth.ScopeDERun, auth.ScopeDERead)

		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})
}

func TestHasScope(t *testing.T) {
	tests := []struct {
		name     string
		scopes   []auth.Scope
		check    auth.Scope
		expected bool
	}{
		{
			name:     "user has scope",
			scopes:   []auth.Scope{auth.ScopeDERun, auth.ScopeDERead},
			check:    auth.ScopeDERun,
			expected: true,
		},
		{
			name:     "user has admin scope (grants all)",
			scopes:   []auth.Scope{auth.ScopeAdmin},
			check:    auth.ScopeDERun,
			expected: true,
		},
		{
			name:     "user missing scope",
			scopes:   []auth.Scope{auth.ScopeDERead},
			check:    auth.ScopeDERun,
			expected: false,
		},
		{
			name:     "user has no scopes",
			scopes:   []auth.Scope{},
			check:    auth.ScopeDERun,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			claims := &auth.Claims{
				Username: "testuser",
				Scopes:   tt.scopes,
			}
			ctx = ContextWithClaims(ctx, claims)

			result := HasScope(ctx, tt.check)
			assert.Equal(t, tt.expected, result)
		})
	}

	t.Run("not authenticated returns false", func(t *testing.T) {
		ctx := context.Background() // No claims

		result := HasScope(ctx, auth.ScopeDERun)
		assert.False(t, result)
	})
}
