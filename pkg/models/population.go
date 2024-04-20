package models

import (
	"errors"
	"math/rand"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
)

type Population []Vector

// ToPB converts the population to a protobuf message.
func (p Population) ToPB() *api.Population {
	popu := &api.Population{
		Vectors: make([]*api.Vector, len(p)),
	}
	for i, v := range p {
		popu.Vectors[i] = v.ToPB()
	}
	return popu
}

// PopulationFromPB converts the population from a protobuf message.
func PopulationFromPB(pop *api.Population) Population {
	p := make([]Vector, len(pop.Vectors))
	for i, v := range pop.Vectors {
		p[i] = VectorFromPB(v)
	}
	return p
}

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
func GeneratePopulation(params PopulationParams) (Population, error) {
	vectors := make([]Vector, params.PopulationSize)
	if len(params.FloorRange) != params.DimensionSize ||
		len(params.CeilRange) != params.DimensionSize {
		return Population{}, errors.New("float range and ceil range must have the same size as the dimension size")
	}

	for i := 0; i < params.PopulationSize; i++ {
		vectors[i] = Vector{
			Elements:         make([]float64, params.DimensionSize),
			Objectives:       make([]float64, params.ObjectivesSize),
			CrowdingDistance: 0.0,
		}

		for j := 0; j < params.DimensionSize; j++ {
			vectors[i].Elements[j] = params.FloorRange[j] +
				(params.CeilRange[j]-params.FloorRange[j])*rand.Float64()
		}

	}
	return vectors, nil
}
