package de

import (
	"context"
	"sync"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// DE contains the necessary methods to setup and execute a Differential
// Evolutionary algorithm.
type de struct {
	constants  Constants
	problem    problems.Interface
	variant    variants.Interface
	population models.Population
	store      Store
	algorithm  Algorithm
}

// New Mode iterface based on the configuration options given.
func New(opts ...ModeOptions) *de {
	m := &de{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// TODO: Make a separate semaphore, responsible for handling how many parallel
// executions the server can take.

func (m *de) Execute(
	ctx context.Context,
	pareto chan<- models.Population,
	maxObjs chan<- []float64,
) error {
	rankedChan := make(chan []models.Vector, m.constants.Executions)

	// TODO: Change this to be just a pipeline pattern, that way there can be a
	// goroutine in the background that would write the last values.
	wg := &sync.WaitGroup{}

	// TODO: generate first population.
	// initialPopulation is the the population in which every execution will start with.
	var initialPopulation models.Population
	GeneratePopulation(&initialPopulation)

	// Runs algorithm for Executions amount of times.
	for i := 0; i < m.constants.Executions; i++ {
		population := initialPopulation.Copy()
		wg.Add(1)

		// Initialize worker responsible for DE execution.
		go func() {
			// cleaning concurrent queue
			defer func() {
				wg.Done()
			}()
			// running one execution of the GDE3
			m.algorithm.Execute(
				ctx,
				population,
				m.problem,
				m.variant,
				m.store,
				rankedChan,
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
	//
	//	// getting biggest objs values
	//	mo := make([]float64, m.constants.Dimensions)
	//	for arr := range maximumObjs {
	//		for i, obj := range arr {
	//			if obj > mo[i] {
	//				mo[i] = obj
	//			}
	//		}
	//	}
	//	fmt.Println(
	//		"maximum objective values found",
	//	)
	//	fmt.Println(maxObjs)
	//
	//	// sends the values to the channel
	//	maxObjs <- mo
	//
	return nil
}
