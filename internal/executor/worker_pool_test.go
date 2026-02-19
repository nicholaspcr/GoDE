package executor

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkerPool_AcquireRelease(t *testing.T) {
	wp := newWorkerPool(3, nil)
	ctx := context.Background()

	releaseFunc, queueWait, err := wp.acquireWorker(ctx)
	require.NoError(t, err)
	require.NotNil(t, releaseFunc)
	assert.GreaterOrEqual(t, queueWait, time.Duration(0))
	assert.Equal(t, int64(1), wp.getActiveCount())

	releaseFunc()
	assert.Equal(t, int64(0), wp.getActiveCount())
}

func TestWorkerPool_GetActiveCount(t *testing.T) {
	wp := newWorkerPool(5, nil)
	ctx := context.Background()

	assert.Equal(t, int64(0), wp.getActiveCount())

	release1, _, err := wp.acquireWorker(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(1), wp.getActiveCount())

	release2, _, err := wp.acquireWorker(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(2), wp.getActiveCount())

	release1()
	assert.Equal(t, int64(1), wp.getActiveCount())

	release2()
	assert.Equal(t, int64(0), wp.getActiveCount())
}

func TestWorkerPool_ContextCancellation(t *testing.T) {
	// Fill pool completely
	wp := newWorkerPool(2, nil)
	ctx := context.Background()

	release1, _, _ := wp.acquireWorker(ctx)
	release2, _, _ := wp.acquireWorker(ctx)
	defer release1()
	defer release2()

	// Pool is full; cancelling context should unblock acquireWorker
	cancelCtx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, _, err := wp.acquireWorker(cancelCtx)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestWorkerPool_ConcurrentAcquire(t *testing.T) {
	maxWorkers := 5
	wp := newWorkerPool(maxWorkers, nil)
	ctx := context.Background()

	// Concurrently acquire all slots
	releases := make([]func(), maxWorkers)
	var wg sync.WaitGroup

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			rel, _, err := wp.acquireWorker(ctx)
			require.NoError(t, err)
			releases[idx] = rel
		}(i)
	}
	wg.Wait()

	assert.Equal(t, int64(maxWorkers), wp.getActiveCount())

	// Release all
	for _, rel := range releases {
		rel()
	}
	assert.Equal(t, int64(0), wp.getActiveCount())
}

func TestWorkerPool_QueueWaitMeasured(t *testing.T) {
	// With 1 worker that's held, second acquire should measure some wait time.
	wp := newWorkerPool(1, nil)
	ctx := context.Background()

	release1, _, err := wp.acquireWorker(ctx)
	require.NoError(t, err)

	// Release after short delay in background
	go func() {
		time.Sleep(20 * time.Millisecond)
		release1()
	}()

	_, queueWait, err := wp.acquireWorker(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, queueWait, 10*time.Millisecond, "should have waited for worker slot")
}

func TestWorkerPool_ConcurrentActiveCountAccuracy(t *testing.T) {
	maxWorkers := 10
	wp := newWorkerPool(maxWorkers, nil)
	ctx := context.Background()

	var peakActive atomic.Int64
	var wg sync.WaitGroup
	totalTasks := 50

	for i := 0; i < totalTasks; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rel, _, err := wp.acquireWorker(ctx)
			require.NoError(t, err)

			current := wp.getActiveCount()
			for {
				peak := peakActive.Load()
				if current <= peak || peakActive.CompareAndSwap(peak, current) {
					break
				}
			}

			time.Sleep(2 * time.Millisecond)
			rel()
		}()
	}
	wg.Wait()

	assert.LessOrEqual(t, peakActive.Load(), int64(maxWorkers),
		"active count should never exceed maxWorkers")
	assert.Equal(t, int64(0), wp.getActiveCount(), "should be zero after all done")
}
