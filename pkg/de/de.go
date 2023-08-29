package de

import (
	"context"
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
	constants Constants
}

// New Mode iterface based on the configuration options given.
func New(opts ...ModeOptions) *de {
	m := &de{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (mode *de) Execute(ctx context.Context) error {
	logger := slog.Default()
	paretoCh := make(chan []models.Vector, 100)
	maxObjsCh := make(chan<- []float64, 100)
	wgExecs := &sync.WaitGroup{}

	// Runs algorithm for Executions amount of times.
	for i := 1; i <= mode.constants.Executions; i++ {
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
				logger.Error("Unexpected error while executing the algorith",
					slog.Int("Execution", idx),
					slog.String("error", err.Error()),
				)
			}
		}(i)
	}

	wgExecs.Wait()
	close(paretoCh)

	// TODO: Filter the rank zero of each algorithm in parallel.
	// This affects the behavior of the algorithm, since the ranking between
	// different rank zeros in succession can lead to getting better vectors
	// than the process of filtering the rank zeros in parallel.
	//
	// Before doing this I should do a benchmark to see if it is a permissable.

	// TODO: Taking into consideration the benchmark mentioned in the TODO
	// above, what should be tested is the degree of deviation from the results
	// known from the test functions. Write both results to different files and
	// define the difference between them.
	now := time.Now()
	finalPareto := filterPareto(ctx, paretoCh)
	logger.Info("Filtering Pareto", slog.Duration("time", time.Since(now)))
	_ = finalPareto

	//now := time.Now()
	//finalParetoParallel := filterParetoParallel(ctx, paretoCh)
	//logger.Info("Filtering Pareto Parallel",
	//	slog.Duration("time", time.Since(now)),
	//)
	//_ = finalParetoParallel

	// TODO: Write the ranked pareto into its own separate section, make it a
	// separate table on the database.

	//	// result of the ranked pareto
	//	f, err := os.Create(
	//		homePath + multiExecutionsPath + "/rankedPareto.csv",
	//	)
	//	// creates writer and writes the elements objs
	//	w := writer.NewWriter(f)
	//	w.Comma = ';'
	//	if err := w.WriteHeader(m.constants.ObjFuncAmount); err != nil {
	//		panic(err)
	//	}
	//	if err := w.ElementsObjs(rankedPareto); err != nil {
	//		panic(err)
	//	}

	// TODO: The biggest objectives values are to be a part of the pareto table.

	//	// getting biggest objs values
	//	mo := make([]float64, m.constants.Dimensions)
	//	for arr := range maximumObjs {
	//		for i, obj := range arr {
	//			if obj > mo[i] {
	//				mo[i] = obj
	//			}
	//		}
	//	}
	return nil
}

func filterPareto(
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

		// TODO: Make it configurable.
		// Limits the amounts of dots to 1k.
		if len(finalPareto) > 1000 {
			finalPareto = finalPareto[:1000]
		}
	}
	return finalPareto
}

func filterParetoParallel(
	ctx context.Context, pareto chan []models.Vector,
) []models.Vector {
	wgRank := &sync.WaitGroup{}
	filterCh := make(chan []models.Vector, 100)
	finalPareto := make([]models.Vector, 0, 2000)
	for v := range pareto {
		v := v
		wgRank.Add(1)
		go func(v []models.Vector) {
			defer wgRank.Done()
			_, v = ReduceByCrowdDistance(ctx, v, len(v))
			filterCh <- v
		}(v)
	}

	wgRank.Wait()
	close(filterCh)

	for v := range filterCh {
		finalPareto = append(
			finalPareto,
			v...,
		)
	}

	// Final filtering of the ranked pareto.
	_, finalPareto = ReduceByCrowdDistance(ctx, finalPareto, len(finalPareto))
	return finalPareto
}
