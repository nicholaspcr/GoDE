package variants

import (
	"errors"
	"math"
	"math/rand"

	"github.com/nicholaspcr/gde3/pkg/problems/models"
)

// PBest implementation
var pbest VariantFn = VariantFn{
	Fn: func(elems, rankZero models.Elements, p Params) (models.Elem, error) {
		ind := make([]int, 3)
		ind[0] = p.CurrPos

		err := generateIndices(1, len(elems), ind)
		if err != nil {
			return models.Elem{}, errors.New(
				"insufficient size for the population, must me equal or greater than 5",
			)
		}

		indexLimit := int(math.Ceil(float64(len(rankZero)) * p.P))
		bestIndex := rand.Int() % indexLimit

		arr := make([]float64, p.DIM)

		r1, r2 := elems[ind[1]], elems[ind[2]]
		curr := elems[p.CurrPos]
		best := rankZero[bestIndex]

		for i := 0; i < p.DIM; i++ {
			arr[i] = curr.X[i] + p.F*(best.X[i]-curr.X[i]) + p.F*(r1.X[i]-r2.X[i])
		}

		ret := models.Elem{
			X: arr,
		}
		return ret, nil
	},
	Name: "pbest",
}
