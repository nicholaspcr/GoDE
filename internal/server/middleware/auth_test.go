package middleware

import (
	"context"
	"testing"
	"time"

	"github.com/nicholaspcr/GoDE/internal/server/auth"
	"github.com/stretchr/testify/assert"
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
	middleware := UnaryAuthMiddleware(jwtService)

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
	middleware := UnaryAuthMiddleware(jwtService)

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
	middleware := UnaryAuthMiddleware(jwtService)

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
	middleware := UnaryAuthMiddleware(jwtService)

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
	middleware := UnaryAuthMiddleware(jwtService)

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
	middleware := UnaryAuthMiddleware(jwtService2)

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
	middleware := UnaryAuthMiddleware(jwtService)

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
