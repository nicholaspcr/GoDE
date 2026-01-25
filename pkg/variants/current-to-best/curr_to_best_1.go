// Package currenttobest implements DE/current-to-best mutation strategies that blend current and best individuals.
package currenttobest

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// currToBest1
type currToBest1 struct{}

// CurrToBest1 returns a DE/current-to-best/1 mutation strategy that blends the current
// individual with the best: current + F(best - r1) + F(r2 - r3).
func CurrToBest1() variants.Interface {
	return &currToBest1{}
}

func (c *currToBest1) Name() string {
	return "currToBest1"
}

func (c *currToBest1) Mutate(
	elems, rankZero []models.Vector,
	p variants.Parameters,
) (models.Vector, error) {
	if len(rankZero) == 0 {
		return models.Vector{}, variants.ErrEmptyRankZero
	}
	// random element from rankZero
	bestIdx := p.Random.Intn(len(rankZero))
	// indices of the elements to be used in the mutation
	ind := make([]int, 4)
	ind[0] = p.CurrPos
	err := variants.GenerateIndices(1, len(elems), ind, p.Random)
	if err != nil {
		return models.Vector{}, variants.ErrInsufficientPopulation
	}

	// Validate elems vectors have non-nil elements
	if err := variants.ValidateVectors(elems, []int{p.CurrPos, ind[1], ind[2], ind[3]}, p.DIM); err != nil {
		return models.Vector{}, err
	}
	// Validate rankZero best vector
	if err := variants.ValidateVectors(rankZero, []int{bestIdx}, p.DIM); err != nil {
		return models.Vector{}, err
	}

	arr := make([]float64, p.DIM)
	for i := 0; i < p.DIM; i++ {
		arr[i] = elems[p.CurrPos].Elements[i] +
			p.F*(rankZero[bestIdx].Elements[i]-elems[ind[1]].Elements[i]) +
			p.F*(elems[ind[2]].Elements[i]-elems[ind[3]].Elements[i])
	}

	ret := models.Vector{
		Elements: arr,
	}
	return ret, nil
}
