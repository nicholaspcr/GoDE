package variants

import (
	"errors"
	"math/rand"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

// currToBest1
type currToBest1 struct{}

func CurrToBest1() models.Variant {
	return &currToBest1{}
}

func (c *currToBest1) Name() string {
	return "currToBest1"
}

func (c *currToBest1) Mutate(
	elems, rankZero models.Population,
	p models.VariantParams,
) (models.Vector, error) {

	ind := make([]int, 4)
	ind[0] = p.CurrPos
	err := generateIndices(1, len(elems), ind)

	if err != nil {
		return models.Vector{}, errors.New(
			"insufficient size for the population, must me equal or greater than 5",
		)
	}

	arr := make([]float64, p.DIM)

	r1, r2, r3 := elems[ind[1]], elems[ind[2]], elems[ind[3]]
	curr := elems[p.CurrPos]
	best := rankZero[rand.Int()%len(rankZero)]

	for i := 0; i < p.DIM; i++ {
		arr[i] = curr.X[i] + p.F*(best.X[i]-r1.X[i]) + p.F*(r2.X[i]-r3.X[i])
	}

	ret := models.Vector{
		X: arr,
	}
	return ret, nil
}
