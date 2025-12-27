// Package gde3 implements the GDE3 multi-objective Differential Evolution algorithm.
package gde3

import (
	"context"
	"log/slog"
	"math/rand"

	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/variants"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func init() {
	de.DefaultRegistry.Register("gde3", de.AlgorithmMetadata{
		Name:        "gde3",
		Description: "GDE3 - Generalized Differential Evolution 3rd version for multi-objective optimization",
	})
}

// gde3 type that contains the definition of the GDE3 algorithm.
type gde3 struct {
	problem           problems.Interface
	variant           variants.Interface
	initialPopulation models.Population
	populationParams  models.PopulationParams
	constants         Constants
	progressCallback  de.ProgressCallback
}

// Option is a functional option for configuring the GDE3 algorithm.
type Option func(*gde3)

// New creates a new GDE3 algorithm instance with the given configuration options.
// GDE3 is a multi-objective Differential Evolution algorithm that uses non-dominated
// sorting and crowding distance for selection.
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
	tracer := otel.Tracer("gde3")
	ctx, span := tracer.Start(ctx, "gde3.Execute",
		trace.WithAttributes(
			attribute.Int("population_size", g.populationParams.PopulationSize),
			attribute.Int("dimensions", g.populationParams.DimensionSize),
			attribute.Int("objectives", g.populationParams.ObjectivesSize),
			attribute.Int("generations", g.constants.DE.Generations),
			attribute.String("variant", g.variant.Name()),
			attribute.String("problem", g.problem.Name()),
		),
	)
	defer span.End()

	logger := slog.Default()
	random := rand.New(rand.NewSource(rand.Int63()))

	execNum := de.FromContextExecutionNumber(ctx)
	logger.Debug("Starting GDE3", slog.Int("execution", execNum))
	span.SetAttributes(attribute.Int("execution_number", execNum))

	population := g.initialPopulation.Copy()
	popuParams := g.populationParams

	maxObjs, err := g.initializePopulation(ctx, population)
	if err != nil {
		span.RecordError(err)
		return err
	}

	// Pre-allocate bestElems with estimated capacity based on generations
	// Each generation typically produces rank-zero elements up to population size
	estimatedCapacity := popuParams.PopulationSize * g.constants.DE.Generations
	bestElems := make([]models.Vector, 0, estimatedCapacity)

	for gen := range g.constants.DE.Generations {
		// Check for cancellation at the start of each generation
		if err := ctx.Err(); err != nil {
			span.RecordError(err)
			return err
		}

		logger.Debug("Running generation",
			slog.Int("execution_n", execNum),
			slog.Int("generation_n", gen),
		)

		newPopulation, bestInGen, err := g.runGeneration(ctx, population, random)
		if err != nil {
			span.RecordError(err)
			return err
		}
		population = newPopulation

		// NOTE: It probably would be a good idea to send the elements into the
		// channel directly instead of appending.
		bestElems = append(bestElems, bestInGen...)

		// Call progress callback if set
		if g.progressCallback != nil {
			g.progressCallback(gen+1, g.constants.DE.Generations, len(bestElems), bestElems)
		}
	}

	span.SetAttributes(attribute.Int("pareto_size", len(bestElems)))
	maxObjCh <- maxObjs
	paretoCh <- bestElems
	return nil
}

func (g *gde3) initializePopulation(ctx context.Context, population models.Population) ([]float64, error) {
	tracer := otel.Tracer("gde3")
	ctx, span := tracer.Start(ctx, "gde3.initializePopulation",
		trace.WithAttributes(
			attribute.Int("population_size", len(population)),
		),
	)
	defer span.End()

	maxObjs := make([]float64, g.populationParams.ObjectivesSize)
	for i := range population {
		// Check for cancellation between evaluations
		if err := ctx.Err(); err != nil {
			span.RecordError(err)
			return nil, err
		}
		if err := g.problem.Evaluate(&population[i], g.populationParams.ObjectivesSize); err != nil {
			span.RecordError(err)
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
	tracer := otel.Tracer("gde3")
	ctx, span := tracer.Start(ctx, "gde3.runGeneration",
		trace.WithAttributes(
			attribute.Int("population_size", len(population)),
		),
	)
	defer span.End()

	genRankZero, _ := de.FilterDominated(population)

	for i := range len(population) {
		trial, err := g.mutateAndCrossover(ctx, population, genRankZero, i, random)
		if err != nil {
			span.RecordError(err)
			return nil, nil, err
		}

		if err := g.problem.Evaluate(&trial, g.populationParams.ObjectivesSize); err != nil {
			span.RecordError(err)
			return nil, nil, err
		}

		population = g.selection(population, trial, i)
	}

	reducedPop, rankZero := de.ReduceByCrowdDistance(
		ctx, population, g.populationParams.PopulationSize,
	)
	span.SetAttributes(
		attribute.Int("rank_zero_size", len(rankZero)),
		attribute.Int("reduced_population_size", len(reducedPop)),
	)
	return reducedPop, rankZero, nil
}

func (g *gde3) mutateAndCrossover(
	ctx context.Context,
	population, genRankZero []models.Vector,
	currentIdx int,
	random *rand.Rand,
) (models.Vector, error) {
	tracer := otel.Tracer("gde3")
	_, span := tracer.Start(ctx, "gde3.mutateAndCrossover")
	defer span.End()

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
		span.RecordError(err)
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
