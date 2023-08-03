package best

import (
	"errors"
	"math/rand"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// best1
type best1 struct{}

func Best1() variants.Interface {
	return &best1{}
}

func (b *best1) Name() string {
	return "best1"
}

func (b *best1) Mutate(
	elems, rankZero []models.Vector,
	p variants.Parameters,
) (models.Vector, error) {
	// random element from rankZero
	bestIdx := rand.Int() % len(rankZero)
	// indices of the elements to be used in the mutation
	index := make([]int, 3)
	index[0] = p.CurrPos
	err := variants.GenerateIndices(1, len(elems), index, p.Random)

	if err != nil {
		return models.Vector{}, errors.New(
			"insufficient size for the population, must me equal or greater than 4",
		)
	}

	arr := make([]float64, p.DIM)
	for i := 0; i < p.DIM; i++ {
		arr[i] = rankZero[bestIdx].Elements[i] +
			p.F*(elems[index[i]].Elements[i]-elems[index[2]].Elements[i])
	}

	ret := models.Vector{
		Elements: arr,
	}
	return ret, nil

}
