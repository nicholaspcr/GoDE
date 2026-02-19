package middleware

import (
	"context"
	"strings"

	"github.com/nicholaspcr/GoDE/internal/server/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

var (
	errMetadataNotFound = status.Errorf(
		codes.Unauthenticated, "metadata is not provided",
	)
	errTokenNotFound = status.Errorf(
		codes.Unauthenticated,
		"authorization token is not provided",
	)
	errTokenInvalid = status.Errorf(
		codes.Unauthenticated,
		"authorization token is invalid",
	)

	usernameCtxKey = &contextKey{"username"}
	claimsCtxKey   = &contextKey{"claims"}
)

// UnaryAuthMiddleware checks for Bearer authentication and validates JWT tokens.
func UnaryAuthMiddleware(
	jwtService auth.JWTService, ignoreMethods ...string,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		ignoreAuth := false

		method, _ := grpc.Method(ctx)
		for _, imethod := range []string{
			"/api.v1.AuthService/Login",
			"/api.v1.AuthService/Register",
		} {
			if imethod == method {
				ignoreAuth = true
			}
		}

		if !ignoreAuth {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, errMetadataNotFound
			}

			values := md["authorization"]
			if len(values) == 0 {
				return nil, errTokenNotFound
			}

			token := values[0]
			token = strings.TrimPrefix(token, "Bearer ")

			claims, err := jwtService.ValidateToken(token)
			if err != nil {
				return nil, errTokenInvalid
			}

			// Add user info and claims to context for downstream handlers
			ctx = context.WithValue(ctx, usernameCtxKey, claims.Username)
			ctx = context.WithValue(ctx, claimsCtxKey, claims)
		}

		return handler(ctx, req)
	}
}

// UsernameFromContext extracts the authenticated username from the context.
// Returns empty string if no username is found.
func UsernameFromContext(ctx context.Context) string {
	username, ok := ctx.Value(usernameCtxKey).(string)
	if !ok {
		return ""
	}
	return username
}

// ContextWithUsername creates a context with the given username.
// This is primarily for testing purposes.
func ContextWithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameCtxKey, username)
}

// ClaimsFromContext extracts the JWT claims from the context.
// Returns nil if no claims are found.
func ClaimsFromContext(ctx context.Context) *auth.Claims {
	claims, ok := ctx.Value(claimsCtxKey).(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}

// ContextWithClaims creates a context with the given claims.
// This is primarily for testing purposes.
func ContextWithClaims(ctx context.Context, claims *auth.Claims) context.Context {
	ctx = context.WithValue(ctx, claimsCtxKey, claims)
	if claims != nil {
		ctx = context.WithValue(ctx, usernameCtxKey, claims.Username)
	}
	return ctx
}

// StreamAuthMiddleware checks for Bearer authentication and validates JWT tokens on stream RPCs.
func StreamAuthMiddleware(jwtService auth.JWTService) grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := ss.Context()

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return errMetadataNotFound
		}

		values := md["authorization"]
		if len(values) == 0 {
			return errTokenNotFound
		}

		token := values[0]
		token = strings.TrimPrefix(token, "Bearer ")

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			return errTokenInvalid
		}

		// Add user info and claims to context for downstream handlers
		ctx = context.WithValue(ctx, usernameCtxKey, claims.Username)
		ctx = context.WithValue(ctx, claimsCtxKey, claims)

		// Wrap the stream with the new context
		wrapped := &wrappedServerStream{ServerStream: ss, ctx: ctx}
		return handler(srv, wrapped)
	}
}

// wrappedServerStream wraps a grpc.ServerStream to override the context.
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context returns the wrapped context.
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

// RequireScope checks if the current user has the required scope.
// Returns a permission denied error if the scope is missing.
func RequireScope(ctx context.Context, scope auth.Scope) error {
	claims := ClaimsFromContext(ctx)
	if claims == nil {
		return status.Error(codes.Unauthenticated, "not authenticated")
	}

	if !claims.HasScope(scope) {
		return status.Errorf(codes.PermissionDenied, "insufficient permissions: requires scope %s", scope)
	}

	return nil
}

// RequireAnyScope checks if the current user has any of the required scopes.
// Returns a permission denied error if none of the scopes are present.
func RequireAnyScope(ctx context.Context, scopes ...auth.Scope) error {
	claims := ClaimsFromContext(ctx)
	if claims == nil {
		return status.Error(codes.Unauthenticated, "not authenticated")
	}

	if !claims.HasAnyScope(scopes...) {
		return status.Errorf(codes.PermissionDenied, "insufficient permissions: requires one of %v", scopes)
	}

	return nil
}

// HasScope checks if the current user has the specified scope.
// Returns false if not authenticated or scope is missing.
func HasScope(ctx context.Context, scope auth.Scope) bool {
	claims := ClaimsFromContext(ctx)
	if claims == nil {
		return false
	}
	return claims.HasScope(scope)
}
