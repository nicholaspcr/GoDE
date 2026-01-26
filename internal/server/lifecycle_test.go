package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/internal/server/handlers"
	"github.com/nicholaspcr/GoDE/internal/server/middleware"
	"github.com/nicholaspcr/GoDE/internal/slo"
	"github.com/nicholaspcr/GoDE/internal/store/mock"
	"github.com/nicholaspcr/GoDE/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// mockExecutor implements ExecutorShutdowner for testing
type mockExecutor struct {
	shutdownCalled bool
	shutdownError  error
}

func (m *mockExecutor) Shutdown(ctx context.Context) error {
	m.shutdownCalled = true
	return m.shutdownError
}

func TestNewLifecycle(t *testing.T) {
	cfg := DefaultConfig()
	srv := &server{}
	exec := &mockExecutor{}

	lc := newLifecycle(cfg, srv, exec)

	assert.NotNil(t, lc)
	assert.Equal(t, cfg, lc.cfg)
	assert.Equal(t, srv, lc.server)
	assert.Equal(t, exec, lc.executor)
}

func TestLifecycle_SetupTelemetry(t *testing.T) {
	tests := []struct {
		name        string
		setupConfig func(*Config)
		wantErr     bool
		validate    func(*testing.T, *lifecycle)
	}{
		{
			name: "all telemetry disabled",
			setupConfig: func(cfg *Config) {
				cfg.TracingEnabled = false
				cfg.MetricsEnabled = false
				cfg.SLOEnabled = false
			},
			wantErr: false,
			validate: func(t *testing.T, lc *lifecycle) {
				assert.Nil(t, lc.tracerProvider)
				assert.Nil(t, lc.meterProvider)
				assert.Nil(t, lc.server.sloTracker)
			},
		},
		{
			name: "metrics enabled with prometheus",
			setupConfig: func(cfg *Config) {
				cfg.TracingEnabled = false
				cfg.MetricsEnabled = true
				cfg.MetricsType = telemetry.MetricsExporterPrometheus
				cfg.SLOEnabled = false
			},
			wantErr: false,
			validate: func(t *testing.T, lc *lifecycle) {
				assert.Nil(t, lc.tracerProvider)
				assert.NotNil(t, lc.meterProvider)
				assert.NotNil(t, lc.server.metrics)
				assert.Nil(t, lc.server.sloTracker)
			},
		},
		{
			name: "SLO enabled without metrics",
			setupConfig: func(cfg *Config) {
				cfg.TracingEnabled = false
				cfg.MetricsEnabled = false
				cfg.SLOEnabled = true
			},
			wantErr: false,
			validate: func(t *testing.T, lc *lifecycle) {
				assert.Nil(t, lc.tracerProvider)
				assert.Nil(t, lc.meterProvider)
				assert.NotNil(t, lc.server.sloTracker)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			tt.setupConfig(&cfg)

			srv := &server{}
			exec := &mockExecutor{}
			lc := newLifecycle(cfg, srv, exec)

			ctx := context.Background()
			err := lc.setupTelemetry(ctx)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.validate != nil {
					tt.validate(t, lc)
				}
			}

			// Cleanup
			if lc.meterProvider != nil {
				_ = lc.meterProvider.Shutdown(ctx)
			}
			if lc.tracerProvider != nil {
				_ = lc.tracerProvider.Shutdown(ctx)
			}
		})
	}
}

func TestLifecycle_SetupPprof(t *testing.T) {
	tests := []struct {
		name        string
		pprofEnabled bool
		wantServer  bool
	}{
		{
			name:         "pprof disabled",
			pprofEnabled: false,
			wantServer:   false,
		},
		{
			name:         "pprof enabled",
			pprofEnabled: true,
			wantServer:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			cfg.PprofEnabled = tt.pprofEnabled
			cfg.PprofPort = ":0" // Use random port

			srv := &server{}
			exec := &mockExecutor{}
			lc := newLifecycle(cfg, srv, exec)

			err := lc.setupPprof()
			assert.NoError(t, err)

			if tt.wantServer {
				assert.NotNil(t, lc.pprofServer)
				// Give the server a moment to start
				time.Sleep(10 * time.Millisecond)
				// Cleanup
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer cancel()
				_ = lc.pprofServer.Shutdown(ctx)
			} else {
				assert.Nil(t, lc.pprofServer)
			}
		})
	}
}

func TestLifecycle_SetupRateLimiter(t *testing.T) {
	cfg := DefaultConfig()
	cfg.RateLimit = RateLimitConfig{
		LoginRequestsPerMinute:    5,
		RegisterRequestsPerMinute: 3,
		DEExecutionsPerUser:       10,
		MaxConcurrentDEPerUser:    3,
		MaxRequestsPerSecond:      100,
	}

	srv := &server{}
	exec := &mockExecutor{}
	lc := newLifecycle(cfg, srv, exec)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lc.setupRateLimiter(ctx)

	assert.NotNil(t, lc.rateLimiter)
	assert.NotNil(t, lc.cleanupDone)

	// Cancel context to stop cleanup goroutine
	cancel()

	// Wait for cleanup goroutine to finish
	select {
	case <-lc.cleanupDone:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("Cleanup goroutine did not stop")
	}
}

func TestLifecycle_SetupGRPCServer(t *testing.T) {
	tests := []struct {
		name        string
		setupConfig func(*Config)
		wantErr     bool
	}{
		{
			name: "basic setup without TLS",
			setupConfig: func(cfg *Config) {
				cfg.TLS.Enabled = false
				cfg.SLOEnabled = false
			},
			wantErr: false,
		},
		{
			name: "setup with SLO enabled",
			setupConfig: func(cfg *Config) {
				cfg.TLS.Enabled = false
				cfg.SLOEnabled = true
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			tt.setupConfig(&cfg)

			srv := &server{
				handlers: []handlers.Handler{
					&mockHandler{},
				},
			}
			exec := &mockExecutor{}
			lc := newLifecycle(cfg, srv, exec)

			// Setup SLO tracker if enabled
			if cfg.SLOEnabled {
				ctx := context.Background()
				var err error
				srv.sloTracker, err = slo.NewTracker(ctx, slo.DefaultObjectives())
				require.NoError(t, err)
			}

			// Initialize rate limiter (required for middleware)
			ctx := context.Background()
			lc.setupRateLimiter(ctx)

			err := lc.setupGRPCServer()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, lc.grpcServer)
				assert.NotNil(t, lc.healthSrv)
			}
		})
	}
}

func TestLifecycle_SetupHTTPGateway(t *testing.T) {
	tests := []struct {
		name           string
		tracingEnabled bool
		metricsEnabled bool
		metricsType    telemetry.MetricsExporterType
	}{
		{
			name:           "basic setup",
			tracingEnabled: false,
			metricsEnabled: false,
		},
		{
			name:           "with tracing",
			tracingEnabled: true,
			metricsEnabled: false,
		},
		{
			name:           "with prometheus metrics",
			tracingEnabled: false,
			metricsEnabled: true,
			metricsType:    telemetry.MetricsExporterPrometheus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			cfg.TLS.Enabled = false
			cfg.TracingEnabled = tt.tracingEnabled
			cfg.MetricsEnabled = tt.metricsEnabled
			cfg.MetricsType = tt.metricsType
			cfg.HTTPPort = ":0" // Random port

			srv := &server{
				handlers: []handlers.Handler{
					&mockHandler{},
				},
			}
			exec := &mockExecutor{}
			lc := newLifecycle(cfg, srv, exec)

			ctx := context.Background()
			err := lc.setupHTTPGateway(ctx)

			assert.NoError(t, err)
			assert.NotNil(t, lc.httpServer)
			assert.Equal(t, cfg.HTTPPort, lc.httpServer.Addr)
		})
	}
}

func TestLifecycle_HealthCheckEndpoints(t *testing.T) {
	tests := []struct {
		name           string
		endpoint       string
		dbHealthy      bool
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "health endpoint always returns UP",
			endpoint:       "/health",
			dbHealthy:      true,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"UP"}`,
		},
		{
			name:           "health endpoint returns UP even if DB down",
			endpoint:       "/health",
			dbHealthy:      false,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"UP"}`,
		},
		{
			name:           "readiness returns UP when DB healthy",
			endpoint:       "/readiness",
			dbHealthy:      true,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"UP"}`,
		},
		{
			name:           "readiness returns DOWN when DB unhealthy",
			endpoint:       "/readiness",
			dbHealthy:      false,
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody:   `{"status":"DOWN","reason":"database unavailable"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			cfg.TLS.Enabled = false

			srv := &server{
				handlers: []handlers.Handler{
					&mockHandler{},
				},
			}
			// Mock database health check
			mockStore := &mock.MockStore{}
			if tt.dbHealthy {
				mockStore.HealthCheckFn = func(ctx context.Context) error {
					return nil
				}
			} else {
				mockStore.HealthCheckFn = func(ctx context.Context) error {
					return assert.AnError
				}
			}
			srv.st = mockStore

			exec := &mockExecutor{}
			lc := newLifecycle(cfg, srv, exec)

			ctx := context.Background()
			err := lc.setupHTTPGateway(ctx)
			require.NoError(t, err)

			// Create test request
			req := httptest.NewRequest("GET", tt.endpoint, nil)
			w := httptest.NewRecorder()

			// Serve the request
			lc.httpServer.Handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestLifecycle_Shutdown(t *testing.T) {
	t.Run("graceful shutdown all components", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.TLS.Enabled = false
		cfg.PprofEnabled = true
		cfg.PprofPort = ":0"

		srv := &server{
			handlers: []handlers.Handler{
				&mockHandler{},
			},
		}
		exec := &mockExecutor{}
		lc := newLifecycle(cfg, srv, exec)

		ctx := context.Background()

		// Setup components
		lc.setupRateLimiter(ctx)
		err := lc.setupPprof()
		require.NoError(t, err)
		err = lc.setupGRPCServer()
		require.NoError(t, err)

		// Setup minimal HTTP server
		lc.httpServer = &http.Server{
			Addr:    ":0",
			Handler: http.NewServeMux(),
		}

		// Perform shutdown
		shutdownCtx := context.Background()
		err = lc.shutdown(shutdownCtx)

		assert.NoError(t, err)
		assert.True(t, exec.shutdownCalled, "Executor shutdown should be called")
	})

	t.Run("shutdown with nil components", func(t *testing.T) {
		cfg := DefaultConfig()
		srv := &server{}
		exec := &mockExecutor{}
		lc := newLifecycle(cfg, srv, exec)

		// Don't setup anything, shutdown should handle nil components
		ctx := context.Background()
		err := lc.shutdown(ctx)

		assert.NoError(t, err)
	})

	t.Run("shutdown with executor error", func(t *testing.T) {
		cfg := DefaultConfig()
		srv := &server{}
		exec := &mockExecutor{
			shutdownError: assert.AnError,
		}
		lc := newLifecycle(cfg, srv, exec)

		ctx := context.Background()
		// Shutdown should not propagate executor error
		err := lc.shutdown(ctx)

		assert.NoError(t, err)
		assert.True(t, exec.shutdownCalled)
	})
}

// mockHandler implements handlers.Handler interface for testing
type mockHandler struct{}

func (m *mockHandler) RegisterService(s *grpc.Server) {}

func (m *mockHandler) RegisterHTTPHandler(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return nil
}

func TestLifecycle_SetupRateLimiterCleanup(t *testing.T) {
	cfg := DefaultConfig()
	srv := &server{}
	exec := &mockExecutor{}
	lc := newLifecycle(cfg, srv, exec)

	ctx, cancel := context.WithCancel(context.Background())

	lc.setupRateLimiter(ctx)

	// Verify rate limiter is set
	assert.NotNil(t, lc.rateLimiter)
	assert.NotNil(t, lc.cleanupDone)

	// Test that cleanup goroutine responds to context cancellation
	cancel()

	// Wait for cleanup to complete
	select {
	case <-lc.cleanupDone:
		// Success - goroutine stopped
	case <-time.After(2 * time.Second):
		t.Fatal("Cleanup goroutine did not stop within timeout")
	}
}

func TestLifecycle_SetupGRPCServer_WithoutHandlers(t *testing.T) {
	cfg := DefaultConfig()
	cfg.TLS.Enabled = false

	srv := &server{
		handlers: []handlers.Handler{}, // Empty handlers
	}
	exec := &mockExecutor{}
	lc := newLifecycle(cfg, srv, exec)

	ctx := context.Background()
	lc.setupRateLimiter(ctx)

	err := lc.setupGRPCServer()

	assert.NoError(t, err)
	assert.NotNil(t, lc.grpcServer)
	assert.NotNil(t, lc.healthSrv)
}

func TestLifecycle_CORSConfiguration(t *testing.T) {
	cfg := DefaultConfig()
	cfg.TLS.Enabled = false
	cfg.CORS = middleware.CORSConfig{
		AllowedOrigins:   []string{"https://example.com"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}

	srv := &server{
		handlers: []handlers.Handler{
			&mockHandler{},
		},
	}
	exec := &mockExecutor{}
	lc := newLifecycle(cfg, srv, exec)

	ctx := context.Background()
	err := lc.setupHTTPGateway(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, lc.httpServer)
	assert.NotNil(t, lc.httpServer.Handler)

	// Test CORS headers with OPTIONS request
	req := httptest.NewRequest("OPTIONS", "/health", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	lc.httpServer.Handler.ServeHTTP(w, req)

	assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
}
