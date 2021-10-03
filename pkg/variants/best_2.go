package variants

import (
	"errors"
	"math/rand"

	"github.com/nicholaspcr/gde3/pkg/models"
)

// best2 current_best + F(a-b) + F(c-d)
// a,b,c,d are random elements
var best2 = models.Variant{
	Fn: func(elems, rankZero models.Population, p models.VariantParams) (models.Vector, error) {
		// indices of the
		ind := make([]int, 5)
		ind[0] = p.CurrPos
		err := generateIndices(1, len(elems), ind)

		if err != nil {
			return models.Vector{}, errors.New(
				"insufficient size for the population, must me equal or greater than 4",
			)
		}

		arr := make([]float64, p.DIM)

		// random element from rankZero
		best := rankZero[rand.Int()%len(rankZero)]
		r1, r2, r3, r4 := elems[ind[1]], elems[ind[2]], elems[ind[3]], elems[ind[4]]
		for i := 0; i < p.DIM; i++ {
			arr[i] = best.X[i] + p.F*(r1.X[i]-r2.X[i]) + p.F*(r3.X[i]-r4.X[i])
		}

		ret := models.Vector{
			X: arr,
		}
		return ret, nil
	},
	VariantName: "best2",
}
