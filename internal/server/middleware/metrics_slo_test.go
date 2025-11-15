package middleware

import (
	"context"
	"testing"
	"time"

	"github.com/nicholaspcr/GoDE/internal/slo"
	"github.com/nicholaspcr/GoDE/internal/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnaryMetricsAndSLOMiddleware_Success(t *testing.T) {
	ctx := context.Background()

	// Create metrics
	metrics, err := telemetry.InitMetrics(ctx, "test-service")
	require.NoError(t, err)

	// Create SLO tracker
	objectives := []slo.Objective{
		{
			Name:        "test_availability",
			Description: "Test availability",
			Target:      99.0,
			Window:      1 * time.Hour,
		},
	}
	tracker, err := slo.NewTracker(ctx, objectives)
	require.NoError(t, err)

	// Create middleware
	middleware := UnaryMetricsAndSLOMiddleware(metrics, tracker)

	// Handler that succeeds
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/api.v1.TestService/TestMethod",
	}

	resp, err := middleware(ctx, nil, info, handler)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp)

	// Verify SLO tracking
	compliance := tracker.GetCompliance("test_availability")
	assert.Equal(t, 100.0, compliance)
}

func TestUnaryMetricsAndSLOMiddleware_Error(t *testing.T) {
	ctx := context.Background()

	metrics, err := telemetry.InitMetrics(ctx, "test-service")
	require.NoError(t, err)

	objectives := []slo.Objective{
		{
			Name:        "test_availability",
			Description: "Test availability",
			Target:      99.0,
			Window:      1 * time.Hour,
		},
	}
	tracker, err := slo.NewTracker(ctx, objectives)
	require.NoError(t, err)

	middleware := UnaryMetricsAndSLOMiddleware(metrics, tracker)

	// Handler that fails
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, status.Error(codes.Internal, "internal error")
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/api.v1.TestService/TestMethod",
	}

	resp, err := middleware(ctx, nil, info, handler)

	assert.Error(t, err)
	assert.Nil(t, resp)

	// Verify SLO tracking recorded failure
	compliance := tracker.GetCompliance("test_availability")
	assert.Equal(t, 0.0, compliance) // 1 request, 0 success = 0%
}

func TestUnaryMetricsAndSLOMiddleware_NilMetrics(t *testing.T) {
	ctx := context.Background()

	objectives := []slo.Objective{
		{
			Name:        "test_availability",
			Description: "Test availability",
			Target:      99.0,
			Window:      1 * time.Hour,
		},
	}
	tracker, err := slo.NewTracker(ctx, objectives)
	require.NoError(t, err)

	// Nil metrics should not cause panic
	middleware := UnaryMetricsAndSLOMiddleware(nil, tracker)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/api.v1.TestService/TestMethod",
	}

	resp, err := middleware(ctx, nil, info, handler)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp)
}

func TestUnaryMetricsAndSLOMiddleware_NilSLOTracker(t *testing.T) {
	ctx := context.Background()

	metrics, err := telemetry.InitMetrics(ctx, "test-service")
	require.NoError(t, err)

	// Nil SLO tracker should not cause panic
	middleware := UnaryMetricsAndSLOMiddleware(metrics, nil)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/api.v1.TestService/TestMethod",
	}

	resp, err := middleware(ctx, nil, info, handler)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp)
}

func TestUnaryMetricsAndSLOMiddleware_MultipleRequests(t *testing.T) {
	ctx := context.Background()

	metrics, err := telemetry.InitMetrics(ctx, "test-service")
	require.NoError(t, err)

	objectives := []slo.Objective{
		{
			Name:        "test_availability",
			Description: "Test availability",
			Target:      99.0,
			Window:      1 * time.Hour,
		},
	}
	tracker, err := slo.NewTracker(ctx, objectives)
	require.NoError(t, err)

	middleware := UnaryMetricsAndSLOMiddleware(metrics, tracker)

	info := &grpc.UnaryServerInfo{
		FullMethod: "/api.v1.TestService/TestMethod",
	}

	// Simulate 10 successful requests
	successHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "success", nil
	}

	for i := 0; i < 10; i++ {
		_, err := middleware(ctx, nil, info, successHandler)
		assert.NoError(t, err)
	}

	// Simulate 1 failed request
	failHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, status.Error(codes.Internal, "error")
	}

	_, err = middleware(ctx, nil, info, failHandler)
	assert.Error(t, err)

	// Verify compliance: 10 success out of 11 total = 90.9%
	compliance := tracker.GetCompliance("test_availability")
	assert.InDelta(t, 90.9, compliance, 0.1)
}

func TestUnaryMetricsAndSLOMiddleware_DifferentStatusCodes(t *testing.T) {
	ctx := context.Background()

	metrics, err := telemetry.InitMetrics(ctx, "test-service")
	require.NoError(t, err)

	objectives := []slo.Objective{
		{
			Name:        "test_availability",
			Description: "Test availability",
			Target:      99.0,
			Window:      1 * time.Hour,
		},
	}
	tracker, err := slo.NewTracker(ctx, objectives)
	require.NoError(t, err)

	middleware := UnaryMetricsAndSLOMiddleware(metrics, tracker)

	info := &grpc.UnaryServerInfo{
		FullMethod: "/api.v1.TestService/TestMethod",
	}

	tests := []struct {
		name          string
		handler       grpc.UnaryHandler
		expectSuccess bool
	}{
		{
			name: "OK",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return "ok", nil
			},
			expectSuccess: true,
		},
		{
			name: "NotFound",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, status.Error(codes.NotFound, "not found")
			},
			expectSuccess: false,
		},
		{
			name: "InvalidArgument",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, status.Error(codes.InvalidArgument, "invalid")
			},
			expectSuccess: false,
		},
		{
			name: "Internal",
			handler: func(ctx context.Context, req interface{}) (interface{}, error) {
				return nil, status.Error(codes.Internal, "internal")
			},
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := middleware(ctx, nil, info, tt.handler)
			if tt.expectSuccess {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestUnaryMetricsAndSLOMiddleware_RecordsDuration(t *testing.T) {
	ctx := context.Background()

	metrics, err := telemetry.InitMetrics(ctx, "test-service")
	require.NoError(t, err)

	objectives := []slo.Objective{
		{
			Name:        "test_availability",
			Description: "Test availability",
			Target:      99.0,
			Window:      1 * time.Hour,
		},
	}
	tracker, err := slo.NewTracker(ctx, objectives)
	require.NoError(t, err)

	middleware := UnaryMetricsAndSLOMiddleware(metrics, tracker)

	// Handler with deliberate delay
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		time.Sleep(10 * time.Millisecond)
		return "success", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/api.v1.TestService/TestMethod",
	}

	start := time.Now()
	resp, err := middleware(ctx, nil, info, handler)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, "success", resp)
	assert.GreaterOrEqual(t, duration, 10*time.Millisecond)
}
