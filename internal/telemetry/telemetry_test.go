package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTracerProvider_Stdout(t *testing.T) {
	cfg := TracingConfig{
		ExporterType: TracingExporterStdout,
		SampleRatio:  1.0,
	}

	tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
	require.NoError(t, err)
	require.NotNil(t, tp)

	err = tp.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestNewTracerProvider_OTLP(t *testing.T) {
	cfg := TracingConfig{
		ExporterType: TracingExporterOTLP,
		OTLPEndpoint: "localhost:4317",
		SampleRatio:  1.0,
	}

	tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
	require.NoError(t, err)
	require.NotNil(t, tp)

	err = tp.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestNewTracerProvider_SampleRatio(t *testing.T) {
	tests := []struct {
		name        string
		sampleRatio float64
	}{
		{"always sample", 1.0},
		{"never sample", 0.0},
		{"half sample", 0.5},
		{"tenth sample", 0.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := TracingConfig{
				ExporterType: TracingExporterStdout,
				SampleRatio:  tt.sampleRatio,
			}

			tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
			require.NoError(t, err)
			require.NotNil(t, tp)

			err = tp.Shutdown(context.Background())
			assert.NoError(t, err)
		})
	}
}

func TestNewTracerProvider_InvalidExporter(t *testing.T) {
	cfg := TracingConfig{
		ExporterType: "invalid",
		SampleRatio:  1.0,
	}

	tp, err := NewTracerProvider(context.Background(), "test-service", cfg)
	assert.Error(t, err)
	assert.Nil(t, tp)
	assert.Contains(t, err.Error(), "unknown tracing exporter type")
}

func TestNewTracerProvider_ServiceName(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
	}{
		{"simple name", "myservice"},
		{"hyphenated name", "my-service"},
		{"underscored name", "my_service"},
		{"with numbers", "service-v1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := TracingConfig{
				ExporterType: TracingExporterStdout,
				SampleRatio:  1.0,
			}

			tp, err := NewTracerProvider(context.Background(), tt.serviceName, cfg)
			require.NoError(t, err)
			require.NotNil(t, tp)

			err = tp.Shutdown(context.Background())
			assert.NoError(t, err)
		})
	}
}

func TestNewMeterProvider_Stdout(t *testing.T) {
	mp, err := NewMeterProvider("test-service", MetricsExporterStdout)
	require.NoError(t, err)
	require.NotNil(t, mp)

	err = mp.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestNewMeterProvider_Prometheus(t *testing.T) {
	mp, err := NewMeterProvider("test-service", MetricsExporterPrometheus)
	require.NoError(t, err)
	require.NotNil(t, mp)

	err = mp.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestNewMeterProvider_InvalidExporter(t *testing.T) {
	mp, err := NewMeterProvider("test-service", "invalid")
	assert.Error(t, err)
	assert.Nil(t, mp)
	assert.Contains(t, err.Error(), "unknown metrics exporter type")
}

func TestNewMeterProvider_ServiceName(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
	}{
		{"simple name", "myservice"},
		{"hyphenated name", "my-service"},
		{"underscored name", "my_service"},
		{"with numbers", "service-v1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp, err := NewMeterProvider(tt.serviceName, MetricsExporterStdout)
			require.NoError(t, err)
			require.NotNil(t, mp)

			err = mp.Shutdown(context.Background())
			assert.NoError(t, err)
		})
	}
}

func TestInitMetrics(t *testing.T) {
	// Setup meter provider first
	mp, err := NewMeterProvider("test-service", MetricsExporterStdout)
	require.NoError(t, err)
	require.NotNil(t, mp)
	defer func() {
		_ = mp.Shutdown(context.Background())
	}()

	m, err := InitMetrics(context.Background(), "test-service")
	require.NoError(t, err)
	require.NotNil(t, m)

	// Verify all API metrics are initialized
	assert.NotNil(t, m.APIRequestsTotal)
	assert.NotNil(t, m.APIRequestDuration)
	assert.NotNil(t, m.APIRequestsInFlight)
	assert.NotNil(t, m.APIErrorsTotal)

	// Verify all DE execution metrics are initialized
	assert.NotNil(t, m.DEExecutionsTotal)
	assert.NotNil(t, m.DEExecutionDuration)
	assert.NotNil(t, m.DEExecutionsInFlight)
	assert.NotNil(t, m.DEGenerationsTotal)
	assert.NotNil(t, m.ParetoSetSize)

	// Verify executor worker pool metrics are initialized
	assert.NotNil(t, m.ExecutorWorkersActive)
	assert.NotNil(t, m.ExecutorWorkersTotal)
	assert.NotNil(t, m.ExecutorQueueWaitDuration)
	assert.NotNil(t, m.ExecutorUtilizationPercent)

	// Verify auth metrics are initialized
	assert.NotNil(t, m.AuthAttemptsTotal)
	assert.NotNil(t, m.AuthSuccessTotal)
	assert.NotNil(t, m.AuthFailuresTotal)

	// Verify rate limiting metrics are initialized
	assert.NotNil(t, m.RateLimitExceeded)

	// Verify panic metrics are initialized
	assert.NotNil(t, m.PanicsTotal)
}

func TestInitMetrics_WithPrometheus(t *testing.T) {
	// Setup Prometheus meter provider
	mp, err := NewMeterProvider("test-service", MetricsExporterPrometheus)
	require.NoError(t, err)
	require.NotNil(t, mp)
	defer func() {
		_ = mp.Shutdown(context.Background())
	}()

	m, err := InitMetrics(context.Background(), "test-service")
	require.NoError(t, err)
	require.NotNil(t, m)

	// Verify metrics can be recorded (basic smoke test)
	ctx := context.Background()
	m.APIRequestsTotal.Add(ctx, 1)
	m.DEExecutionsTotal.Add(ctx, 1)
	m.AuthAttemptsTotal.Add(ctx, 1)
}
