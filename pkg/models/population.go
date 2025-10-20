package models

import (
	"fmt"
	"math/rand"
)

type Population []Vector

// Copy returns a copy of the population.
func (p Population) Copy() Population {
	newP := make([]Vector, len(p))
	copy(newP, p)
	return newP
}

// PopulationParams is the set of parameters to generate a population.
type PopulationParams struct {
	DimensionSize  int
	PopulationSize int
	ObjectivesSize int
	FloorRange     []float64
	CeilRange      []float64
}

// GeneratePopulation generates a population with the given parameters.
func GeneratePopulation(params PopulationParams, random *rand.Rand) (Population, error) {
	vectors := make([]Vector, params.PopulationSize)
	if len(params.FloorRange) != params.DimensionSize ||
		len(params.CeilRange) != params.DimensionSize {
		return Population{}, fmt.Errorf(
			"floor range and ceil range must have the same size as the dimension size, got %d, %d and %d",
			len(params.FloorRange),
			len(params.CeilRange),
			params.DimensionSize,
		)
	}

	for i := 0; i < params.PopulationSize; i++ {
		vectors[i] = Vector{
			Elements:         make([]float64, params.DimensionSize),
			Objectives:       make([]float64, params.ObjectivesSize),
			CrowdingDistance: 0.0,
		}

		for j := 0; j < params.DimensionSize; j++ {
			vectors[i].Elements[j] = params.FloorRange[j] +
				(params.CeilRange[j]-params.FloorRange[j])*random.Float64()
		}

	}
	return vectors, nil
}
