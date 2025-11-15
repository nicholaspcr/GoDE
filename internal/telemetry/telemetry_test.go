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
