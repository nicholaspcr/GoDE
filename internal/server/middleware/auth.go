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

	usernameCtxKey struct{} = struct{}{}
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

			// Add user info to context for downstream handlers
			ctx = context.WithValue(ctx, usernameCtxKey, claims.Username)
		}

		return handler(ctx, req)
	}
}
