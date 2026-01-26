// Package telemetry provides observability through metrics and tracing.
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
	APIRequestsTotal    metric.Int64Counter
	APIRequestDuration  metric.Float64Histogram
	APIRequestsInFlight metric.Int64UpDownCounter
	APIErrorsTotal      metric.Int64Counter

	// DE execution metrics
	DEExecutionsTotal    metric.Int64Counter
	DEExecutionDuration  metric.Float64Histogram
	DEExecutionsInFlight metric.Int64UpDownCounter
	DEGenerationsTotal   metric.Int64Counter
	ParetoSetSize        metric.Int64Histogram

	// DE algorithm-specific metrics
	DEObjectiveEvaluations   metric.Int64Counter   // Total objective function evaluations
	DEMutationsTotal         metric.Int64Counter   // Total mutation operations
	DECrossoverTotal         metric.Int64Counter   // Total crossover operations
	DEPopulationDiversity    metric.Float64Histogram // Measure of population diversity
	DEConvergenceRate        metric.Float64Histogram // Rate of improvement per generation
	DENonDominatedCount      metric.Int64Histogram   // Count of non-dominated solutions
	DERankZeroSize           metric.Int64Histogram   // Size of rank 0 front
	DECrowdingDistanceAvg    metric.Float64Histogram // Average crowding distance
	DEVariantPerformance     metric.Float64Histogram // Performance by variant type
	DEProblemComplexity      metric.Int64Histogram   // Dimensions × objectives

	// Executor worker pool metrics
	ExecutorWorkersActive      metric.Int64UpDownCounter
	ExecutorWorkersTotal       metric.Int64UpDownCounter
	ExecutorQueueWaitDuration  metric.Float64Histogram
	ExecutorUtilizationPercent metric.Float64Histogram

	// Auth metrics
	AuthAttemptsTotal metric.Int64Counter
	AuthSuccessTotal  metric.Int64Counter
	AuthFailuresTotal metric.Int64Counter

	// Rate limiting metrics
	RateLimitExceeded metric.Int64Counter

	// Panic metrics
	PanicsTotal metric.Int64Counter
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

	// DE algorithm-specific metrics
	m.DEObjectiveEvaluations, err = meter.Int64Counter(
		"de_objective_evaluations_total",
		metric.WithDescription("Total number of objective function evaluations"),
		metric.WithUnit("{evaluation}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create de_objective_evaluations_total counter: %w", err)
	}

	m.DEMutationsTotal, err = meter.Int64Counter(
		"de_mutations_total",
		metric.WithDescription("Total number of mutation operations performed"),
		metric.WithUnit("{mutation}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create de_mutations_total counter: %w", err)
	}

	m.DECrossoverTotal, err = meter.Int64Counter(
		"de_crossover_total",
		metric.WithDescription("Total number of crossover operations performed"),
		metric.WithUnit("{crossover}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create de_crossover_total counter: %w", err)
	}

	m.DEPopulationDiversity, err = meter.Float64Histogram(
		"de_population_diversity",
		metric.WithDescription("Measure of population diversity (standard deviation of objectives)"),
		metric.WithUnit("{diversity}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create de_population_diversity histogram: %w", err)
	}

	m.DEConvergenceRate, err = meter.Float64Histogram(
		"de_convergence_rate",
		metric.WithDescription("Rate of improvement per generation"),
		metric.WithUnit("{improvement}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create de_convergence_rate histogram: %w", err)
	}

	m.DENonDominatedCount, err = meter.Int64Histogram(
		"de_non_dominated_count",
		metric.WithDescription("Number of non-dominated solutions in current generation"),
		metric.WithUnit("{solution}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create de_non_dominated_count histogram: %w", err)
	}

	m.DERankZeroSize, err = meter.Int64Histogram(
		"de_rank_zero_size",
		metric.WithDescription("Size of rank 0 (Pareto front) in current generation"),
		metric.WithUnit("{solution}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create de_rank_zero_size histogram: %w", err)
	}

	m.DECrowdingDistanceAvg, err = meter.Float64Histogram(
		"de_crowding_distance_avg",
		metric.WithDescription("Average crowding distance of solutions in Pareto front"),
		metric.WithUnit("{distance}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create de_crowding_distance_avg histogram: %w", err)
	}

	m.DEVariantPerformance, err = meter.Float64Histogram(
		"de_variant_performance",
		metric.WithDescription("Performance metrics by DE variant (rand, best, current-to-best, pbest)"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create de_variant_performance histogram: %w", err)
	}

	m.DEProblemComplexity, err = meter.Int64Histogram(
		"de_problem_complexity",
		metric.WithDescription("Problem complexity (dimensions × objectives)"),
		metric.WithUnit("{complexity}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create de_problem_complexity histogram: %w", err)
	}

	// Executor worker pool metrics
	m.ExecutorWorkersActive, err = meter.Int64UpDownCounter(
		"executor_workers_active",
		metric.WithDescription("Number of executor workers currently processing tasks"),
		metric.WithUnit("{worker}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create executor_workers_active gauge: %w", err)
	}

	m.ExecutorWorkersTotal, err = meter.Int64UpDownCounter(
		"executor_workers_total",
		metric.WithDescription("Total number of executor workers available"),
		metric.WithUnit("{worker}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create executor_workers_total gauge: %w", err)
	}

	m.ExecutorQueueWaitDuration, err = meter.Float64Histogram(
		"executor_queue_wait_duration_seconds",
		metric.WithDescription("Time spent waiting in executor queue before worker assignment"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create executor_queue_wait_duration_seconds histogram: %w", err)
	}

	m.ExecutorUtilizationPercent, err = meter.Float64Histogram(
		"executor_utilization_percent",
		metric.WithDescription("Executor worker pool utilization percentage"),
		metric.WithUnit("%"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create executor_utilization_percent histogram: %w", err)
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
