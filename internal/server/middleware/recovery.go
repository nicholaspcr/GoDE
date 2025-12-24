package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryPanicRecoveryMiddleware recovers from panics in unary handlers.
func UnaryPanicRecoveryMiddleware() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic with stack trace
				stack := debug.Stack()
				slog.Error("panic recovered in gRPC handler",
					slog.String("method", info.FullMethod),
					slog.Any("panic", r),
					slog.String("stack", string(stack)),
				)

				// Return internal error to client
				err = status.Errorf(codes.Internal,
					"internal server error occurred")
			}
		}()

		resp, err = handler(ctx, req)
		return resp, err
	}
}

// StreamPanicRecoveryMiddleware recovers from panics in stream handlers.
func StreamPanicRecoveryMiddleware() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic with stack trace
				stack := debug.Stack()
				slog.Error("panic recovered in gRPC stream handler",
					slog.String("method", info.FullMethod),
					slog.Any("panic", r),
					slog.String("stack", string(stack)),
				)

				// Return internal error to client
				err = status.Errorf(codes.Internal,
					"internal server error occurred")
			}
		}()

		err = handler(srv, ss)
		return err
	}
}

// RecoverDEExecution wraps a DE execution goroutine with panic recovery.
// It should be used as: defer RecoverDEExecution(executionID)
func RecoverDEExecution(executionID int) {
	if r := recover(); r != nil {
		// Log the panic with stack trace
		stack := debug.Stack()
		slog.Error("panic recovered in DE execution goroutine",
			slog.Int("execution_id", executionID),
			slog.Any("panic", r),
			slog.String("stack", string(stack)),
		)
	}
}

// SafeGo runs a function in a goroutine with panic recovery.
// If the function panics, it logs the error but doesn't crash the program.
func SafeGo(name string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				slog.Error(fmt.Sprintf("panic recovered in goroutine '%s'", name),
					slog.Any("panic", r),
					slog.String("stack", string(stack)),
				)
			}
		}()
		fn()
	}()
}
