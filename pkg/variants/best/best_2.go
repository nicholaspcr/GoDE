package best

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// best2
type best2 struct{}

// Best2 returns a DE/best/2 mutation strategy that uses the best individual
// and two difference vectors: best + F(r1 - r2) + F(r3 - r4).
func Best2() variants.Interface {
	return &best2{}
}

func (b *best2) Name() string {
	return "best2"
}

func (b *best2) Mutate(
	elems, rankZero []models.Vector,
	p variants.Parameters,
) (models.Vector, error) {
	// random element from rankZero
	bestIdx := p.Random.Intn(len(rankZero))
	// indices of the elements to be used in the mutation
	ind := make([]int, 5)
	ind[0] = p.CurrPos
	err := variants.GenerateIndices(1, len(elems), ind, p.Random)
	if err != nil {
		return models.Vector{}, variants.ErrInsufficientPopulation
	}

	// Validate vectors have non-nil elements
	if rankZero[bestIdx].Elements == nil || len(rankZero[bestIdx].Elements) != p.DIM {
		return models.Vector{}, variants.ErrInvalidVector
	}
	for _, idx := range []int{ind[1], ind[2], ind[3], ind[4]} {
		if elems[idx].Elements == nil || len(elems[idx].Elements) != p.DIM {
			return models.Vector{}, variants.ErrInvalidVector
		}
	}

	arr := make([]float64, p.DIM)
	for i := 0; i < p.DIM; i++ {
		arr[i] = rankZero[bestIdx].Elements[i] +
			p.F*(elems[ind[1]].Elements[i]-elems[ind[2]].Elements[i]) +
			p.F*(elems[ind[3]].Elements[i]-elems[ind[4]].Elements[i])
	}

	ret := models.Vector{
		Elements: arr,
	}
	return ret, nil
}
