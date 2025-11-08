package gde3

import (
	"context"
	"log/slog"
	"math/rand"

	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// gde3 type that contains the definition of the GDE3 algorithm.
type gde3 struct {
	problem           problems.Interface
	variant           variants.Interface
	initialPopulation models.Population
	populationParams  models.PopulationParams
	constants         Constants
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
	paretoCh chan<- []models.Vector,
	maxObjCh chan<- []float64,
) error {
	logger := slog.Default()
	random := rand.New(rand.NewSource(rand.Int63()))

	execNum := de.FromContextExecutionNumber(ctx)
	logger.Debug("Starting GDE3", slog.Int("execution", execNum))

	population := g.initialPopulation.Copy()
	popuParams := g.populationParams

	maxObjs, err := g.initializePopulation(population)
	if err != nil {
		return err
	}

	bestElems := make([]models.Vector, 0, popuParams.PopulationSize)

	for gen := range g.constants.DE.Generations {
		logger.Debug("Running generation",
			slog.Int("execution_n", execNum),
			slog.Int("generation_n", gen),
		)

		newPopulation, bestInGen, err := g.runGeneration(ctx, population, random)
		if err != nil {
			return err
		}
		population = newPopulation

		// NOTE: It probably would be a good idea to send the elements into the
		// channel directly instead of appending.
		bestElems = append(bestElems, bestInGen...)
	}

	maxObjCh <- maxObjs
	paretoCh <- bestElems
	return nil
}

func (g *gde3) initializePopulation(population models.Population) ([]float64, error) {
	maxObjs := make([]float64, g.populationParams.ObjectivesSize)
	for i := range population {
		if err := g.problem.Evaluate(&population[i], g.populationParams.ObjectivesSize); err != nil {
			return nil, err
		}
		for j, obj := range population[i].Objectives {
			if obj > maxObjs[j] {
				maxObjs[j] = obj
			}
		}
	}
	return maxObjs, nil
}

func (g *gde3) runGeneration(
	ctx context.Context,
	population models.Population,
	random *rand.Rand,
) (models.Population, []models.Vector, error) {
	genRankZero, _ := de.FilterDominated(population)

	for i := range len(population) {
		trial, err := g.mutateAndCrossover(population, genRankZero, i, random)
		if err != nil {
			return nil, nil, err
		}

		if err := g.problem.Evaluate(&trial, g.populationParams.ObjectivesSize); err != nil {
			return nil, nil, err
		}

		population = g.selection(population, trial, i)
	}

	reducedPop, rankZero := de.ReduceByCrowdDistance(
		ctx, population, g.populationParams.PopulationSize,
	)
	return reducedPop, rankZero, nil
}

func (g *gde3) mutateAndCrossover(
	population, genRankZero []models.Vector,
	currentIdx int,
	random *rand.Rand,
) (models.Vector, error) {
	popuParams := g.populationParams

	vr, err := g.variant.Mutate(
		population,
		genRankZero,
		variants.Parameters{
			DIM:     popuParams.DimensionSize,
			F:       g.constants.F,
			CurrPos: currentIdx,
			P:       g.constants.P,
			Random:  random,
		},
	)
	if err != nil {
		return models.Vector{}, err
	}

	trial := population[currentIdx].Copy()

	currInd := random.Int() % popuParams.DimensionSize
	luckyIndex := random.Int() % popuParams.DimensionSize

	for range popuParams.DimensionSize {
		changeProb := random.Float64()
		if changeProb < g.constants.CR || currInd == luckyIndex {
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

	return trial, nil
}

func (g *gde3) selection(population []models.Vector, trial models.Vector, currentIdx int) models.Population {
	comp := de.DominanceTest(
		population[currentIdx].Objectives, trial.Objectives,
	)
	if comp == 1 {
		population[currentIdx] = trial
	} else if comp == 0 && len(population) <= 2*g.populationParams.PopulationSize {
		population = append(population, trial)
	}
	return population
}
