package de

import (
	"context"
	"sync"

	"github.com/nicholaspcr/GoDE/internal/log"
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
	logger := log.FromContext(ctx)
	logger.Debug("Starting execution")

	rankedChan := make(chan []models.Vector, mode.constants.Executions)

	pareto := make(chan<- []models.Vector, 10)
	maxObjs := make(chan<- []float64, 10)

	// TODO: Change this to be just a pipeline pattern, that way there can be a
	// goroutine in the background that would write the last values.
	wg := &sync.WaitGroup{}

	// Runs algorithm for Executions amount of times.
	for i := 0; i < mode.constants.Executions; i++ {
		wg.Add(1)

		// Initialize worker responsible for DE execution.
		go func() {
			// cleaning concurrent queue
			defer func() {
				wg.Done()
			}()
			// running the algorithm execution.
			mode.algorithm.Execute(
				ctx,
				pareto,
				maxObjs,
			)
		}()
	}

	// closer
	go func() {
		wg.Wait()
		close(rankedChan)
	}()

	// gets data from the pareto created by rank[0] of each gen
	var rankedPareto []models.Vector
	for v := range rankedChan {
		rankedPareto = append(
			rankedPareto,
			v...,
		)

		// gets non dominated and filters by crowdingDistance
		_, rankedPareto = ReduceByCrowdDistance(
			rankedPareto,
			len(rankedPareto),
		)

		// TODO: Make it configurable.
		// Limits the amounts of dots to 1k.
		if len(rankedPareto) > 1000 {
			rankedPareto = rankedPareto[:1000]
		}
	}

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

	logger.Debug("Finished execution")
	return nil
}
