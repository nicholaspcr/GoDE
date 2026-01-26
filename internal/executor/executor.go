// Package executor provides background execution of Differential Evolution algorithms.
package executor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"runtime/debug"
	"sync"
	"time"

	"github.com/google/uuid"
	api "github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/de/gde3"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/variants"

	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/internal/telemetry"
)

// Executor manages background execution of DE algorithms.
type Executor struct {
	store           store.Store
	executionTTL    time.Duration
	resultTTL       time.Duration
	progressTTL     time.Duration
	workers         *workerPool
	progress        *progressTracker
	activeExecs     map[string]context.CancelFunc
	activeExecsMu   sync.RWMutex
	problemRegistry map[string]problems.Interface
	variantRegistry map[string]variants.Interface
}

// Config holds configuration for the Executor.
type Config struct {
	Store                store.Store
	MaxWorkers           int
	MaxVectorsInProgress int // Maximum vectors to include in progress updates (default: 100)
	ExecutionTTL         time.Duration
	ResultTTL            time.Duration
	ProgressTTL          time.Duration
	Metrics              *telemetry.Metrics
}

// New creates a new Executor instance.
func New(cfg Config) *Executor {
	maxVectorsInProgress := cfg.MaxVectorsInProgress
	if maxVectorsInProgress <= 0 {
		maxVectorsInProgress = 100 // Default value
	}

	e := &Executor{
		store:           cfg.Store,
		executionTTL:    cfg.ExecutionTTL,
		resultTTL:       cfg.ResultTTL,
		progressTTL:     cfg.ProgressTTL,
		workers:         newWorkerPool(cfg.MaxWorkers, cfg.Metrics),
		progress:        newProgressTracker(cfg.Store, maxVectorsInProgress),
		activeExecs:     make(map[string]context.CancelFunc),
		problemRegistry: make(map[string]problems.Interface),
		variantRegistry: make(map[string]variants.Interface),
	}

	return e
}

// RegisterProblem registers a problem implementation.
func (e *Executor) RegisterProblem(name string, p problems.Interface) {
	e.problemRegistry[name] = p
}

// RegisterVariant registers a variant implementation.
func (e *Executor) RegisterVariant(name string, v variants.Interface) {
	e.variantRegistry[name] = v
}

// SubmitExecution submits a new DE execution to run in the background.
func (e *Executor) SubmitExecution(ctx context.Context, userID, algorithm, problem, variant string, config *api.DEConfig) (string, error) {
	// Validate problem and variant exist before creating execution record
	if _, exists := e.problemRegistry[problem]; !exists {
		return "", fmt.Errorf("unknown problem: %s", problem)
	}
	if _, exists := e.variantRegistry[variant]; !exists {
		return "", fmt.Errorf("unknown variant: %s", variant)
	}

	// Generate execution ID
	executionID := uuid.New().String()

	// Create execution record
	execution := &store.Execution{
		ID:        executionID,
		UserID:    userID,
		Status:    store.ExecutionStatusPending,
		Config:    config,
		Algorithm: algorithm,
		Variant:   variant,
		Problem:   problem,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := e.store.CreateExecution(ctx, execution); err != nil {
		return "", fmt.Errorf("failed to create execution: %w", err)
	}

	// Extract trace context values from parent context before spawning goroutine.
	// This allows background execution to propagate parent context values (e.g., tracing spans)
	// while still being cancellable independently.
	parentCtx := context.WithoutCancel(ctx)

	// Submit to worker pool (non-blocking)
	go e.executeInBackground(parentCtx, executionID, userID, algorithm, problem, variant, config)

	return executionID, nil
}

// CancelExecution cancels a running execution.
func (e *Executor) CancelExecution(ctx context.Context, executionID, userID string) error {
	// Mark for cancellation in store
	if err := e.store.MarkExecutionForCancellation(ctx, executionID, userID); err != nil {
		return err
	}

	// Cancel the context if execution is active
	e.activeExecsMu.Lock()
	if cancel, exists := e.activeExecs[executionID]; exists {
		cancel()
		delete(e.activeExecs, executionID)
	}
	e.activeExecsMu.Unlock()

	return nil
}

// Shutdown gracefully stops all active executions and waits for workers to finish.
// It cancels all running executions and waits up to 30 seconds for them to complete.
func (e *Executor) Shutdown(ctx context.Context) error {
	// Copy cancel functions to avoid holding lock during cancellation
	e.activeExecsMu.Lock()
	activeCount := len(e.activeExecs)
	cancelFuncs := make([]context.CancelFunc, 0, activeCount)
	executionIDs := make([]string, 0, activeCount)
	for id, cancel := range e.activeExecs {
		cancelFuncs = append(cancelFuncs, cancel)
		executionIDs = append(executionIDs, id)
	}
	e.activeExecsMu.Unlock()

	slog.Info("shutting down executor",
		slog.Int("active_executions", activeCount),
		slog.Int("max_workers", e.workers.maxWorkers),
	)

	// Cancel all active executions
	for i, cancel := range cancelFuncs {
		slog.Info("cancelling execution during shutdown",
			slog.String("execution_id", executionIDs[i]),
		)
		cancel()
	}

	// Wait for all workers to finish by acquiring all slots
	// This ensures all executeInBackground goroutines have completed
	shutdownTimeout := 30 * time.Second
	deadline := time.Now().Add(shutdownTimeout)

	acquired := 0
	for acquired < e.workers.maxWorkers {
		remainingTime := time.Until(deadline)
		if remainingTime <= 0 {
			slog.Warn("shutdown timeout - some workers still active",
				slog.Int("workers_remaining", e.workers.maxWorkers-acquired),
			)
			return fmt.Errorf("shutdown timeout: %d workers still active", e.workers.maxWorkers-acquired)
		}

		// Use NewTimer instead of time.After to avoid timer leaks
		timer := time.NewTimer(remainingTime)
		select {
		case e.workers.pool <- struct{}{}:
			timer.Stop()
			acquired++
		case <-ctx.Done():
			timer.Stop()
			slog.Warn("shutdown context cancelled",
				slog.Int("workers_remaining", e.workers.maxWorkers-acquired),
			)
			return ctx.Err()
		case <-timer.C:
			slog.Warn("shutdown timeout - some workers still active",
				slog.Int("workers_remaining", e.workers.maxWorkers-acquired),
			)
			return fmt.Errorf("shutdown timeout: %d workers still active", e.workers.maxWorkers-acquired)
		}
	}

	slog.Info("executor shutdown complete")
	return nil
}

func (e *Executor) executeInBackground(parentCtx context.Context, executionID, userID, algorithm, problem, variant string, config *api.DEConfig) {
	// Create base context for worker acquisition
	ctx := context.Background()

	// Acquire worker slot (blocks until available)
	releaseWorker, _ := e.workers.acquireWorker(ctx)
	defer releaseWorker()

	// Create cancellable context derived from parent (preserves trace context).
	// The parent context is already detached from request cancellation via WithoutCancel.
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	// Register active execution
	e.activeExecsMu.Lock()
	e.activeExecs[executionID] = cancel
	e.activeExecsMu.Unlock()

	// Cleanup on exit
	defer func() {
		e.activeExecsMu.Lock()
		delete(e.activeExecs, executionID)
		e.activeExecsMu.Unlock()

		if r := recover(); r != nil {
			stack := debug.Stack()
			slog.Error("panic in execution",
				slog.String("execution_id", executionID),
				slog.Any("panic", r),
				slog.String("stack", string(stack)),
			)
			if updateErr := e.store.UpdateExecutionStatus(ctx, executionID, store.ExecutionStatusFailed, fmt.Sprintf("panic: %v", r)); updateErr != nil {
				slog.Error("failed to update execution status after panic",
					slog.String("execution_id", executionID),
					slog.String("panic", fmt.Sprintf("%v", r)),
					slog.Any("update_error", updateErr),
				)
			}
		}
	}()

	// Update status to running
	if err := e.store.UpdateExecutionStatus(ctx, executionID, store.ExecutionStatusRunning, ""); err != nil {
		slog.Error("failed to update execution status",
			slog.String("execution_id", executionID),
			slog.String("error", err.Error()),
		)
		return
	}

	// Execute the algorithm
	pareto, maxObjs, err := e.runAlgorithm(ctx, executionID, problem, variant, config)
	if err != nil {
		var updateErr error
		if errors.Is(err, context.Canceled) {
			updateErr = e.store.UpdateExecutionStatus(ctx, executionID, store.ExecutionStatusCancelled, "")
			if updateErr != nil {
				slog.Error("failed to update execution status to cancelled",
					slog.String("execution_id", executionID),
					slog.Any("update_error", updateErr),
				)
			}
		} else {
			updateErr = e.store.UpdateExecutionStatus(ctx, executionID, store.ExecutionStatusFailed, err.Error())
			if updateErr != nil {
				slog.Error("failed to update execution status to failed",
					slog.String("execution_id", executionID),
					slog.String("original_error", err.Error()),
					slog.Any("update_error", updateErr),
				)
			}
		}
		slog.Info("execution failed",
			slog.String("execution_id", executionID),
			slog.String("error", err.Error()),
		)
		return
	}

	// Save results
	paretoID, err := e.saveResults(ctx, userID, algorithm, problem, variant, pareto, maxObjs)
	if err != nil {
		if updateErr := e.store.UpdateExecutionStatus(ctx, executionID, store.ExecutionStatusFailed, err.Error()); updateErr != nil {
			slog.Error("failed to update execution status after save failure",
				slog.String("execution_id", executionID),
				slog.String("save_error", err.Error()),
				slog.Any("update_error", updateErr),
			)
		}
		slog.Error("failed to save results",
			slog.String("execution_id", executionID),
			slog.String("error", err.Error()),
		)
		return
	}

	// Update execution with result
	if err := e.store.UpdateExecutionResult(ctx, executionID, paretoID); err != nil {
		slog.Error("failed to update execution result",
			slog.String("execution_id", executionID),
			slog.String("error", err.Error()),
		)
	}

	// Mark as completed
	if err := e.store.UpdateExecutionStatus(ctx, executionID, store.ExecutionStatusCompleted, ""); err != nil {
		slog.Error("failed to mark execution as completed",
			slog.String("execution_id", executionID),
			slog.String("error", err.Error()),
		)
	}
}

func (e *Executor) runAlgorithm(ctx context.Context, executionID, problemName, variantName string, config *api.DEConfig) ([]models.Vector, [][]float64, error) {
	// Register execution for progress tracking
	counter, cleanup := e.progress.registerExecution(executionID)
	defer cleanup()

	// Get problem
	problemImpl, exists := e.problemRegistry[problemName]
	if !exists {
		return nil, nil, fmt.Errorf("unknown problem: %s", problemName)
	}

	// Get variant
	variantImpl, exists := e.variantRegistry[variantName]
	if !exists {
		return nil, nil, fmt.Errorf("unknown variant: %s", variantName)
	}

	// Build population parameters
	popParams := models.PopulationParams{
		PopulationSize: int(config.PopulationSize),
		DimensionSize:  int(config.DimensionsSize),
		ObjectivesSize: int(config.ObjectivesSize),
		FloorRange:     make([]float64, config.DimensionsSize),
		CeilRange:      make([]float64, config.DimensionsSize),
	}

	for i := range config.DimensionsSize {
		popParams.FloorRange[i] = float64(config.FloorLimiter)
		popParams.CeilRange[i] = float64(config.CeilLimiter)
	}

	// Generate initial population
	// #nosec G404 - Using math/rand for DE algorithm randomness, not cryptographic purposes
	initialPop, err := models.GeneratePopulation(popParams, rand.New(rand.NewSource(time.Now().UnixNano())))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate population: %w", err)
	}

	// Build GDE3 constants
	gde3Config := config.GetGde3()
	if gde3Config == nil {
		return nil, nil, fmt.Errorf("GDE3 configuration is required")
	}

	gde3Constants := gde3.Constants{
		DE: de.Constants{
			Executions:    int(config.Executions),
			Generations:   int(config.Generations),
			Dimensions:    int(config.DimensionsSize),
			ObjFuncAmount: int(config.ObjectivesSize),
		},
		CR: float64(gde3Config.Cr),
		F:  float64(gde3Config.F),
		P:  float64(gde3Config.P),
	}

	// Create progress callback using progress tracker
	progressCallback := e.progress.createProgressCallback(
		ctx,
		executionID,
		counter,
		int32(config.Executions),
	)

	// Create GDE3 algorithm
	algorithm := gde3.New(
		gde3.WithProblem(problemImpl),
		gde3.WithVariant(variantImpl),
		gde3.WithPopulationParams(popParams),
		gde3.WithConstants(gde3Constants),
		gde3.WithInitialPopulation(initialPop),
		gde3.WithProgressCallback(progressCallback),
	)

	// Create DE mode
	deConfig := de.Config{
		ParetoChannelLimiter: int(config.Executions),
		MaxChannelLimiter:    int(config.Executions),
		ResultLimiter:        1000,
	}

	mode, err := de.New(deConfig,
		de.WithAlgorithm(algorithm),
		de.WithExecutions(int(config.Executions)),
		de.WithGenerations(int(config.Generations)),
		de.WithDimensions(int(config.DimensionsSize)),
		de.WithObjFuncAmount(int(config.ObjectivesSize)),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create DE mode: %w", err)
	}

	// Execute
	return mode.Execute(ctx)
}

func (e *Executor) saveResults(ctx context.Context, userID, algorithm, problem, variant string, pareto []models.Vector, maxObjs [][]float64) (uint64, error) {
	// Convert to API vectors
	apiVectors := make([]*api.Vector, len(pareto))
	for i := range pareto {
		vec := &pareto[i]
		apiVectors[i] = &api.Vector{
			Elements:         vec.Elements,
			Objectives:       vec.Objectives,
			CrowdingDistance: vec.CrowdingDistance,
		}
	}

	// Convert max objectives
	storeMaxObjs := make([]*store.MaxObjectives, len(maxObjs))
	for i := range maxObjs {
		storeMaxObjs[i] = &store.MaxObjectives{Values: maxObjs[i]}
	}

	// Create pareto set
	paretoSet := &store.ParetoSet{
		UserID:        userID,
		Algorithm:     algorithm,
		Problem:       problem,
		Variant:       variant,
		Vectors:       apiVectors,
		MaxObjectives: storeMaxObjs,
		CreatedAt:     time.Now(),
	}

	if err := e.store.CreateParetoSet(ctx, paretoSet); err != nil {
		return 0, fmt.Errorf("failed to create pareto set: %w", err)
	}

	return paretoSet.ID, nil
}
