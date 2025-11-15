package redis

import (
	"testing"
	"time"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultBreakerConfig(t *testing.T) {
	cfg := DefaultBreakerConfig()

	assert.Equal(t, uint32(5), cfg.MaxFailures)
	assert.Equal(t, 60*time.Second, cfg.Timeout)
	assert.Equal(t, uint32(10), cfg.MaxRequests)
}

func TestNewCircuitBreaker(t *testing.T) {
	cfg := BreakerConfig{
		MaxFailures: 3,
		Timeout:     30 * time.Second,
		MaxRequests: 5,
	}

	breaker := newCircuitBreaker("test-breaker", cfg)
	require.NotNil(t, breaker)

	// Verify initial state is closed
	state := breaker.State()
	assert.Equal(t, gobreaker.StateClosed, state)
}

func TestCircuitBreaker_OpensAfterFailures(t *testing.T) {
	cfg := BreakerConfig{
		MaxFailures: 3,
		Timeout:     100 * time.Millisecond,
		MaxRequests: 2,
	}

	breaker := newCircuitBreaker("test-breaker", cfg)

	// Execute successful request
	_, err := breaker.Execute(func() (interface{}, error) {
		return "success", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, gobreaker.StateClosed, breaker.State())

	// Execute 3 failing requests to trigger circuit breaker
	// Need 3 requests with >= 60% failure ratio
	for i := 0; i < 3; i++ {
		_, _ = breaker.Execute(func() (interface{}, error) {
			return nil, assert.AnError
		})
	}

	// Circuit should be open now (after 3 requests with 100% failure)
	assert.Equal(t, gobreaker.StateOpen, breaker.State())
}

func TestCircuitBreaker_HalfOpenAfterTimeout(t *testing.T) {
	cfg := BreakerConfig{
		MaxFailures: 2,
		Timeout:     100 * time.Millisecond,
		MaxRequests: 2,
	}

	breaker := newCircuitBreaker("test-breaker", cfg)

	// Trigger circuit to open (3 failures)
	for i := 0; i < 3; i++ {
		_, _ = breaker.Execute(func() (interface{}, error) {
			return nil, assert.AnError
		})
	}

	assert.Equal(t, gobreaker.StateOpen, breaker.State())

	// Wait for timeout to allow half-open
	time.Sleep(150 * time.Millisecond)

	// Next request should transition to half-open
	_, _ = breaker.Execute(func() (interface{}, error) {
		return "success", nil
	})

	// Should now be in half-open or closed state (depending on success)
	state := breaker.State()
	assert.True(t, state == gobreaker.StateHalfOpen || state == gobreaker.StateClosed)
}

func TestCircuitBreaker_FailureRatioThreshold(t *testing.T) {
	cfg := BreakerConfig{
		MaxFailures: 5,
		Timeout:     60 * time.Second,
		MaxRequests: 10,
	}

	breaker := newCircuitBreaker("test-breaker", cfg)

	// Execute 3 requests: 1 success, 2 failures
	// Failure ratio: 2/3 = 66.6% > 60% threshold
	// Should trigger circuit breaker
	_, _ = breaker.Execute(func() (interface{}, error) {
		return "success", nil
	})
	_, _ = breaker.Execute(func() (interface{}, error) {
		return nil, assert.AnError
	})
	_, _ = breaker.Execute(func() (interface{}, error) {
		return nil, assert.AnError
	})

	// Circuit should be open due to failure ratio
	assert.Equal(t, gobreaker.StateOpen, breaker.State())
}

func TestCircuitBreaker_MinimumRequestsThreshold(t *testing.T) {
	cfg := BreakerConfig{
		MaxFailures: 5,
		Timeout:     60 * time.Second,
		MaxRequests: 10,
	}

	breaker := newCircuitBreaker("test-breaker", cfg)

	// Execute only 2 failing requests (less than minimum of 3)
	// Should NOT open circuit
	_, _ = breaker.Execute(func() (interface{}, error) {
		return nil, assert.AnError
	})
	_, _ = breaker.Execute(func() (interface{}, error) {
		return nil, assert.AnError
	})

	// Circuit should still be closed (not enough requests)
	assert.Equal(t, gobreaker.StateClosed, breaker.State())
}

func TestCircuitBreaker_SuccessfulRecovery(t *testing.T) {
	cfg := BreakerConfig{
		MaxFailures: 3,
		Timeout:     50 * time.Millisecond,
		MaxRequests: 2,
	}

	breaker := newCircuitBreaker("test-breaker", cfg)

	// Open the circuit
	for i := 0; i < 3; i++ {
		_, _ = breaker.Execute(func() (interface{}, error) {
			return nil, assert.AnError
		})
	}
	assert.Equal(t, gobreaker.StateOpen, breaker.State())

	// Wait for timeout
	time.Sleep(100 * time.Millisecond)

	// Execute successful requests to close circuit
	for i := 0; i < 3; i++ {
		_, err := breaker.Execute(func() (interface{}, error) {
			return "success", nil
		})
		assert.NoError(t, err)
	}

	// Circuit should be closed again
	assert.Equal(t, gobreaker.StateClosed, breaker.State())
}

func TestCircuitBreaker_CustomConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		config BreakerConfig
	}{
		{
			name: "high_threshold",
			config: BreakerConfig{
				MaxFailures: 10,
				Timeout:     30 * time.Second,
				MaxRequests: 5,
			},
		},
		{
			name: "low_threshold",
			config: BreakerConfig{
				MaxFailures: 2,
				Timeout:     10 * time.Second,
				MaxRequests: 3,
			},
		},
		{
			name: "quick_recovery",
			config: BreakerConfig{
				MaxFailures: 5,
				Timeout:     5 * time.Second,
				MaxRequests: 15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			breaker := newCircuitBreaker(tt.name, tt.config)
			require.NotNil(t, breaker)
			assert.Equal(t, gobreaker.StateClosed, breaker.State())
		})
	}
}
