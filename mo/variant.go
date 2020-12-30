package mo

import "fmt"

// VariantFn function type of the multiple variants
type VariantFn func(elems Elements, p Params) (Elem, error)

// Rand1 variant -> a + F(b - c)
func rand1(elems Elements, p Params) (Elem, error) {
	if len(elems) < 3 {
		return Elem{}, fmt.Errorf("no sufficient amount of elements in the population, should be bigger than three")
	}
	inds := make([]int, 3)
	err := generateIndices(0, len(elems), inds)
	if err != nil {
		return Elem{}, err
	}

	e := Elem{}
	e.X = make([]float64, p.DIM)
	for i := 0; i < p.DIM; i++ {
		e.X[i] = elems[inds[0]].X[i] + p.F*(elems[inds[1]].X[i]-elems[inds[2]].X[i])
	}

	return e, nil
}
