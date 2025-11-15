package redis

import (
	"time"

	"github.com/sony/gobreaker"
)

// BreakerConfig holds circuit breaker configuration.
type BreakerConfig struct {
	MaxFailures uint32        // Number of failures before opening (default: 5)
	Timeout     time.Duration // Time to wait before half-open (default: 60s)
	MaxRequests uint32        // Max requests in half-open state (default: 10)
}

// DefaultBreakerConfig returns sensible defaults for Redis circuit breaker.
func DefaultBreakerConfig() BreakerConfig {
	return BreakerConfig{
		MaxFailures: 5,
		Timeout:     60 * time.Second,
		MaxRequests: 10,
	}
}

// newCircuitBreaker creates a circuit breaker with the given configuration.
func newCircuitBreaker(name string, cfg BreakerConfig) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: cfg.MaxRequests,
		Interval:    60 * time.Second, // Clear failure counts every 60s
		Timeout:     cfg.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Open circuit if:
			// 1. At least 3 requests have been made
			// 2. Failure ratio >= 60%
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			// Log state changes for observability
			// Note: using fmt.Printf here for simplicity, could use slog
			// fmt.Printf("Circuit breaker %s: %s -> %s\n", name, from, to)
		},
	}

	return gobreaker.NewCircuitBreaker(settings)
}
