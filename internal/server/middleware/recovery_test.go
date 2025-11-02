package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnaryPanicRecoveryMiddleware(t *testing.T) {
	middleware := UnaryPanicRecoveryMiddleware()

	t.Run("normal execution without panic", func(t *testing.T) {
		ctx := context.Background()
		info := &grpc.UnaryServerInfo{
			FullMethod: "/api.v1.TestService/TestMethod",
		}

		mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "response", nil
		}

		resp, err := middleware(ctx, nil, info, mockHandler)
		assert.NoError(t, err)
		assert.Equal(t, "response", resp)
	})

	t.Run("panic is recovered and returns internal error", func(t *testing.T) {
		ctx := context.Background()
		info := &grpc.UnaryServerInfo{
			FullMethod: "/api.v1.TestService/TestMethod",
		}

		mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
			panic("something went wrong")
		}

		resp, err := middleware(ctx, nil, info, mockHandler)
		assert.Nil(t, resp)
		require.Error(t, err)

		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "internal server error")
	})

	t.Run("panic with nil value", func(t *testing.T) {
		ctx := context.Background()
		info := &grpc.UnaryServerInfo{
			FullMethod: "/api.v1.TestService/TestMethod",
		}

		mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
			panic(nil)
		}

		resp, err := middleware(ctx, nil, info, mockHandler)
		assert.Nil(t, resp)
		require.Error(t, err)

		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
	})

	t.Run("panic with custom error", func(t *testing.T) {
		ctx := context.Background()
		info := &grpc.UnaryServerInfo{
			FullMethod: "/api.v1.TestService/TestMethod",
		}

		mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
			panic(assert.AnError)
		}

		resp, err := middleware(ctx, nil, info, mockHandler)
		assert.Nil(t, resp)
		require.Error(t, err)

		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
	})

	t.Run("handler returns error normally", func(t *testing.T) {
		ctx := context.Background()
		info := &grpc.UnaryServerInfo{
			FullMethod: "/api.v1.TestService/TestMethod",
		}

		mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, assert.AnError
		}

		resp, err := middleware(ctx, nil, info, mockHandler)
		assert.Nil(t, resp)
		assert.Equal(t, assert.AnError, err)
	})
}

type mockServerStream struct {
	grpc.ServerStream
}

func TestStreamPanicRecoveryMiddleware(t *testing.T) {
	middleware := StreamPanicRecoveryMiddleware()

	t.Run("normal execution without panic", func(t *testing.T) {
		info := &grpc.StreamServerInfo{
			FullMethod: "/api.v1.TestService/TestStream",
		}

		mockHandler := func(srv interface{}, ss grpc.ServerStream) error {
			return nil
		}

		err := middleware(nil, &mockServerStream{}, info, mockHandler)
		assert.NoError(t, err)
	})

	t.Run("panic is recovered and returns internal error", func(t *testing.T) {
		info := &grpc.StreamServerInfo{
			FullMethod: "/api.v1.TestService/TestStream",
		}

		mockHandler := func(srv interface{}, ss grpc.ServerStream) error {
			panic("stream panic")
		}

		err := middleware(nil, &mockServerStream{}, info, mockHandler)
		require.Error(t, err)

		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "internal server error")
	})

	t.Run("handler returns error normally", func(t *testing.T) {
		info := &grpc.StreamServerInfo{
			FullMethod: "/api.v1.TestService/TestStream",
		}

		mockHandler := func(srv interface{}, ss grpc.ServerStream) error {
			return assert.AnError
		}

		err := middleware(nil, &mockServerStream{}, info, mockHandler)
		assert.Equal(t, assert.AnError, err)
	})
}

func TestRecoverDEExecution(t *testing.T) {
	t.Run("normal execution without panic", func(t *testing.T) {
		defer RecoverDEExecution(1)
		// Should not panic
	})

	t.Run("panic is recovered", func(t *testing.T) {
		executed := false
		func() {
			defer RecoverDEExecution(1)
			executed = true
			panic("DE execution panic")
		}()

		// Should have executed before panic
		assert.True(t, executed)
		// Test passes if we reach here (panic was recovered)
	})
}

func TestSafeGo(t *testing.T) {
	t.Run("normal execution without panic", func(t *testing.T) {
		done := make(chan bool)
		SafeGo("test", func() {
			done <- true
		})

		select {
		case <-done:
			// Success
		case <-context.Background().Done():
			t.Fatal("goroutine did not complete")
		}
	})

	t.Run("panic is recovered", func(t *testing.T) {
		done := make(chan bool)
		SafeGo("test-panic", func() {
			defer func() {
				done <- true
			}()
			panic("goroutine panic")
		})

		select {
		case <-done:
			// Success - panic was recovered and deferred function ran
		case <-context.Background().Done():
			t.Fatal("goroutine did not complete")
		}
	})
}
