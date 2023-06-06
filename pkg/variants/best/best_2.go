package best

import (
	"errors"
	"math/rand"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// best2
type best2 struct{}

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
	bestIdx := rand.Int() % len(rankZero)
	// indices of the elements to be used in the mutation
	ind := make([]int, 5)
	ind[0] = p.CurrPos
	err := variants.GenerateIndices(1, len(elems), ind)

	if err != nil {
		return models.Vector{}, errors.New(
			"insufficient size for the population, must me equal or greater than 4",
		)
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
