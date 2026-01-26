package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/nicholaspcr/GoDE/internal/store/mock"
)

func TestSetupHealthCheck(t *testing.T) {
	t.Run("registers health server with all services", func(t *testing.T) {
		srv := &server{}
		grpcSrv := grpc.NewServer()
		defer grpcSrv.Stop()

		healthSrv := srv.setupHealthCheck(grpcSrv)

		require.NotNil(t, healthSrv)

		// Verify all services are registered as SERVING
		ctx := context.Background()

		// Check overall server health
		resp, err := healthSrv.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: ""})
		require.NoError(t, err)
		assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.Status)

		// Check individual services
		services := []string{
			"api.v1.AuthService",
			"api.v1.UserService",
			"api.v1.DifferentialEvolutionService",
			"api.v1.ParetoSetService",
		}

		for _, service := range services {
			resp, err := healthSrv.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: service})
			require.NoError(t, err, "service %s should be registered", service)
			assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.Status,
				"service %s should be SERVING", service)
		}
	})

	t.Run("health server can be queried via Watch", func(t *testing.T) {
		srv := &server{}
		grpcSrv := grpc.NewServer()
		defer grpcSrv.Stop()

		healthSrv := srv.setupHealthCheck(grpcSrv)
		require.NotNil(t, healthSrv)

		// Create a Watch stream
		ctx := t.Context()

		stream := &mockHealthWatchServer{ctx: ctx, updates: make(chan *grpc_health_v1.HealthCheckResponse, 1)}

		// Start watching in background
		go func() {
			_ = healthSrv.Watch(&grpc_health_v1.HealthCheckRequest{Service: ""}, stream)
		}()

		// Should receive initial status
		select {
		case resp := <-stream.updates:
			assert.Equal(t, grpc_health_v1.HealthCheckResponse_SERVING, resp.Status)
		case <-ctx.Done():
			t.Fatal("did not receive health update")
		}
	})
}

func TestCheckDatabaseHealth(t *testing.T) {
	tests := []struct {
		name       string
		setupStore func() *mock.MockStore
		nilStore   bool
		wantResult bool
	}{
		{
			name:       "returns false when store is nil",
			nilStore:   true,
			wantResult: false,
		},
		{
			name: "returns false when store health check fails",
			setupStore: func() *mock.MockStore {
				m := &mock.MockStore{}
				m.HealthCheckFn = func(ctx context.Context) error {
					return assert.AnError
				}
				return m
			},
			wantResult: false,
		},
		{
			name: "returns true when store health check succeeds",
			setupStore: func() *mock.MockStore {
				m := &mock.MockStore{}
				m.HealthCheckFn = func(ctx context.Context) error {
					return nil
				}
				return m
			},
			wantResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &server{}

			if !tt.nilStore && tt.setupStore != nil {
				srv.st = tt.setupStore()
			}

			ctx := context.Background()
			result := srv.checkDatabaseHealth(ctx)

			assert.Equal(t, tt.wantResult, result)
		})
	}
}

func TestCheckDatabaseHealth_ContextCancellation(t *testing.T) {
	t.Run("respects context cancellation", func(t *testing.T) {
		srv := &server{}
		mockStore := &mock.MockStore{}

		// Mock that waits for context cancellation
		mockStore.HealthCheckFn = func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		}
		srv.st = mockStore

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		result := srv.checkDatabaseHealth(ctx)
		assert.False(t, result, "should return false when context is cancelled")
	})
}

// mockHealthWatchServer implements grpc_health_v1.Health_WatchServer for testing
type mockHealthWatchServer struct {
	grpc.ServerStream
	ctx     context.Context
	updates chan *grpc_health_v1.HealthCheckResponse
}

func (m *mockHealthWatchServer) Send(resp *grpc_health_v1.HealthCheckResponse) error {
	select {
	case m.updates <- resp:
		return nil
	case <-m.ctx.Done():
		return m.ctx.Err()
	}
}

func (m *mockHealthWatchServer) Context() context.Context {
	return m.ctx
}
