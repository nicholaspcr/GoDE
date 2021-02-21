package variants

import (
	"errors"
	"math/rand"

	"gitlab.com/nicholaspcr/go-de/pkg/problems/models"
)

// currToBest1 -> variant defined in JADE paper
var currToBest1 VariantFn = VariantFn{
	Fn: func(elems, rankZero models.Elements, p Params) (models.Elem, error) {
		ind := make([]int, 4)
		ind[0] = p.CurrPos
		err := generateIndices(1, len(elems), ind)

		if err != nil {
			return models.Elem{}, errors.New("insufficient size for the population, must me equal or greater than 5")
		}

		arr := make([]float64, p.DIM)

		r1, r2, r3 := elems[ind[1]], elems[ind[2]], elems[ind[3]]
		curr := elems[p.CurrPos]
		best := rankZero[rand.Int()%len(rankZero)]

		for i := 0; i < p.DIM; i++ {
			arr[i] = curr.X[i] + p.F*(best.X[i]-r1.X[i]) + p.F*(r2.X[i]-r3.X[i])
		}

		ret := models.Elem{
			X: arr,
		}
		return ret, nil
	},
	Name: "currtobest1",
}
