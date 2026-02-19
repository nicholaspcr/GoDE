package handlers

import (
	"context"

	"github.com/nicholaspcr/GoDE/internal/server/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// usernameFromContext extracts the authenticated username from the context.
// Returns a gRPC Unauthenticated error if the username is not present.
func usernameFromContext(ctx context.Context) (string, error) {
	userID := middleware.UsernameFromContext(ctx)
	if userID == "" {
		return "", status.Error(codes.Unauthenticated, "user not authenticated")
	}
	return userID, nil
}
