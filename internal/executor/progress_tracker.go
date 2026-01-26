package executor

import (
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	api "github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/models"

	"github.com/nicholaspcr/GoDE/internal/store"
)

// progressTracker manages execution progress tracking and callbacks.
type progressTracker struct {
	store                store.Store
	maxVectorsInProgress int
	completionCounters   map[string]*atomic.Int32
	countersMu           sync.RWMutex
}

// newProgressTracker creates a new progress tracker.
func newProgressTracker(store store.Store, maxVectorsInProgress int) *progressTracker {
	return &progressTracker{
		store:                store,
		maxVectorsInProgress: maxVectorsInProgress,
		completionCounters:   make(map[string]*atomic.Int32),
	}
}

// registerExecution creates a completion counter for an execution.
// Returns the counter and a cleanup function that should be called when execution completes.
func (pt *progressTracker) registerExecution(executionID string) (*atomic.Int32, func()) {
	counter := &atomic.Int32{}

	pt.countersMu.Lock()
	pt.completionCounters[executionID] = counter
	pt.countersMu.Unlock()

	cleanup := func() {
		pt.countersMu.Lock()
		delete(pt.completionCounters, executionID)
		pt.countersMu.Unlock()
	}

	return counter, cleanup
}

// createProgressCallback creates a progress callback function for DE execution.
// The callback saves progress to the store and increments the completion counter on final generation.
func (pt *progressTracker) createProgressCallback(
	ctx context.Context,
	executionID string,
	counter *atomic.Int32,
	totalExecutions int32,
) de.ProgressCallback {
	return func(generation int, totalGenerations int, paretoSize int, currentPareto []models.Vector) {
		// If this is the final generation, increment completion counter
		if generation == totalGenerations {
			counter.Add(1)
		}

		// Convert to API vectors (limit to avoid excessive data)
		maxVectors := pt.maxVectorsInProgress
		apiVectors := make([]*api.Vector, 0, min(len(currentPareto), maxVectors))
		for i := 0; i < len(currentPareto) && i < maxVectors; i++ {
			vec := &currentPareto[i]
			apiVectors = append(apiVectors, &api.Vector{
				Elements:         vec.Elements,
				Objectives:       vec.Objectives,
				CrowdingDistance: vec.CrowdingDistance,
			})
		}

		// Read current completion count
		completedCount := counter.Load()

		// #nosec G115 - Values validated in ValidateDEConfig, safe to convert
		progress := &store.ExecutionProgress{
			ExecutionID:         executionID,
			CurrentGeneration:   int32(generation),
			TotalGenerations:    int32(totalGenerations),
			CompletedExecutions: completedCount,
			TotalExecutions:     totalExecutions,
			PartialPareto:       apiVectors,
			UpdatedAt:           time.Now(),
		}

		if err := pt.store.SaveProgress(ctx, progress); err != nil {
			slog.Warn("failed to save progress",
				slog.String("execution_id", executionID),
				slog.String("error", err.Error()),
			)
		}
	}
}
