package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// TracingExporterType defines the type of tracing exporter to use.
type TracingExporterType string

const (
	// TracingExporterStdout exports traces to stdout (for development).
	TracingExporterStdout TracingExporterType = "stdout"
	// TracingExporterOTLP exports traces via OTLP (for Jaeger, Tempo, etc.).
	TracingExporterOTLP TracingExporterType = "otlp"
)

// TracingConfig holds configuration for tracing.
type TracingConfig struct {
	ExporterType TracingExporterType
	OTLPEndpoint string // e.g., "localhost:4317" for gRPC
	SampleRatio  float64 // 0.0 to 1.0, where 1.0 = sample everything
}

// NewTracerProvider creates a new tracer provider with the specified configuration.
func NewTracerProvider(ctx context.Context, appName string, cfg TracingConfig) (*trace.TracerProvider, error) {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(appName),
	)

	var exporter trace.SpanExporter
	var err error

	switch cfg.ExporterType {
	case TracingExporterStdout:
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout trace exporter: %w", err)
		}

	case TracingExporterOTLP:
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithInsecure(), // TODO: Add TLS support in production
		}
		if cfg.OTLPEndpoint != "" {
			opts = append(opts, otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint))
		}
		exporter, err = otlptracegrpc.New(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
		}

	default:
		return nil, fmt.Errorf("unknown tracing exporter type: %s", cfg.ExporterType)
	}

	// Configure sampling
	sampler := trace.ParentBased(trace.TraceIDRatioBased(cfg.SampleRatio))
	if cfg.SampleRatio >= 1.0 {
		sampler = trace.AlwaysSample()
	} else if cfg.SampleRatio <= 0.0 {
		sampler = trace.NeverSample()
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
		trace.WithSampler(sampler),
	)
	otel.SetTracerProvider(tp)

	return tp, nil
}
