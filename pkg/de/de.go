package de

import (
	"context"
	"errors"
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
	algorithm Algorithm
	config    Config
	constants Constants
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
	paretoCh := make(chan []models.Vector, mode.config.ParetoChannelLimiter)
	maxObjsCh := make(chan []float64, mode.config.MaxChannelLimiter)
	wgExecs := &sync.WaitGroup{}

	// Collect max objectives from all executions
	allMaxObjs := make([][]float64, 0, mode.constants.Executions)
	var maxObjsMu sync.Mutex

	// Goroutine to consume max objectives (prevents deadlock)
	wgMaxObjs := &sync.WaitGroup{}
	wgMaxObjs.Add(1)
	go func() {
		defer wgMaxObjs.Done()
		for maxObjs := range maxObjsCh {
			maxObjsMu.Lock()
			allMaxObjs = append(allMaxObjs, maxObjs)
			maxObjsMu.Unlock()
		}
	}()

	// Runs algorithm for Executions amount of times.
	for i := range mode.constants.Executions {
		wgExecs.Add(1)
		// Initialize worker responsible for DE execution.
		go func(idx int) {
			defer wgExecs.Done()
			// running the algorithm execution.
			if err := mode.algorithm.Execute(
				WithContextExecutionNumber(ctx, idx),
				paretoCh,
				maxObjsCh,
			); err != nil {
				slog.Error("Unexpected error while executing the algorithm",
					slog.Int("Execution", idx),
					slog.String("error", err.Error()),
				)
			}
		}(i)
	}

	wgExecs.Wait()
	close(paretoCh)
	close(maxObjsCh)
	wgMaxObjs.Wait()

	now := time.Now()
	finalPareto := mode.filterPareto(ctx, paretoCh)
	slog.Info("Filtering Pareto", slog.Duration("time", time.Since(now)))

	return finalPareto, allMaxObjs, nil
}

func (mode *de) filterPareto(
	ctx context.Context, pareto chan []models.Vector,
) []models.Vector {
	finalPareto := make([]models.Vector, 0, 2000)
	for v := range pareto {
		finalPareto = append(
			finalPareto,
			v...,
		)
		// gets non dominated and filters by crowdingDistance
		_, finalPareto = ReduceByCrowdDistance(
			ctx, finalPareto, len(finalPareto),
		)

		if len(finalPareto) > mode.config.ResultLimiter {
			finalPareto = finalPareto[:mode.config.ResultLimiter]
		}
	}
	return finalPareto
}
