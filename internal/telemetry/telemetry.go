package telemetry

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc/credentials"
)

// TracingExporterType defines the type of tracing exporter to use.
type TracingExporterType string

const (
	// TracingExporterStdout exports traces to stdout (for development).
	TracingExporterStdout TracingExporterType = "stdout"
	// TracingExporterOTLP exports traces via OTLP (for Jaeger, Tempo, etc.).
	TracingExporterOTLP TracingExporterType = "otlp"
	// TracingExporterFile exports traces to a file.
	TracingExporterFile TracingExporterType = "file"
	// TracingExporterNone disables trace export (traces still recorded but not exported).
	TracingExporterNone TracingExporterType = "none"
)

// TLSConfig holds TLS configuration for secure connections.
type TLSConfig struct {
	Enabled    bool   // Enable TLS
	CertFile   string // Path to client certificate file (PEM)
	KeyFile    string // Path to client key file (PEM)
	CAFile     string // Path to CA certificate file (PEM) for server verification
	SkipVerify bool   // Skip server certificate verification (not recommended for production)
}

// TracingConfig holds configuration for tracing.
type TracingConfig struct {
	ExporterType TracingExporterType
	OTLPEndpoint string  // e.g., "localhost:4317" for gRPC
	SampleRatio  float64 // 0.0 to 1.0, where 1.0 = sample everything
	TLS          TLSConfig
	FilePath     string // Path to trace output file (for file exporter)
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
		opts := []otlptracegrpc.Option{}

		if cfg.OTLPEndpoint != "" {
			opts = append(opts, otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint))
		}

		// Configure TLS or insecure connection
		if cfg.TLS.Enabled {
			tlsCreds, tlsErr := buildTLSCredentials(cfg.TLS)
			if tlsErr != nil {
				return nil, fmt.Errorf("failed to build TLS credentials: %w", tlsErr)
			}
			opts = append(opts, otlptracegrpc.WithTLSCredentials(tlsCreds))
		} else {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}

		exporter, err = otlptracegrpc.New(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
		}

	case TracingExporterFile:
		if cfg.FilePath == "" {
			cfg.FilePath = "traces.json"
		}
		file, fileErr := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if fileErr != nil {
			return nil, fmt.Errorf("failed to open trace file %s: %w", cfg.FilePath, fileErr)
		}
		exporter, err = stdouttrace.New(stdouttrace.WithWriter(file), stdouttrace.WithPrettyPrint())
		if err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to create file trace exporter: %w", err)
		}

	case TracingExporterNone:
		// No exporter - create a minimal tracer provider without batcher
		sampler := trace.ParentBased(trace.TraceIDRatioBased(cfg.SampleRatio))
		if cfg.SampleRatio >= 1.0 {
			sampler = trace.AlwaysSample()
		} else if cfg.SampleRatio <= 0.0 {
			sampler = trace.NeverSample()
		}

		tp := trace.NewTracerProvider(
			trace.WithResource(res),
			trace.WithSampler(sampler),
		)
		otel.SetTracerProvider(tp)
		return tp, nil

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

// buildTLSCredentials creates gRPC transport credentials from TLS configuration.
func buildTLSCredentials(cfg TLSConfig) (credentials.TransportCredentials, error) {
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: cfg.SkipVerify, //nolint:gosec // Configurable for testing environments
	}

	// Load client certificate if provided
	if cfg.CertFile != "" && cfg.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	// Load CA certificate for server verification if provided
	if cfg.CAFile != "" {
		caCert, err := os.ReadFile(cfg.CAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig.RootCAs = caCertPool
	}

	return credentials.NewTLS(tlsConfig), nil
}
