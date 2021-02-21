package variants

import (
	"errors"

	"gitlab.com/nicholaspcr/go-de/pkg/problems/models"
)

// rand2 a + F(b-c) + F(d-e)
// a,b,c,d,e are random elements
var rand2 VariantFn = VariantFn{
	Fn: func(elems, rankZero models.Elements, p Params) (models.Elem, error) {
		// generating random indices different from current pos
		ind := make([]int, 6)
		ind[0] = p.CurrPos
		err := generateIndices(1, len(elems), ind)
		if err != nil {
			return models.Elem{}, errors.New("insufficient size for the population, must me equal or greater than 4")
		}

		arr := make([]float64, p.DIM)
		r1, r2, r3, r4, r5 := elems[ind[1]], elems[ind[2]], elems[ind[3]], elems[ind[4]], elems[ind[5]]
		for i := 0; i < p.DIM; i++ {
			arr[i] = r1.X[i] + p.F*(r2.X[i]-r3.X[i]) + p.F*(r4.X[i]-r5.X[i])
		}
		ret := models.Elem{
			X: arr,
		}
		return ret, nil
	},
	Name: "rand2",
}
