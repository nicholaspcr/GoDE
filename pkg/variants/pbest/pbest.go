// Package pbest implements DE/pbest mutation strategies using top-performing individuals.
package pbest

import (
	"math"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// pbest
type pbest struct{}

// Pbest returns a DE/pbest mutation strategy that uses a randomly selected individual
// from the top P% of the population: current + F(pbest - current) + F(r1 - r2).
func Pbest() variants.Interface {
	return &pbest{}
}

func (p *pbest) Name() string {
	return "pbest"
}

func (p *pbest) Mutate(
	elems, rankZero []models.Vector,
	params variants.Parameters,
) (models.Vector, error) {
	if len(rankZero) == 0 {
		return models.Vector{}, variants.ErrEmptyRankZero
	}
	ind := make([]int, 3)
	ind[0] = params.CurrPos

	err := variants.GenerateIndices(1, len(elems), ind, params.Random)
	if err != nil {
		return models.Vector{}, variants.ErrInsufficientPopulation
	}

	indexLimit := int(math.Ceil(float64(len(rankZero)) * params.P))
	bestIndex := params.Random.Intn(indexLimit)

	// Validate elems vectors have non-nil elements
	for _, idx := range []int{params.CurrPos, ind[1], ind[2]} {
		if elems[idx].Elements == nil || len(elems[idx].Elements) != params.DIM {
			return models.Vector{}, variants.ErrInvalidVector
		}
	}
	// Validate rankZero best vector
	if rankZero[bestIndex].Elements == nil || len(rankZero[bestIndex].Elements) != params.DIM {
		return models.Vector{}, variants.ErrInvalidVector
	}

	arr := make([]float64, params.DIM)
	for i := 0; i < params.DIM; i++ {
		arr[i] = elems[params.CurrPos].Elements[i] +
			params.F*(rankZero[bestIndex].Elements[i]-elems[params.CurrPos].Elements[i]) +
			params.F*(elems[ind[1]].Elements[i]-elems[ind[2]].Elements[i])
	}

	ret := models.Vector{
		Elements: arr,
	}
	return ret, nil
}
