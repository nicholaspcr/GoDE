package middleware

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(5, 3, 10, 3, 100)

	assert.NotNil(t, rl)
	assert.NotNil(t, rl.globalLimiter)
	assert.NotNil(t, rl.loginLimiters)
	assert.NotNil(t, rl.registerLimiters)
	assert.NotNil(t, rl.userDELimiters)
	assert.Equal(t, 3, rl.maxConcurrentDE)
}

func TestRateLimiter_UnaryGlobalRateLimitMiddleware(t *testing.T) {
	// Create rate limiter with very low limit for testing
	rl := NewRateLimiter(5, 3, 10, 3, 2) // 2 requests per second

	middleware := rl.UnaryGlobalRateLimitMiddleware()
	ctx := context.Background()
	info := &grpc.UnaryServerInfo{
		FullMethod: "/api.v1.TestService/TestMethod",
	}

	handlerCalled := 0
	mockHandler := func(ctx context.Context, req any) (any, error) {
		handlerCalled++
		return "response", nil
	}

	// First two requests should succeed
	resp, err := middleware(ctx, nil, info, mockHandler)
	assert.NoError(t, err)
	assert.Equal(t, "response", resp)
	assert.Equal(t, 1, handlerCalled)

	resp, err = middleware(ctx, nil, info, mockHandler)
	assert.NoError(t, err)
	assert.Equal(t, "response", resp)
	assert.Equal(t, 2, handlerCalled)

	// Third request should be rate limited
	_, err = middleware(ctx, nil, info, mockHandler)
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.ResourceExhausted, st.Code())
	assert.Contains(t, st.Message(), "global rate limit exceeded")
	assert.Equal(t, 2, handlerCalled) // Handler should not be called

	// Wait for rate limiter to allow more requests
	time.Sleep(600 * time.Millisecond)

	resp, err = middleware(ctx, nil, info, mockHandler)
	assert.NoError(t, err)
	assert.Equal(t, "response", resp)
	assert.Equal(t, 3, handlerCalled)
}

func TestRateLimiter_UnaryAuthRateLimitMiddleware(t *testing.T) {
	// Create rate limiter with very low auth limit for testing
	rl := NewRateLimiter(60, 60, 10, 3, 100) // 60 login/register per minute = 1 per second

	middleware := rl.UnaryAuthRateLimitMiddleware()

	handlerCalled := 0
	mockHandler := func(ctx context.Context, req any) (any, error) {
		handlerCalled++
		return "response", nil
	}

	t.Run("non-auth endpoints are not rate limited", func(t *testing.T) {
		ctx := context.Background()
		info := &grpc.UnaryServerInfo{
			FullMethod: "/api.v1.SomeService/SomeMethod",
		}

		// Should not be rate limited
		for i := 0; i < 10; i++ {
			resp, err := middleware(ctx, nil, info, mockHandler)
			assert.NoError(t, err)
			assert.Equal(t, "response", resp)
		}
	})

	t.Run("auth endpoints are rate limited", func(t *testing.T) {
		handlerCalled = 0
		ctx := context.Background()
		info := &grpc.UnaryServerInfo{
			FullMethod: "/api.v1.AuthService/Login",
		}

		// First request should succeed
		resp, err := middleware(ctx, nil, info, mockHandler)
		assert.NoError(t, err)
		assert.Equal(t, "response", resp)
		assert.Equal(t, 1, handlerCalled)

		// Second request should succeed (burst allowance)
		resp, err = middleware(ctx, nil, info, mockHandler)
		assert.NoError(t, err)
		assert.Equal(t, "response", resp)
		assert.Equal(t, 2, handlerCalled)

		// Third request should be rate limited
		_, err = middleware(ctx, nil, info, mockHandler)
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.ResourceExhausted, st.Code())
		assert.Contains(t, st.Message(), "too many login attempts")
		assert.Equal(t, 2, handlerCalled)
	})

	t.Run("register endpoint is also rate limited", func(t *testing.T) {
		handlerCalled = 0
		ctx := context.Background()
		info := &grpc.UnaryServerInfo{
			FullMethod: "/api.v1.AuthService/Register",
		}

		// Allow previous rate limiter to reset
		time.Sleep(2 * time.Second)

		// First request should succeed
		resp, err := middleware(ctx, nil, info, mockHandler)
		assert.NoError(t, err)
		assert.Equal(t, "response", resp)
		assert.Equal(t, 1, handlerCalled)

		// Second request should succeed (burst allowance of 1 for register)
		// Actually, register has burst=1, so second request should fail
		_, err = middleware(ctx, nil, info, mockHandler)
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.ResourceExhausted, st.Code())
		assert.Contains(t, st.Message(), "too many registration attempts")
		assert.Equal(t, 1, handlerCalled)
	})
}

func TestRateLimiter_UnaryDERateLimitMiddleware(t *testing.T) {
	// Create rate limiter with low limits for testing
	rl := NewRateLimiter(5, 3, 60, 2, 100) // 60 DE per minute = 1 per second, max 2 concurrent

	middleware := rl.UnaryDERateLimitMiddleware()

	handlerCalled := 0
	mockHandler := func(ctx context.Context, req any) (any, error) {
		handlerCalled++
		time.Sleep(100 * time.Millisecond) // Simulate work
		return "response", nil
	}

	t.Run("non-DE endpoints are not rate limited", func(t *testing.T) {
		ctx := context.Background()
		info := &grpc.UnaryServerInfo{
			FullMethod: "/api.v1.SomeService/SomeMethod",
		}

		// Should not be rate limited
		for i := 0; i < 10; i++ {
			resp, err := middleware(ctx, nil, info, mockHandler)
			assert.NoError(t, err)
			assert.Equal(t, "response", resp)
		}
	})

	t.Run("DE endpoints without username return auth error", func(t *testing.T) {
		handlerCalled = 0
		ctx := context.Background()
		info := &grpc.UnaryServerInfo{
			FullMethod: "/api.v1.DifferentialEvolutionService/Run",
		}

		// Should return auth error (no username in context)
		resp, err := middleware(ctx, nil, info, mockHandler)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Unauthenticated, status.Code(err))
		assert.Equal(t, 0, handlerCalled) // Handler should not be called
	})

	t.Run("DE endpoints are rate limited per user", func(t *testing.T) {
		handlerCalled = 0

		md := metadata.New(map[string]string{
			"username": "testuser",
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)
		info := &grpc.UnaryServerInfo{
			FullMethod: "/api.v1.DifferentialEvolutionService/Run",
		}

		// First request should succeed
		resp, err := middleware(ctx, nil, info, mockHandler)
		assert.NoError(t, err)
		assert.Equal(t, "response", resp)
		assert.Equal(t, 1, handlerCalled)

		// Second request should succeed (burst allows it)
		resp, err = middleware(ctx, nil, info, mockHandler)
		assert.NoError(t, err)
		assert.Equal(t, "response", resp)
		assert.Equal(t, 2, handlerCalled)

		// Third request should be rate limited (burst=2, rate=1/sec, so 3rd immediate request fails)
		_, err = middleware(ctx, nil, info, mockHandler)
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.ResourceExhausted, st.Code())
		assert.Contains(t, st.Message(), "too many DE execution requests")
		assert.Equal(t, 2, handlerCalled) // Handler should not be called
	})

	t.Run("concurrency limit is enforced", func(t *testing.T) {
		// Create a new rate limiter with higher rate limit but low concurrency
		rl2 := NewRateLimiter(5, 3, 600, 2, 100) // 600 DE per minute = 10 per second, max 2 concurrent
		middleware2 := rl2.UnaryDERateLimitMiddleware()

		handlerCalled = 0

		md := metadata.New(map[string]string{
			"username": "concurrencyuser",
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)
		info := &grpc.UnaryServerInfo{
			FullMethod: "/api.v1.DifferentialEvolutionService/Run",
		}

		// Protect handlerCalled with mutex to prevent data race
		var (
			concurrentCalled int
			mu               sync.Mutex
		)

		slowHandler := func(ctx context.Context, req any) (any, error) {
			mu.Lock()
			concurrentCalled++
			mu.Unlock()
			time.Sleep(500 * time.Millisecond) // Simulate long work
			return "response", nil
		}

		// Start first two concurrent requests (should succeed)
		done1 := make(chan error, 1)
		done2 := make(chan error, 1)
		done3 := make(chan error, 1)

		go func() {
			_, err := middleware2(ctx, nil, info, slowHandler)
			done1 <- err
		}()

		go func() {
			_, err := middleware2(ctx, nil, info, slowHandler)
			done2 <- err
		}()

		// Wait a bit to ensure the first two are running
		time.Sleep(50 * time.Millisecond)

		// Third concurrent request should fail
		go func() {
			_, err := middleware2(ctx, nil, info, slowHandler)
			done3 <- err
		}()

		err3 := <-done3
		require.Error(t, err3)
		st, ok := status.FromError(err3)
		require.True(t, ok)
		assert.Equal(t, codes.ResourceExhausted, st.Code())
		assert.Contains(t, st.Message(), "maximum concurrent DE executions reached")

		// First two should complete successfully
		err1 := <-done1
		err2 := <-done2
		assert.NoError(t, err1)
		assert.NoError(t, err2)

		// Check concurrent handler calls with mutex protection
		mu.Lock()
		count := concurrentCalled
		mu.Unlock()
		assert.Equal(t, 2, count)
	})
}

func TestRateLimiter_Cleanup(t *testing.T) {
	t.Run("recent entries are not cleaned up", func(t *testing.T) {
		rl := NewRateLimiter(5, 3, 10, 3, 100)

		// Create some limiters
		rl.getLoginLimiter("ip1")
		rl.getLoginLimiter("ip2")
		rl.getRegisterLimiter("ip3")
		rl.getUserDELimiter("user1")
		rl.getUserDELimiter("user2")

		assert.Len(t, rl.loginLimiters, 2)
		assert.Len(t, rl.registerLimiters, 1)
		assert.Len(t, rl.userDELimiters, 2)

		// Cleanup with maxAge of 1 hour - recent entries should not be removed
		rl.Cleanup(time.Hour)

		// All entries should still exist since they were created recently
		assert.Len(t, rl.loginLimiters, 2)
		assert.Len(t, rl.registerLimiters, 1)
		assert.Len(t, rl.userDELimiters, 2)
	})

	t.Run("old entries are cleaned up", func(t *testing.T) {
		rl := NewRateLimiter(5, 3, 10, 3, 100)

		// Create some limiters
		rl.getLoginLimiter("ip1")
		rl.getLoginLimiter("ip2")
		rl.getRegisterLimiter("ip3")
		rl.getUserDELimiter("user1")
		rl.getUserDELimiter("user2")

		assert.Len(t, rl.loginLimiters, 2)
		assert.Len(t, rl.registerLimiters, 1)
		assert.Len(t, rl.userDELimiters, 2)

		// Set last access to a time in the past to simulate old entries
		pastTime := time.Now().Add(-2 * time.Hour)
		for _, lwt := range rl.loginLimiters {
			lwt.mu.Lock()
			lwt.lastAccess = pastTime
			lwt.mu.Unlock()
		}
		for _, lwt := range rl.registerLimiters {
			lwt.mu.Lock()
			lwt.lastAccess = pastTime
			lwt.mu.Unlock()
		}
		for _, limiter := range rl.userDELimiters {
			limiter.concurMutex.Lock()
			limiter.lastAccess = pastTime
			limiter.concurMutex.Unlock()
		}

		// Cleanup with maxAge of 1 hour - old entries should be removed
		rl.Cleanup(time.Hour)

		assert.Len(t, rl.loginLimiters, 0)
		assert.Len(t, rl.registerLimiters, 0)
		assert.Len(t, rl.userDELimiters, 0)
	})

	t.Run("DE limiters with active concurrent executions are not cleaned up", func(t *testing.T) {
		rl := NewRateLimiter(5, 3, 10, 3, 100)

		// Create and access a DE limiter
		limiter := rl.getUserDELimiter("user1")

		// Simulate an active execution
		err := limiter.acquireConcurrent()
		assert.NoError(t, err)

		// Set last access to a time in the past
		limiter.concurMutex.Lock()
		limiter.lastAccess = time.Now().Add(-2 * time.Hour)
		limiter.concurMutex.Unlock()

		// Cleanup - entry should NOT be removed because of active execution
		rl.Cleanup(time.Hour)

		assert.Len(t, rl.userDELimiters, 1)

		// Release the execution
		limiter.releaseConcurrent()

		// Now cleanup should remove the entry
		rl.Cleanup(time.Hour)

		assert.Len(t, rl.userDELimiters, 0)
	})
}

func TestGetIPFromContext(t *testing.T) {
	ctx := context.Background()
	ip := getIPFromContext(ctx)
	assert.Equal(t, "unknown", ip)
}

func TestGetUsernameFromContext(t *testing.T) {
	t.Run("no metadata", func(t *testing.T) {
		ctx := context.Background()
		username := getUsernameFromContext(ctx)
		assert.Equal(t, "", username)
	})

	t.Run("no username in metadata", func(t *testing.T) {
		md := metadata.New(map[string]string{
			"other": "value",
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)
		username := getUsernameFromContext(ctx)
		assert.Equal(t, "", username)
	})

	t.Run("username in metadata", func(t *testing.T) {
		md := metadata.New(map[string]string{
			"username": "testuser",
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)
		username := getUsernameFromContext(ctx)
		assert.Equal(t, "testuser", username)
	})
}
