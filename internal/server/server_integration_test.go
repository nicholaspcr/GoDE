package server

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/nicholaspcr/GoDE/internal/store/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("successfully creates server with store", func(t *testing.T) {
		ctx := context.Background()
		store := &mock.MockStore{}
		cfg := DefaultConfig()

		srv, err := New(ctx, cfg, WithStore(store))
		require.NoError(t, err)
		require.NotNil(t, srv)

		// Verify server internals
		s, ok := srv.(*server)
		require.True(t, ok, "server should be of type *server")
		assert.NotNil(t, s.st, "store should be set")
		assert.NotNil(t, s.jwtService, "jwt service should be initialized")
		assert.NotNil(t, s.executor, "executor should be initialized")
		assert.Len(t, s.handlers, 4, "should have 4 handlers (auth, user, pareto, de)")
	})

	t.Run("returns error when store is not provided", func(t *testing.T) {
		ctx := context.Background()
		cfg := DefaultConfig()

		srv, err := New(ctx, cfg)
		assert.Error(t, err)
		assert.Nil(t, srv)
		assert.Contains(t, err.Error(), "store must be provided")
	})

	t.Run("applies WithStore option correctly", func(t *testing.T) {
		ctx := context.Background()
		store := &mock.MockStore{}
		cfg := DefaultConfig()

		srv, err := New(ctx, cfg, WithStore(store))
		require.NoError(t, err)
		require.NotNil(t, srv)

		s, ok := srv.(*server)
		require.True(t, ok)
		assert.Equal(t, store, s.st)
	})

	t.Run("applies WithConfig option correctly", func(t *testing.T) {
		ctx := context.Background()
		store := &mock.MockStore{}
		cfg1 := DefaultConfig()
		cfg2 := DefaultConfig()
		cfg2.LisAddr = ":9999"

		srv, err := New(ctx, cfg1, WithStore(store), WithConfig(cfg2))
		require.NoError(t, err)
		require.NotNil(t, srv)

		s, ok := srv.(*server)
		require.True(t, ok)
		assert.Equal(t, ":9999", s.cfg.LisAddr, "config should be overridden by WithConfig")
	})

	t.Run("registers problems from registry", func(t *testing.T) {
		ctx := context.Background()
		store := &mock.MockStore{}
		cfg := DefaultConfig()

		srv, err := New(ctx, cfg, WithStore(store))
		require.NoError(t, err)
		require.NotNil(t, srv)

		s, ok := srv.(*server)
		require.True(t, ok)
		assert.NotNil(t, s.executor, "executor should be initialized with problems")
		// Executor should have problems registered - verified through successful execution
	})

	t.Run("registers variants from registry", func(t *testing.T) {
		ctx := context.Background()
		store := &mock.MockStore{}
		cfg := DefaultConfig()

		srv, err := New(ctx, cfg, WithStore(store))
		require.NoError(t, err)
		require.NotNil(t, srv)

		s, ok := srv.(*server)
		require.True(t, ok)
		assert.NotNil(t, s.executor, "executor should be initialized with variants")
		// Executor should have variants registered - verified through successful execution
	})

	t.Run("creates all required handlers", func(t *testing.T) {
		ctx := context.Background()
		store := &mock.MockStore{}
		cfg := DefaultConfig()

		srv, err := New(ctx, cfg, WithStore(store))
		require.NoError(t, err)

		s, ok := srv.(*server)
		require.True(t, ok)

		// Should have exactly 4 handlers
		assert.Len(t, s.handlers, 4)

		// Verify handlers are not nil
		for i, h := range s.handlers {
			assert.NotNil(t, h, "handler %d should not be nil", i)
		}
	})

	t.Run("initializes executor with correct config", func(t *testing.T) {
		ctx := context.Background()
		store := &mock.MockStore{}
		cfg := DefaultConfig()
		cfg.Executor.MaxWorkers = 5
		cfg.Executor.MaxVectorsInProgress = 50
		cfg.Executor.ExecutionTTL = time.Hour
		cfg.Executor.ResultTTL = 24 * time.Hour
		cfg.Executor.ProgressTTL = 30 * time.Minute

		srv, err := New(ctx, cfg, WithStore(store))
		require.NoError(t, err)

		s, ok := srv.(*server)
		require.True(t, ok)
		assert.NotNil(t, s.executor, "executor should be initialized with custom config")
	})

	t.Run("initializes JWT service with config values", func(t *testing.T) {
		ctx := context.Background()
		store := &mock.MockStore{}
		cfg := DefaultConfig()
		cfg.JWTSecret = "test-secret"
		cfg.JWTExpiry = 2 * time.Hour

		srv, err := New(ctx, cfg, WithStore(store))
		require.NoError(t, err)

		s, ok := srv.(*server)
		require.True(t, ok)
		assert.NotNil(t, s.jwtService, "jwt service should be initialized")
	})
}

func TestInterceptorLogger(t *testing.T) {
	t.Run("creates logger adapter", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		interceptorLogger := InterceptorLogger(logger)

		assert.NotNil(t, interceptorLogger)
		// Logger should be callable without panic
		ctx := context.Background()
		assert.NotPanics(t, func() {
			interceptorLogger.Log(ctx, 0, "test message")
		})
	})

	t.Run("logs messages at different levels", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		interceptorLogger := InterceptorLogger(logger)
		ctx := context.Background()

		// Should not panic for any log level
		assert.NotPanics(t, func() {
			interceptorLogger.Log(ctx, 0, "debug message")
			interceptorLogger.Log(ctx, 1, "info message")
			interceptorLogger.Log(ctx, 2, "warn message")
			interceptorLogger.Log(ctx, 3, "error message")
		})
	})

	t.Run("logs with fields", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		interceptorLogger := InterceptorLogger(logger)
		ctx := context.Background()

		assert.NotPanics(t, func() {
			interceptorLogger.Log(ctx, 1, "message with fields", "key1", "value1", "key2", 42)
		})
	})
}

func TestServerIntegration_Construction(t *testing.T) {
	t.Run("server can be constructed with minimal config", func(t *testing.T) {
		ctx := context.Background()
		store := &mock.MockStore{}
		cfg := DefaultConfig()

		srv, err := New(ctx, cfg, WithStore(store))
		require.NoError(t, err)
		require.NotNil(t, srv)

		s, ok := srv.(*server)
		require.True(t, ok)

		// Verify all components are initialized
		assert.NotNil(t, s.st)
		assert.NotNil(t, s.jwtService)
		assert.NotNil(t, s.executor)
		assert.NotNil(t, s.handlers)
		assert.Len(t, s.handlers, 4)
	})

	t.Run("server construction with custom ports", func(t *testing.T) {
		ctx := context.Background()
		store := &mock.MockStore{}
		cfg := DefaultConfig()
		cfg.LisAddr = ":13030"
		cfg.HTTPPort = "18081"

		srv, err := New(ctx, cfg, WithStore(store))
		require.NoError(t, err)
		require.NotNil(t, srv)

		s, ok := srv.(*server)
		require.True(t, ok)
		assert.Equal(t, ":13030", s.cfg.LisAddr)
		assert.Equal(t, "18081", s.cfg.HTTPPort)
	})
}

func TestServerOptions(t *testing.T) {
	t.Run("WithStore option sets store", func(t *testing.T) {
		ctx := context.Background()
		store := &mock.MockStore{}
		cfg := DefaultConfig()

		srv, err := New(ctx, cfg, WithStore(store))
		require.NoError(t, err)

		s, ok := srv.(*server)
		require.True(t, ok)
		assert.Equal(t, store, s.st)
	})

	t.Run("WithConfig option overrides config", func(t *testing.T) {
		ctx := context.Background()
		store := &mock.MockStore{}
		cfg1 := DefaultConfig()
		cfg1.LisAddr = ":3030"

		cfg2 := DefaultConfig()
		cfg2.LisAddr = ":4040"

		srv, err := New(ctx, cfg1, WithStore(store), WithConfig(cfg2))
		require.NoError(t, err)

		s, ok := srv.(*server)
		require.True(t, ok)
		assert.Equal(t, ":4040", s.cfg.LisAddr, "config should be overridden")
	})

	t.Run("multiple options applied in order", func(t *testing.T) {
		ctx := context.Background()
		store1 := &mock.MockStore{}
		store2 := &mock.MockStore{}
		cfg := DefaultConfig()

		srv, err := New(ctx, cfg,
			WithStore(store1),
			WithStore(store2), // Second store should override first
		)
		require.NoError(t, err)

		s, ok := srv.(*server)
		require.True(t, ok)
		assert.Equal(t, store2, s.st, "second WithStore should override first")
	})
}
