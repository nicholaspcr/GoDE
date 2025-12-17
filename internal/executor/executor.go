// Package executor provides background execution of Differential Evolution algorithms.
package executor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	api "github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/de/gde3"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/variants"

	"github.com/nicholaspcr/GoDE/internal/store"
)

// Executor manages background execution of DE algorithms.
type Executor struct {
	store              store.Store
	maxWorkers         int
	executionTTL       time.Duration
	resultTTL          time.Duration
	progressTTL        time.Duration
	workerPool         chan struct{}
	activeExecs        map[string]context.CancelFunc
	activeExecsMu      sync.RWMutex
	problemRegistry    map[string]problems.Interface
	variantRegistry    map[string]variants.Interface
	completionCounters map[string]*atomic.Int32
	countersMu         sync.RWMutex
}

// Config holds configuration for the Executor.
type Config struct {
	Store        store.Store
	MaxWorkers   int
	ExecutionTTL time.Duration
	ResultTTL    time.Duration
	ProgressTTL  time.Duration
}

// New creates a new Executor instance.
func New(cfg Config) *Executor {
	return &Executor{
		store:              cfg.Store,
		maxWorkers:         cfg.MaxWorkers,
		executionTTL:       cfg.ExecutionTTL,
		resultTTL:          cfg.ResultTTL,
		progressTTL:        cfg.ProgressTTL,
		workerPool:         make(chan struct{}, cfg.MaxWorkers),
		activeExecs:        make(map[string]context.CancelFunc),
		problemRegistry:    make(map[string]problems.Interface),
		variantRegistry:    make(map[string]variants.Interface),
		completionCounters: make(map[string]*atomic.Int32),
	}
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
	// Generate execution ID
	executionID := uuid.New().String()

	// Create execution record
	execution := &store.Execution{
		ID:        executionID,
		UserID:    userID,
		Status:    store.ExecutionStatusPending,
		Config:    config,
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

func (e *Executor) executeInBackground(parentCtx context.Context, executionID, userID, algorithm, problem, variant string, config *api.DEConfig) {
	// Acquire worker slot
	e.workerPool <- struct{}{}
	defer func() { <-e.workerPool }()

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
			slog.Error("panic in execution",
				slog.String("execution_id", executionID),
				slog.Any("panic", r),
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
	paretoID, err := e.saveResults(ctx, userID, pareto, maxObjs)
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
	// Initialize completion counter for this execution
	counter := &atomic.Int32{}
	e.countersMu.Lock()
	e.completionCounters[executionID] = counter
	e.countersMu.Unlock()

	// Clean up counter when done
	defer func() {
		e.countersMu.Lock()
		delete(e.completionCounters, executionID)
		e.countersMu.Unlock()
	}()

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
	gde3Constants := gde3.Constants{
		DE: de.Constants{
			Executions:    int(config.Executions),
			Generations:   int(config.Generations),
			Dimensions:    int(config.DimensionsSize),
			ObjFuncAmount: int(config.ObjectivesSize),
		},
		CR: float64(config.GetGde3().Cr),
		F:  float64(config.GetGde3().F),
		P:  float64(config.GetGde3().P),
	}

	// Create progress callback that tracks completions
	progressCallback := func(generation int, totalGenerations int, paretoSize int, currentPareto []models.Vector) {
		// If this is the final generation, increment completion counter
		if generation == totalGenerations {
			counter.Add(1)
		}

		// Convert to API vectors (limit to avoid excessive data)
		const maxVectorsInProgress = 100
		apiVectors := make([]*api.Vector, 0, min(len(currentPareto), maxVectorsInProgress))
		for i := 0; i < len(currentPareto) && i < maxVectorsInProgress; i++ {
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
			TotalExecutions:     int32(config.Executions),
			PartialPareto:       apiVectors,
			UpdatedAt:           time.Now(),
		}

		if err := e.store.SaveProgress(ctx, progress); err != nil {
			slog.Warn("failed to save progress",
				slog.String("execution_id", executionID),
				slog.String("error", err.Error()),
			)
		}
	}

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

	mode, err := de.New(deConfig, de.WithAlgorithm(algorithm))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create DE mode: %w", err)
	}

	// Execute
	return mode.Execute(ctx)
}

func (e *Executor) saveResults(ctx context.Context, userID string, pareto []models.Vector, maxObjs [][]float64) (uint64, error) {
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
		Vectors:       apiVectors,
		MaxObjectives: storeMaxObjs,
		CreatedAt:     time.Now(),
	}

	if err := e.store.CreateParetoSet(ctx, paretoSet); err != nil {
		return 0, fmt.Errorf("failed to create pareto set: %w", err)
	}

	return paretoSet.ID, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
