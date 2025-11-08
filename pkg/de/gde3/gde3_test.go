package gde3

import (
	"context"
	"math/rand"
	"testing"

	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems/multi"
	variantsrand "github.com/nicholaspcr/GoDE/pkg/variants/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestPopulation(popSize, dimSize, objSize int) (models.Population, models.PopulationParams) {
	params := models.PopulationParams{
		PopulationSize: popSize,
		DimensionSize:  dimSize,
		ObjectivesSize: objSize,
		FloorRange:     make([]float64, dimSize),
		CeilRange:      make([]float64, dimSize),
	}
	for i := range params.CeilRange {
		params.CeilRange[i] = 1.0
	}

	random := rand.New(rand.NewSource(1))
	population, _ := models.GeneratePopulation(params, random)
	return population, params
}

func TestGDE3_New(t *testing.T) {
	t.Run("create GDE3 instance with options", func(t *testing.T) {
		problem := multi.Zdt1()
		variant := variantsrand.Rand1()
		constants := Constants{
			CR: 0.5,
			F:  0.5,
			P:  0.1,
			DE: de.Constants{
				Generations: 10,
			},
		}
		population, params := createTestPopulation(10, 5, 2)

		algorithm := New(
			WithProblem(problem),
			WithVariant(variant),
			WithConstants(constants),
			WithInitialPopulation(population),
			WithPopulationParams(params),
		)

		assert.NotNil(t, algorithm)
	})

	t.Run("create GDE3 without options", func(t *testing.T) {
		algorithm := New()
		assert.NotNil(t, algorithm)
	})
}

func TestGDE3_Execute(t *testing.T) {
	t.Run("execute with small population and generations", func(t *testing.T) {
		problem := multi.Zdt1()
		variant := variantsrand.Rand1()
		constants := Constants{
			CR: 0.9,
			F:  0.5,
			P:  0.1,
			DE: de.Constants{
				Generations: 5,
			},
		}

		population, params := createTestPopulation(10, 5, 2)

		algorithm := New(
			WithProblem(problem),
			WithVariant(variant),
			WithConstants(constants),
			WithInitialPopulation(population),
			WithPopulationParams(params),
		)

		ctx := de.WithContextExecutionNumber(context.Background(), 1)
		paretoCh := make(chan []models.Vector, 1)
		maxObjCh := make(chan []float64, 1)

		go func() {
			err := algorithm.Execute(ctx, paretoCh, maxObjCh)
			assert.NoError(t, err)
		}()

		// Receive results
		pareto := <-paretoCh
		maxObjs := <-maxObjCh

		assert.NotEmpty(t, pareto, "should have pareto front")
		assert.Len(t, maxObjs, 2, "should have 2 objectives")
		assert.Greater(t, len(pareto), 0, "pareto front should not be empty")

		// Verify all vectors have objectives
		for _, vec := range pareto {
			assert.Len(t, vec.Objectives, 2)
			assert.Len(t, vec.Elements, 5)
		}
	})

	t.Run("execute with different CR values", func(t *testing.T) {
		problem := multi.Zdt1()
		variant := variantsrand.Rand1()

		testCases := []struct {
			name string
			cr   float64
		}{
			{"low CR", 0.1},
			{"medium CR", 0.5},
			{"high CR", 0.9},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				constants := Constants{
					CR: tc.cr,
					F:  0.5,
					P:  0.1,
					DE: de.Constants{Generations: 3},
				}

				population, params := createTestPopulation(5, 3, 2)

				algorithm := New(
					WithProblem(problem),
					WithVariant(variant),
					WithConstants(constants),
					WithInitialPopulation(population),
					WithPopulationParams(params),
				)

				ctx := de.WithContextExecutionNumber(context.Background(), 1)
				paretoCh := make(chan []models.Vector, 1)
				maxObjCh := make(chan []float64, 1)

				go func() {
					err := algorithm.Execute(ctx, paretoCh, maxObjCh)
					assert.NoError(t, err)
				}()

				pareto := <-paretoCh
				maxObjs := <-maxObjCh

				assert.NotEmpty(t, pareto)
				assert.Len(t, maxObjs, 2)
			})
		}
	})
}

func TestGDE3_InitializePopulation(t *testing.T) {
	t.Run("initialize population evaluates all individuals", func(t *testing.T) {
		problem := multi.Zdt1()
		population, params := createTestPopulation(5, 3, 2)

		algorithm := New(
			WithProblem(problem),
			WithPopulationParams(params),
		).(*gde3)

		maxObjs, err := algorithm.initializePopulation(context.Background(), population)
		require.NoError(t, err)

		// Verify max objectives
		assert.Len(t, maxObjs, 2)
		assert.Greater(t, maxObjs[0], 0.0)
		assert.Greater(t, maxObjs[1], 0.0)

		// Verify all individuals have objectives
		for _, ind := range population {
			assert.Len(t, ind.Objectives, 2)
			assert.NotEmpty(t, ind.Objectives)
		}
	})

	t.Run("max objectives are correctly computed", func(t *testing.T) {
		problem := multi.Zdt1()
		population, params := createTestPopulation(10, 5, 2)

		algorithm := New(
			WithProblem(problem),
			WithPopulationParams(params),
		).(*gde3)

		maxObjs, err := algorithm.initializePopulation(context.Background(), population)
		require.NoError(t, err)

		// Max objectives should be >= all individual objectives
		for _, ind := range population {
			for i, obj := range ind.Objectives {
				assert.LessOrEqual(t, obj, maxObjs[i],
					"max objective should be >= individual objective")
			}
		}
	})
}

func TestGDE3_RunGeneration(t *testing.T) {
	t.Run("run generation produces valid population", func(t *testing.T) {
		problem := multi.Zdt1()
		variant := variantsrand.Rand1()
		constants := Constants{
			CR: 0.9,
			F:  0.5,
			P:  0.1,
			DE: de.Constants{Generations: 1},
		}

		population, params := createTestPopulation(10, 5, 2)

		algorithm := New(
			WithProblem(problem),
			WithVariant(variant),
			WithConstants(constants),
			WithInitialPopulation(population),
			WithPopulationParams(params),
		).(*gde3)

		// Initialize first
		_, err := algorithm.initializePopulation(context.Background(), population)
		require.NoError(t, err)

		random := rand.New(rand.NewSource(1))
		ctx := context.Background()
		newPop, rankZero, err := algorithm.runGeneration(ctx, population, random)
		require.NoError(t, err)

		assert.NotNil(t, newPop)
		assert.NotNil(t, rankZero)
		assert.LessOrEqual(t, len(newPop), params.PopulationSize*2)
		assert.Greater(t, len(rankZero), 0, "should have at least one non-dominated solution")
	})
}

func TestGDE3_Selection(t *testing.T) {
	problem := multi.Zdt1()
	population, params := createTestPopulation(5, 3, 2)

	algorithm := New(
		WithProblem(problem),
		WithPopulationParams(params),
	).(*gde3)

	t.Run("trial dominates current, trial replaces current", func(t *testing.T) {
		testPop := population.Copy()
		// Initialize
		_, _ = algorithm.initializePopulation(context.Background(), testPop)

		// Create a trial that is better
		trial := testPop[0].Copy()
		for i := range trial.Objectives {
			trial.Objectives[i] = trial.Objectives[i] * 0.5
		}

		originalObj := testPop[0].Objectives[0]
		newPop := algorithm.selection(testPop, trial, 0)

		// Trial should have replaced the current or population grew
		assert.True(t, newPop[0].Objectives[0] != originalObj || len(newPop) > len(testPop))
	})

	t.Run("population can grow when non-dominating", func(t *testing.T) {
		testPop := population.Copy()
		_, _ = algorithm.initializePopulation(context.Background(), testPop)

		trial := testPop[0].Copy()
		trial.Objectives[0] = trial.Objectives[0] * 0.8
		trial.Objectives[1] = trial.Objectives[1] * 1.2

		originalLen := len(testPop)
		newPop := algorithm.selection(testPop, trial, 0)

		// Population may grow if there's space
		assert.GreaterOrEqual(t, len(newPop), originalLen)
	})
}

func TestGDE3_EdgeCases(t *testing.T) {
	t.Run("single generation", func(t *testing.T) {
		problem := multi.Zdt1()
		variant := variantsrand.Rand1()
		constants := Constants{
			CR: 0.9,
			F:  0.5,
			P:  0.1,
			DE: de.Constants{Generations: 1},
		}

		population, params := createTestPopulation(5, 3, 2)

		algorithm := New(
			WithProblem(problem),
			WithVariant(variant),
			WithConstants(constants),
			WithInitialPopulation(population),
			WithPopulationParams(params),
		)

		ctx := de.WithContextExecutionNumber(context.Background(), 1)
		paretoCh := make(chan []models.Vector, 1)
		maxObjCh := make(chan []float64, 1)

		go func() {
			err := algorithm.Execute(ctx, paretoCh, maxObjCh)
			assert.NoError(t, err)
		}()

		pareto := <-paretoCh
		maxObjs := <-maxObjCh

		assert.NotEmpty(t, pareto)
		assert.Len(t, maxObjs, 2)
	})

	t.Run("minimal population of 5", func(t *testing.T) {
		problem := multi.Zdt1()
		variant := variantsrand.Rand1()
		constants := Constants{
			CR: 0.9,
			F:  0.5,
			P:  0.1,
			DE: de.Constants{Generations: 2},
		}

		population, params := createTestPopulation(5, 2, 2)

		algorithm := New(
			WithProblem(problem),
			WithVariant(variant),
			WithConstants(constants),
			WithInitialPopulation(population),
			WithPopulationParams(params),
		)

		ctx := de.WithContextExecutionNumber(context.Background(), 1)
		paretoCh := make(chan []models.Vector, 1)
		maxObjCh := make(chan []float64, 1)

		go func() {
			err := algorithm.Execute(ctx, paretoCh, maxObjCh)
			assert.NoError(t, err)
		}()

		pareto := <-paretoCh
		maxObjs := <-maxObjCh

		assert.NotEmpty(t, pareto)
		assert.Len(t, maxObjs, 2)
	})
}
