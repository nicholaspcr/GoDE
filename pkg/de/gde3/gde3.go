package gde3

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// gde3 type that contains the definition of the GDE3 algorithm.
type gde3 struct {
	initialPopulation models.Population
	populationParams  models.PopulationParams
	problem           problems.Interface
	variant           variants.Interface
	contants          de.Constants
	store             store.Store
}

type Option func(*gde3) *gde3

// GDE3 Returns an instance of an object that implements the GDE3
// algorithm. It is compliant with the Mode
func New(opts ...Option) de.Algorithm {
	d := &gde3{}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// Execute is responsible for receiving the standard parameters defined in the
// Mode and executing the gde3 algorithm
func (g *gde3) Execute(
	ctx context.Context,
	pareto chan<- []models.Vector,
	maxObjectives chan<- []float64,
) error {
	logger := slog.Default()
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	execNum := de.FromContextExecutionNumber(ctx)
	logger.Debug("Starting GDE3", slog.Int("execution", execNum))

	population := g.initialPopulation.Copy()
	popuParams := g.populationParams
	dimSize := popuParams.DimensionSize
	objFuncAmount := g.populationParams.ObjectivesSize
	maxObjs := make([]float64, dimSize)

	// calculates the objectives of the initial population
	for i := range population.Vectors {
		err := g.problem.Evaluate(&population.Vectors[i], objFuncAmount)
		if err != nil {
			return err
		}
		for j, obj := range population.Vectors[i].Objectives {
			if obj > maxObjs[j] {
				maxObjs[j] = obj
			}
		}
	}

	//// writes the header in this execution's file
	//if err := store.Header(); err != nil {
	//	// TODO: Add header contents to methods
	//	return err
	//}
	// writes the initial generation

	// TODO: Update how the population is written
	//if err := store.Population(population); err != nil {
	//	panic(err)
	//}

	// stores the rank[0] of each generation
	bestElems := make([]models.Vector, 0, popuParams.DimensionSize)

	var genRankZero []models.Vector
	var bestInGen []models.Vector

	for gen := 0; gen < g.contants.Generations; gen++ {
		logger.Debug("Running generation",
			slog.Int("execution_n", execNum),
			slog.Int("generation_n", gen),
		)
		// gets non dominated of the current population
		genRankZero, _ = de.FilterDominated(population.Vectors)

		for i := 0; i < len(population.Vectors); i++ {
			// generates the mutatated vector
			vr, err := g.variant.Mutate(
				population.Vectors,
				genRankZero,
				variants.Parameters{
					DIM:     popuParams.DimensionSize,
					F:       g.contants.F,
					CurrPos: i,
					P:       g.contants.P,

					Random: random,
				})
			if err != nil {
				return err
			}

			// trial element
			trial := population.Vectors[i].Copy()

			// CROSS OVER
			currInd := random.Int() % popuParams.DimensionSize
			luckyIndex := random.Int() % popuParams.DimensionSize

			for j := 0; j < popuParams.DimensionSize; j++ {
				changeProb := random.Float64()
				if changeProb < g.contants.CR || currInd == luckyIndex {
					trial.Elements[currInd] = vr.Elements[currInd]
				}

				if trial.Elements[currInd] < popuParams.FloorRange[currInd] {
					trial.Elements[currInd] = popuParams.FloorRange[currInd]
				}
				if trial.Elements[currInd] > popuParams.CeilRange[currInd] {
					trial.Elements[currInd] = popuParams.CeilRange[currInd]
				}
				currInd = (currInd + 1) % popuParams.DimensionSize
			}

			if err := g.problem.Evaluate(&trial, objFuncAmount); err != nil {
				return err
			}

			// SELECTION
			comp := de.DominanceTest(
				population.Vectors[i].Objectives, trial.Objectives,
			)
			if comp == 1 {
				population.Vectors[i] = trial
			} else if comp == 0 && len(population.Vectors) <= 2*popuParams.DimensionSize {
				population.Vectors = append(population.Vectors, trial)
			}
		}

		population.Vectors, bestInGen = de.ReduceByCrowdDistance(
			ctx, population.Vectors, popuParams.DimensionSize,
		)
		bestElems = append(bestElems, bestInGen...)

		//// writes the objectives of the population
		//if err := store.Population(population); err != nil {
		//	return err
		//}

		// TODO: Update how the population is written
		//// checks for the biggest objective
		//for _, vector := range population.Vectors {
		//	for j, obj := range vector.Objs {
		//		if obj > maxObjs[j] {
		//			maxObjs[j] = obj
		//		}
		//	}
		//}
	}

	//logger.Debug("bestElems: ", bestElems)
	// sending via channel the data
	// maximumObjs <- maxObjs
	pareto <- bestElems
	return nil
}
