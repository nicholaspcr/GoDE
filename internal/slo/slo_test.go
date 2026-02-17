package slo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTracker(t *testing.T) {
	objectives := []Objective{
		{
			Name:        "availability",
			Description: "99.9% availability",
			Target:      99.9,
			Window:      1 * time.Hour,
		},
	}

	tracker, err := NewTracker(context.Background(), objectives)
	require.NoError(t, err)
	require.NotNil(t, tracker)
	assert.Equal(t, 1, len(tracker.objectives))
}

func TestTracker_RecordRequest(t *testing.T) {
	objectives := []Objective{
		{
			Name:        "availability",
			Description: "99% availability",
			Target:      99.0,
			Window:      1 * time.Hour,
		},
	}

	tracker, err := NewTracker(context.Background(), objectives)
	require.NoError(t, err)

	ctx := context.Background()

	// Record 100 successful requests
	for range 100 {
		tracker.RecordRequest(ctx, "test-service", true, 0.1)
	}

	compliance := tracker.GetCompliance("availability")
	assert.Equal(t, 100.0, compliance)
}

func TestTracker_RecordRequest_WithErrors(t *testing.T) {
	objectives := []Objective{
		{
			Name:        "availability",
			Description: "99% availability",
			Target:      99.0,
			Window:      1 * time.Hour,
		},
	}

	tracker, err := NewTracker(context.Background(), objectives)
	require.NoError(t, err)

	ctx := context.Background()

	// Record 95 successful requests and 5 errors
	for range 95 {
		tracker.RecordRequest(ctx, "test-service", true, 0.1)
	}
	for range 5 {
		tracker.RecordRequest(ctx, "test-service", false, 0.1)
	}

	compliance := tracker.GetCompliance("availability")
	assert.Equal(t, 95.0, compliance)
}

func TestTracker_GetViolations(t *testing.T) {
	objectives := []Objective{
		{
			Name:        "availability",
			Description: "99% availability",
			Target:      99.0,
			Window:      1 * time.Hour,
		},
	}

	tracker, err := NewTracker(context.Background(), objectives)
	require.NoError(t, err)

	ctx := context.Background()

	// Record requests that violate SLO (90% success rate)
	for range 90 {
		tracker.RecordRequest(ctx, "test-service", true, 0.1)
	}
	for range 10 {
		tracker.RecordRequest(ctx, "test-service", false, 0.1)
	}

	violations := tracker.GetViolations()
	assert.Equal(t, 1, len(violations))
	assert.Contains(t, violations, "availability")
}

func TestTracker_NoViolations(t *testing.T) {
	objectives := []Objective{
		{
			Name:        "availability",
			Description: "99% availability",
			Target:      99.0,
			Window:      1 * time.Hour,
		},
	}

	tracker, err := NewTracker(context.Background(), objectives)
	require.NoError(t, err)

	ctx := context.Background()

	// Record requests that meet SLO (99.5% success rate)
	for range 995 {
		tracker.RecordRequest(ctx, "test-service", true, 0.1)
	}
	for range 5 {
		tracker.RecordRequest(ctx, "test-service", false, 0.1)
	}

	violations := tracker.GetViolations()
	assert.Equal(t, 0, len(violations))
}

func TestTracker_GetObjectives(t *testing.T) {
	objectives := []Objective{
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
	}

	tracker, err := NewTracker(context.Background(), objectives)
	require.NoError(t, err)

	retrievedObjectives := tracker.GetObjectives()
	assert.Equal(t, 2, len(retrievedObjectives))
}

func TestTracker_MultipleObjectives(t *testing.T) {
	objectives := []Objective{
		{
			Name:        "availability_strict",
			Description: "99.9% availability",
			Target:      99.9,
			Window:      1 * time.Hour,
		},
		{
			Name:        "availability_relaxed",
			Description: "95% availability",
			Target:      95.0,
			Window:      1 * time.Hour,
		},
	}

	tracker, err := NewTracker(context.Background(), objectives)
	require.NoError(t, err)

	ctx := context.Background()

	// Record 97% success rate
	for range 97 {
		tracker.RecordRequest(ctx, "test-service", true, 0.1)
	}
	for range 3 {
		tracker.RecordRequest(ctx, "test-service", false, 0.1)
	}

	// Should violate strict but not relaxed
	violations := tracker.GetViolations()
	assert.Equal(t, 1, len(violations))
	assert.Contains(t, violations, "availability_strict")

	strictCompliance := tracker.GetCompliance("availability_strict")
	relaxedCompliance := tracker.GetCompliance("availability_relaxed")

	assert.Equal(t, 97.0, strictCompliance)
	assert.Equal(t, 97.0, relaxedCompliance)
}

func TestDefaultObjectives(t *testing.T) {
	objectives := DefaultObjectives()
	assert.Greater(t, len(objectives), 0)

	// Verify all objectives have required fields
	for _, obj := range objectives {
		assert.NotEmpty(t, obj.Name)
		assert.NotEmpty(t, obj.Description)
		assert.Greater(t, obj.Target, 0.0)
		assert.LessOrEqual(t, obj.Target, 100.0)
		assert.Greater(t, obj.Window, time.Duration(0))
	}
}

func TestTracker_EmptyCompliance(t *testing.T) {
	objectives := []Objective{
		{
			Name:        "availability",
			Description: "99% availability",
			Target:      99.0,
			Window:      1 * time.Hour,
		},
	}

	tracker, err := NewTracker(context.Background(), objectives)
	require.NoError(t, err)

	// No requests recorded, should return 100% compliance
	compliance := tracker.GetCompliance("availability")
	assert.Equal(t, 100.0, compliance)
}

func TestTracker_NonExistentObjective(t *testing.T) {
	objectives := []Objective{
		{
			Name:        "availability",
			Description: "99% availability",
			Target:      99.0,
			Window:      1 * time.Hour,
		},
	}

	tracker, err := NewTracker(context.Background(), objectives)
	require.NoError(t, err)

	compliance := tracker.GetCompliance("nonexistent")
	assert.Equal(t, 0.0, compliance)
}
