package variants

import (
	"fmt"

	"github.com/nicholaspcr/gde3/pkg/models"
)

// Rand1 variant -> a + F(b - c)
// a,b,c are random elements
var rand1 = models.Variant{
	Fn: func(elems, rankZero models.Population, p models.VariantParams) (models.Vector, error) {
		if len(elems) < 3 {
			return models.Vector{}, fmt.Errorf(
				"no sufficient amount of elements in the population, should be bigger than three",
			)
		}

		// generating random indices different from current pos
		inds := make([]int, 4)
		inds[0] = p.CurrPos
		err := generateIndices(1, len(elems), inds)
		if err != nil {
			return models.Vector{}, err
		}

		result := models.Vector{}
		result.X = make([]float64, p.DIM)

		r1, r2, r3 := elems[inds[1]], elems[inds[2]], elems[inds[3]]
		for i := 0; i < p.DIM; i++ {
			result.X[i] = r1.X[i] + p.F*(r2.X[i]-r3.X[i])
		}
		return result, nil
	},
	VariantName: "rand1",
}
