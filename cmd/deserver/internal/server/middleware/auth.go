package middleware

import (
	"context"
	"strings"

	"github.com/nicholaspcr/GoDE/cmd/deserver/internal/server/session"
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
)

// UnaryAuthMiddleware checks for the Basic authentication and validates if the
// provided token matches with the server's store.
func UnaryAuthMiddleware(
	sessionStore session.Store, ignoreMethods ...string,
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
			token = strings.TrimPrefix(token, "Basic ")

			if !sessionStore.Get(string(token)) {
				return nil, errTokenInvalid
			}
		}

		return handler(ctx, req)
	}
}
