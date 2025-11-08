package currenttobest

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// currToBest1
type currToBest1 struct{}

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
	// random element from rankZero
	bestIdx := p.Random.Intn(len(rankZero))
	// indices of the elements to be used in the mutation
	ind := make([]int, 4)
	ind[0] = p.CurrPos
	err := variants.GenerateIndices(1, len(elems), ind, p.Random)
	if err != nil {
		return models.Vector{}, variants.ErrInsufficientPopulation
	}

	// Validate vectors have non-nil elements
	for _, idx := range []int{p.CurrPos, bestIdx, ind[1], ind[2], ind[3]} {
		if elems[idx].Elements == nil || len(elems[idx].Elements) != p.DIM {
			return models.Vector{}, variants.ErrInvalidVector
		}
	}

	arr := make([]float64, p.DIM)
	for i := 0; i < p.DIM; i++ {
		arr[i] = elems[p.CurrPos].Elements[i] +
			p.F*(elems[bestIdx].Elements[i]-elems[ind[1]].Elements[i]) +
			p.F*(elems[ind[2]].Elements[i]-elems[ind[3]].Elements[i])
	}

	ret := models.Vector{
		Elements: arr,
	}
	return ret, nil
}
