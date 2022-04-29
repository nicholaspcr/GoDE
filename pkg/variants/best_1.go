package variants

import (
	"errors"
	"math/rand"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

// best1
type best1 struct{}

func Best1() models.Variant {
	return &best1{}
}

func (b *best1) Name() string {
	return "best1"
}

func (b *best1) Mutate(
	elems, rankZero models.Population,
	p models.VariantParams,
) (models.Vector, error) {
	index := make([]int, 3)
	index[0] = p.CurrPos
	err := generateIndices(1, len(elems), index)

	if err != nil {
		return models.Vector{}, errors.New(
			"insufficient size for the population, must me equal or greater than 4",
		)
	}

	arr := make([]float64, p.DIM)

	best := rankZero[rand.Int()%len(rankZero)]
	r1, r2 := elems[index[1]], elems[index[2]]
	for i := 0; i < p.DIM; i++ {
		arr[i] = best.X[i] + p.F*(r1.X[i]-r2.X[i])
	}

	ret := models.Vector{
		X: arr,
	}
	return ret, nil

}
