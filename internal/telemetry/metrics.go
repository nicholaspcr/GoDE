package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// MetricsExporterType defines the type of metrics exporter to use.
type MetricsExporterType string

const (
	// MetricsExporterStdout exports metrics to stdout (for development).
	MetricsExporterStdout MetricsExporterType = "stdout"
	// MetricsExporterPrometheus exports metrics in Prometheus format.
	MetricsExporterPrometheus MetricsExporterType = "prometheus"
)

// Metrics holds all application metrics.
type Metrics struct {
	// API metrics
	APIRequestsTotal      metric.Int64Counter
	APIRequestDuration    metric.Float64Histogram
	APIRequestsInFlight   metric.Int64UpDownCounter
	APIErrorsTotal        metric.Int64Counter

	// DE execution metrics
	DEExecutionsTotal     metric.Int64Counter
	DEExecutionDuration   metric.Float64Histogram
	DEExecutionsInFlight  metric.Int64UpDownCounter
	DEGenerationsTotal    metric.Int64Counter
	ParetoSetSize         metric.Int64Histogram

	// Auth metrics
	AuthAttemptsTotal     metric.Int64Counter
	AuthSuccessTotal      metric.Int64Counter
	AuthFailuresTotal     metric.Int64Counter

	// Rate limiting metrics
	RateLimitExceeded     metric.Int64Counter

	// Panic metrics
	PanicsTotal           metric.Int64Counter
}

// NewMeterProvider creates a new meter provider with the specified exporter type.
func NewMeterProvider(appName string, exporterType MetricsExporterType) (*sdkmetric.MeterProvider, error) {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(appName),
	)

	var mp *sdkmetric.MeterProvider

	switch exporterType {
	case MetricsExporterStdout:
		exporter, err := stdoutmetric.New()
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout metrics exporter: %w", err)
		}
		mp = sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(res),
			sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		)

	case MetricsExporterPrometheus:
		exporter, err := prometheus.New()
		if err != nil {
			return nil, fmt.Errorf("failed to create prometheus metrics exporter: %w", err)
		}
		mp = sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(res),
			sdkmetric.WithReader(exporter),
		)

	default:
		return nil, fmt.Errorf("unknown metrics exporter type: %s", exporterType)
	}

	otel.SetMeterProvider(mp)
	return mp, nil
}

// InitMetrics initializes all application metrics.
func InitMetrics(ctx context.Context, appName string) (*Metrics, error) {
	meter := otel.Meter(appName)

	m := &Metrics{}
	var err error

	// API metrics
	m.APIRequestsTotal, err = meter.Int64Counter(
		"api_requests_total",
		metric.WithDescription("Total number of API requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create api_requests_total counter: %w", err)
	}

	m.APIRequestDuration, err = meter.Float64Histogram(
		"api_request_duration_seconds",
		metric.WithDescription("API request duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create api_request_duration_seconds histogram: %w", err)
	}

	m.APIRequestsInFlight, err = meter.Int64UpDownCounter(
		"api_requests_in_flight",
		metric.WithDescription("Number of API requests currently being processed"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create api_requests_in_flight gauge: %w", err)
	}

	m.APIErrorsTotal, err = meter.Int64Counter(
		"api_errors_total",
		metric.WithDescription("Total number of API errors"),
		metric.WithUnit("{error}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create api_errors_total counter: %w", err)
	}

	// DE execution metrics
	m.DEExecutionsTotal, err = meter.Int64Counter(
		"de_executions_total",
		metric.WithDescription("Total number of DE executions"),
		metric.WithUnit("{execution}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create de_executions_total counter: %w", err)
	}

	m.DEExecutionDuration, err = meter.Float64Histogram(
		"de_execution_duration_seconds",
		metric.WithDescription("DE execution duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create de_execution_duration_seconds histogram: %w", err)
	}

	m.DEExecutionsInFlight, err = meter.Int64UpDownCounter(
		"de_executions_in_flight",
		metric.WithDescription("Number of DE executions currently running"),
		metric.WithUnit("{execution}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create de_executions_in_flight gauge: %w", err)
	}

	m.DEGenerationsTotal, err = meter.Int64Counter(
		"de_generations_total",
		metric.WithDescription("Total number of DE generations executed"),
		metric.WithUnit("{generation}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create de_generations_total counter: %w", err)
	}

	m.ParetoSetSize, err = meter.Int64Histogram(
		"pareto_set_size",
		metric.WithDescription("Size of Pareto set returned from DE execution"),
		metric.WithUnit("{solution}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create pareto_set_size histogram: %w", err)
	}

	// Auth metrics
	m.AuthAttemptsTotal, err = meter.Int64Counter(
		"auth_attempts_total",
		metric.WithDescription("Total number of authentication attempts"),
		metric.WithUnit("{attempt}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth_attempts_total counter: %w", err)
	}

	m.AuthSuccessTotal, err = meter.Int64Counter(
		"auth_success_total",
		metric.WithDescription("Total number of successful authentications"),
		metric.WithUnit("{success}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth_success_total counter: %w", err)
	}

	m.AuthFailuresTotal, err = meter.Int64Counter(
		"auth_failures_total",
		metric.WithDescription("Total number of failed authentications"),
		metric.WithUnit("{failure}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth_failures_total counter: %w", err)
	}

	// Rate limiting metrics
	m.RateLimitExceeded, err = meter.Int64Counter(
		"rate_limit_exceeded_total",
		metric.WithDescription("Total number of rate limit exceeded events"),
		metric.WithUnit("{event}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create rate_limit_exceeded_total counter: %w", err)
	}

	// Panic metrics
	m.PanicsTotal, err = meter.Int64Counter(
		"panics_total",
		metric.WithDescription("Total number of recovered panics"),
		metric.WithUnit("{panic}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create panics_total counter: %w", err)
	}

	return m, nil
}
