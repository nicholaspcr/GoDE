package models

import (
	"errors"

	"github.com/nicholaspcr/GoDE/pkg/api"
)

// Population is the set of elements in a DE generation.
type Population struct {
	Vectors []Vector
}

// ToPB converts the population to a protobuf message.
func (p *Population) ToPB() *api.Population {
	popu := &api.Population{
		Vectors: make([]*api.Vector, len(p.Vectors)),
	}
	for i, v := range p.Vectors {
		popu.Vectors[i] = v.ToPB()
	}
	return popu
}

// PopulationFromPB converts the population from a protobuf message.
func PopulationFromPB(popu *api.Population) Population {
	vectors := make([]Vector, len(popu.Vectors))
	for i, v := range popu.Vectors {
		vectors[i] = VectorFromPB(v)
	}
	return Population{
		Vectors: vectors,
	}
}

// Copy returns a copy of the population.
func (p *Population) Copy() Population {
	vectors := make([]Vector, len(p.Vectors))
	for i, v := range p.Vectors {
		vectors[i] = v.Copy()
	}
	return Population{
		Vectors: vectors,
	}
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
	}
	return Population{
		Vectors: vectors,
	}, nil
}
