package auth

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
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

// UnaryMiddleware checks for the Basic authentication and validates if the
// provided token matches with the server's store.
func UnaryMiddleware(st store.UserOperations) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
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
		b, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return nil, err
		}

		kv := strings.Split(string(b), ":")
		if len(kv) != 2 {
			return nil, errTokenInvalid
		}
		username := kv[0]
		password := kv[1]

		usr, err := st.GetUser(ctx, &api.UserIDs{Email: username})
		if err != nil {
			return nil, err
		}

		if usr.Password != password {
			return nil, errTokenInvalid
		}

		return handler(ctx, req)
	}
}
