package variants

import (
	"fmt"

	"gitlab.com/nicholaspcr/go-de/pkg/problems/models"
)

// Rand1 variant -> a + F(b - c)
// a,b,c are random elements
var rand1 VariantFn = VariantFn{
	Fn: func(elems, rankZero models.Elements, p Params) (models.Elem, error) {
		if len(elems) < 3 {
			return models.Elem{}, fmt.Errorf("no sufficient amount of elements in the population, should be bigger than three")
		}

		// generating random indices different from current pos
		inds := make([]int, 4)
		inds[0] = p.CurrPos
		err := generateIndices(1, len(elems), inds)
		if err != nil {
			return models.Elem{}, err
		}

		result := models.Elem{}
		result.X = make([]float64, p.DIM)

		r1, r2, r3 := elems[inds[1]], elems[inds[2]], elems[inds[3]]
		for i := 0; i < p.DIM; i++ {
			result.X[i] = r1.X[i] + p.F*(r2.X[i]-r3.X[i])
		}
		return result, nil
	},
	Name: "rand1",
}
