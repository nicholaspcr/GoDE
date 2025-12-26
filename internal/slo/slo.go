// Package slo provides Service Level Objective (SLO) tracking and monitoring.
// SLOs define the target level of service quality and measure actual performance against those targets.
package slo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Objective defines a specific SLO target.
type Objective struct {
	Name        string  // Name of the objective (e.g., "availability", "latency_p95")
	Description string  // Human-readable description
	Target      float64 // Target percentage (0.0-100.0), e.g., 99.9 for 99.9%
	Window      time.Duration // Time window for evaluation (e.g., 24h, 7d)
}

// SLOType represents the type of SLO being tracked.
type SLOType string

const (
	// SLOTypeAvailability tracks the percentage of successful requests.
	SLOTypeAvailability SLOType = "availability"
	// SLOTypeLatency tracks the percentage of requests below a latency threshold.
	SLOTypeLatency SLOType = "latency"
	// SLOTypeErrorRate tracks the percentage of requests that error.
	SLOTypeErrorRate SLOType = "error_rate"
)

// Tracker tracks SLO compliance and records measurements.
type Tracker struct {
	objectives map[string]*Objective
	mu         sync.RWMutex

	// Metrics
	sloCompliance    metric.Float64Gauge
	sloViolations    metric.Int64Counter
	totalRequests    metric.Int64Counter
	successRequests  metric.Int64Counter
	errorRequests    metric.Int64Counter
	latencyHistogram metric.Float64Histogram

	// Time-windowed tracking
	windows map[string]*window
}

// window tracks metrics within a specific time window.
type window struct {
	startTime time.Time
	duration  time.Duration
	total     int64
	success   int64
	errors    int64
	mu        sync.Mutex
}

// NewTracker creates a new SLO tracker with the given objectives.
func NewTracker(ctx context.Context, objectives []Objective) (*Tracker, error) {
	meter := otel.Meter("slo")

	t := &Tracker{
		objectives: make(map[string]*Objective),
		windows:    make(map[string]*window),
	}

	// Initialize objectives
	for i := range objectives {
		t.objectives[objectives[i].Name] = &objectives[i]
		t.windows[objectives[i].Name] = &window{
			startTime: time.Now(),
			duration:  objectives[i].Window,
		}
	}

	// Initialize metrics
	var err error

	t.sloCompliance, err = meter.Float64Gauge(
		"slo_compliance_percentage",
		metric.WithDescription("Current SLO compliance percentage"),
		metric.WithUnit("%"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create slo_compliance metric: %w", err)
	}

	t.sloViolations, err = meter.Int64Counter(
		"slo_violations_total",
		metric.WithDescription("Total number of SLO violations"),
		metric.WithUnit("{violation}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create slo_violations metric: %w", err)
	}

	t.totalRequests, err = meter.Int64Counter(
		"slo_requests_total",
		metric.WithDescription("Total number of requests tracked for SLO"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create slo_requests_total metric: %w", err)
	}

	t.successRequests, err = meter.Int64Counter(
		"slo_requests_success",
		metric.WithDescription("Total number of successful requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create slo_requests_success metric: %w", err)
	}

	t.errorRequests, err = meter.Int64Counter(
		"slo_requests_error",
		metric.WithDescription("Total number of failed requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create slo_requests_error metric: %w", err)
	}

	t.latencyHistogram, err = meter.Float64Histogram(
		"slo_request_duration_seconds",
		metric.WithDescription("Request duration for SLO tracking"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create slo_request_duration metric: %w", err)
	}

	return t, nil
}

// RecordRequest records a request for SLO tracking.
func (t *Tracker) RecordRequest(ctx context.Context, service string, success bool, durationSeconds float64) {
	attrs := []attribute.KeyValue{
		attribute.String("service", service),
	}

	// Record total requests
	t.totalRequests.Add(ctx, 1, metric.WithAttributes(attrs...))

	// Record success/error
	if success {
		t.successRequests.Add(ctx, 1, metric.WithAttributes(attrs...))
	} else {
		t.errorRequests.Add(ctx, 1, metric.WithAttributes(attrs...))
	}

	// Record latency
	t.latencyHistogram.Record(ctx, durationSeconds, metric.WithAttributes(attrs...))

	// Update time windows
	t.mu.Lock()
	defer t.mu.Unlock()

	for name, w := range t.windows {
		w.mu.Lock()
		// Reset window if expired
		if time.Since(w.startTime) > w.duration {
			w.startTime = time.Now()
			w.total = 0
			w.success = 0
			w.errors = 0
		}

		w.total++
		if success {
			w.success++
		} else {
			w.errors++
		}

		// Calculate compliance (window already locked, use internal method)
		objective := t.objectives[name]
		compliance := t.calculateComplianceInternal(w)
		w.mu.Unlock() // Unlock before recording metrics

		t.sloCompliance.Record(ctx, compliance,
			metric.WithAttributes(
				attribute.String("objective", name),
				attribute.String("service", service),
			),
		)

		// Check for violations
		if compliance < objective.Target {
			t.sloViolations.Add(ctx, 1,
				metric.WithAttributes(
					attribute.String("objective", name),
					attribute.String("service", service),
					attribute.Float64("target", objective.Target),
					attribute.Float64("actual", compliance),
				),
			)
		}
	}
}

// calculateComplianceInternal calculates compliance from a window (assumes window is already locked).
func (t *Tracker) calculateComplianceInternal(w *window) float64 {
	if w.total == 0 {
		return 100.0
	}

	// For availability, calculate success rate
	return (float64(w.success) / float64(w.total)) * 100.0
}

// calculateCompliance calculates the current compliance percentage for an objective.
func (t *Tracker) calculateCompliance(objectiveName string, target float64) float64 {
	w, exists := t.windows[objectiveName]
	if !exists {
		return 100.0
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	return t.calculateComplianceInternal(w)
}

// GetCompliance returns the current compliance percentage for an objective.
func (t *Tracker) GetCompliance(objectiveName string) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	objective, exists := t.objectives[objectiveName]
	if !exists {
		return 0.0
	}

	return t.calculateCompliance(objectiveName, objective.Target)
}

// GetObjectives returns all configured objectives.
func (t *Tracker) GetObjectives() []Objective {
	t.mu.RLock()
	defer t.mu.RUnlock()

	objectives := make([]Objective, 0, len(t.objectives))
	for _, obj := range t.objectives {
		objectives = append(objectives, *obj)
	}
	return objectives
}

// GetViolations returns objectives that are currently in violation.
func (t *Tracker) GetViolations() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	violations := make([]string, 0)
	for name, obj := range t.objectives {
		compliance := t.calculateCompliance(name, obj.Target)
		if compliance < obj.Target {
			violations = append(violations, name)
		}
	}
	return violations
}

// DefaultObjectives returns a set of common SLO objectives.
func DefaultObjectives() []Objective {
	return []Objective{
		{
			Name:        "availability_24h",
			Description: "99.9% availability over 24 hours",
			Target:      99.9,
			Window:      24 * time.Hour,
		},
		{
			Name:        "availability_7d",
			Description: "99.5% availability over 7 days",
			Target:      99.5,
			Window:      7 * 24 * time.Hour,
		},
		{
			Name:        "latency_p95_1h",
			Description: "95% of requests under 500ms in 1 hour",
			Target:      95.0,
			Window:      1 * time.Hour,
		},
	}
}
