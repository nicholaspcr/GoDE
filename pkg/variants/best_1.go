package variants

import (
	"errors"
	"math/rand"

	"github.com/nicholaspcr/gde3/pkg/problems/models"
)

// best1  current_best + F(a-b)
// a,b are random elements
var best1 VariantFn = VariantFn{
	Fn: func(elems, rankZero models.Elements, p Params) (models.Elem, error) {
		index := make([]int, 3)
		index[0] = p.CurrPos
		err := generateIndices(1, len(elems), index)

		if err != nil {
			return models.Elem{}, errors.New(
				"insufficient size for the population, must me equal or greater than 4",
			)
		}

		arr := make([]float64, p.DIM)

		best := rankZero[rand.Int()%len(rankZero)]
		r1, r2 := elems[index[1]], elems[index[2]]
		for i := 0; i < p.DIM; i++ {
			arr[i] = best.X[i] + p.F*(r1.X[i]-r2.X[i])
		}

		ret := models.Elem{
			X: arr,
		}
		return ret, nil
	},
	Name: "best1",
}
