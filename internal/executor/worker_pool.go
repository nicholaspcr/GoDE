package executor

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/nicholaspcr/GoDE/internal/telemetry"
)

// workerPool manages worker concurrency and metrics.
type workerPool struct {
	pool              chan struct{}    // Semaphore for limiting concurrency
	maxWorkers        int              // Maximum number of concurrent workers
	activeWorkerCount atomic.Int64     // Current number of active workers (safe to read)
	metrics           *telemetry.Metrics
}

// newWorkerPool creates a new worker pool with the specified capacity.
func newWorkerPool(maxWorkers int, metrics *telemetry.Metrics) *workerPool {
	wp := &workerPool{
		pool:       make(chan struct{}, maxWorkers),
		maxWorkers: maxWorkers,
		metrics:    metrics,
	}

	// Initialize total workers metric
	if metrics != nil && metrics.ExecutorWorkersTotal != nil {
		metrics.ExecutorWorkersTotal.Add(context.Background(), int64(maxWorkers))
	}

	return wp
}

// acquireWorker acquires a worker slot and returns a release function.
// This function blocks until a worker slot is available or the context is cancelled.
// The release function must be called (typically via defer) to return the worker to the pool.
// Returns an error if the context is cancelled while waiting.
func (wp *workerPool) acquireWorker(ctx context.Context) (releaseFunc func(), queueWait time.Duration, err error) {
	queueStart := time.Now()

	// Acquire worker slot, respecting context cancellation
	select {
	case wp.pool <- struct{}{}:
		// slot acquired
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
	queueWait = time.Since(queueStart)

	// Increment active worker count
	activeWorkers := wp.activeWorkerCount.Add(1)

	// Record queue wait and active workers metrics
	if wp.metrics != nil {
		if wp.metrics.ExecutorQueueWaitDuration != nil {
			wp.metrics.ExecutorQueueWaitDuration.Record(ctx, queueWait.Seconds())
		}

		if wp.metrics.ExecutorWorkersActive != nil {
			wp.metrics.ExecutorWorkersActive.Add(ctx, 1)
		}

		if wp.metrics.ExecutorUtilizationPercent != nil {
			utilization := float64(activeWorkers) / float64(wp.maxWorkers) * 100
			wp.metrics.ExecutorUtilizationPercent.Record(ctx, utilization)
		}
	}

	// Return release function
	releaseFunc = func() {
		// Decrement active worker count
		wp.activeWorkerCount.Add(-1)

		// Release worker slot
		<-wp.pool

		// Decrement active workers metric
		if wp.metrics != nil && wp.metrics.ExecutorWorkersActive != nil {
			wp.metrics.ExecutorWorkersActive.Add(ctx, -1)
		}
	}

	return releaseFunc, queueWait, nil
}

// getActiveCount returns the current number of active workers.
// This is safe to call from any goroutine.
func (wp *workerPool) getActiveCount() int64 {
	return wp.activeWorkerCount.Load()
}
