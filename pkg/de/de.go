package de

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

// Algorithm defines the methods that a Differential Evolution algorithm should
// implement, this method will be executed in each generation.
type Algorithm interface {
	Execute(
		ctx context.Context,
		pareto chan<- []models.Vector,
		maxObjectives chan<- []float64,
	) error
}

// DE contains the necessary methods to setup and execute a Differential
// Evolutionary algorithm.
type de struct {
	algorithm        Algorithm
	config           Config
	constants        Constants
	progressCallback ProgressCallback
}

// New creates a new DE instance based on the configuration options given.
func New(cfg Config, opts ...ModeOptions) (*de, error) {
	m := &de{config: cfg}
	for _, opt := range opts {
		opt(m)
	}

	if m.algorithm == nil {
		return nil, errors.New("no algorithm set")
	}

	return m, nil
}

func (mode *de) Execute(ctx context.Context) ([]models.Vector, [][]float64, error) {
	// Check if context is already cancelled
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	paretoCh := make(chan []models.Vector, mode.config.ParetoChannelLimiter)
	maxObjsCh := make(chan []float64, mode.config.MaxChannelLimiter)
	wgExecs := &sync.WaitGroup{}

	// Collect max objectives from all executions
	allMaxObjs := make([][]float64, 0, mode.constants.Executions)
	var maxObjsMu sync.Mutex

	// Collect pareto results from all executions
	allPareto := make([][]models.Vector, 0, mode.constants.Executions)
	var paretoMu sync.Mutex

	// Collect errors from executions
	var execErrors []error
	var execErrorsMu sync.Mutex

	// WaitGroup for channel consumers
	wgConsumers := &sync.WaitGroup{}

	// Goroutine to consume max objectives (prevents deadlock)
	wgConsumers.Add(1)
	go func() {
		defer wgConsumers.Done()
		for maxObjs := range maxObjsCh {
			maxObjsMu.Lock()
			allMaxObjs = append(allMaxObjs, maxObjs)
			maxObjsMu.Unlock()
		}
	}()

	// Goroutine to consume pareto results (prevents deadlock)
	wgConsumers.Add(1)
	go func() {
		defer wgConsumers.Done()
		for pareto := range paretoCh {
			paretoMu.Lock()
			allPareto = append(allPareto, pareto)
			paretoMu.Unlock()
		}
	}()

	// Runs algorithm for Executions amount of times.
	for i := range mode.constants.Executions {
		wgExecs.Add(1)
		// Initialize worker responsible for DE execution.
		go func(idx int) {
			defer wgExecs.Done()
			defer func() {
				if r := recover(); r != nil {
					slog.Error("panic recovered in DE execution goroutine",
						slog.Int("execution_id", idx),
						slog.Any("panic", r),
					)
					execErrorsMu.Lock()
					execErrors = append(execErrors, fmt.Errorf("panic in execution %d: %v", idx, r))
					execErrorsMu.Unlock()
				}
			}()
			// running the algorithm execution.
			if err := mode.algorithm.Execute(
				WithContextExecutionNumber(ctx, idx),
				paretoCh,
				maxObjsCh,
			); err != nil {
				execErrorsMu.Lock()
				execErrors = append(execErrors, fmt.Errorf("execution %d: %w", idx, err))
				execErrorsMu.Unlock()

				// Check if error is due to cancellation
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					slog.Info("Execution cancelled",
						slog.Int("Execution", idx),
					)
				} else {
					slog.Error("Unexpected error while executing the algorithm",
						slog.Int("Execution", idx),
						slog.String("error", err.Error()),
					)
				}
			}
		}(i)
	}

	wgExecs.Wait()
	close(paretoCh)
	close(maxObjsCh)
	wgConsumers.Wait()

	// If all executions failed, return the combined error
	if len(execErrors) == mode.constants.Executions {
		return nil, nil, errors.Join(execErrors...)
	}

	// Check if cancelled before filtering
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	now := time.Now()
	finalPareto := mode.filterCollectedPareto(ctx, allPareto)
	slog.Info("Filtering Pareto", slog.Duration("time", time.Since(now)))

	return finalPareto, allMaxObjs, nil
}

func (mode *de) filterCollectedPareto(
	ctx context.Context, allPareto [][]models.Vector,
) []models.Vector {
	finalPareto := make([]models.Vector, 0, 2000)
	for _, pareto := range allPareto {
		// Check for cancellation between batches
		if ctx.Err() != nil {
			return finalPareto
		}

		// Use incremental update instead of full re-ranking
		var rankZero []models.Vector
		finalPareto, rankZero = IncrementalParetoUpdate(
			ctx, finalPareto, pareto, mode.config.ResultLimiter,
		)
		_ = rankZero // Available for future features
	}
	return finalPareto
}
