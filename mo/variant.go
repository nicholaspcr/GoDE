package mo

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
)

type varParams struct {
	DIM     int
	F       float64
	currPos int
	P       float64
}

// VariantFn function type of the multiple variants
type VariantFn struct {
	fn   func(elems Elements, p varParams) (Elem, error)
	Name string
}

// GetVariantByName -> Returns the variant function
func GetVariantByName(name string) VariantFn {
	name = strings.ToLower(name)
	variants := map[string]VariantFn{
		"rand1":       rand1,
		"rand2":       rand2,
		"best1":       best1,
		"best2":       best2,
		"currtobest1": currToBest1,
		"currtobest2": currToBest2,
		"pbest":       pbest,
	}
	for k, v := range variants {
		if name == k {
			return v
		}
	}
	return VariantFn{}
}

// Rand1 variant -> a + F(b - c)
// a,b,c are random elements
var rand1 VariantFn = VariantFn{
	fn: func(elems Elements, p varParams) (Elem, error) {
		if len(elems) < 3 {
			return Elem{}, fmt.Errorf("no sufficient amount of elements in the population, should be bigger than three")
		}
		inds := make([]int, 3)
		err := generateIndices(0, len(elems), inds)
		if err != nil {
			return Elem{}, err
		}

		result := Elem{}
		result.X = make([]float64, p.DIM)
		for i := 0; i < p.DIM; i++ {
			result.X[i] = elems[inds[0]].X[i] + p.F*(elems[inds[1]].X[i]-elems[inds[2]].X[i])
		}
		return result, nil
	},
	Name: "rand1",
}

// rand2 a + F(b-c) + F(d-e)
// a,b,c,d,e are random elements
var rand2 VariantFn = VariantFn{
	fn: func(elems Elements, p varParams) (Elem, error) {
		ind := make([]int, 7)
		ind[0] = p.currPos
		err := generateIndices(1, len(elems), ind)
		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 4")
		}

		arr := make([]float64, p.DIM)
		a, b, c, d, e := elems[ind[1]], elems[ind[2]], elems[ind[3]], elems[ind[4]], elems[ind[5]]
		for i := 0; i < p.DIM; i++ {
			arr[i] = a.X[i] + p.F*(b.X[i]-c.X[i]) + p.F*(d.X[i]-e.X[i])
		}
		ret := Elem{
			X: arr,
		}
		return ret, nil
	},
	Name: "rand2",
}

// best1  current_best + F(a-b)
// a,b are random elements
var best1 VariantFn = VariantFn{
	fn: func(elems Elements, p varParams) (Elem, error) {
		index := make([]int, 4)
		index[0] = p.currPos
		index[1] = 0 // best in pop
		err := generateIndices(2, len(elems), index)
		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 4")
		}

		arr := make([]float64, p.DIM)
		a, b, c := elems[index[1]], elems[index[2]], elems[index[3]]
		for i := 0; i < p.DIM; i++ {
			arr[i] = a.X[i] + p.F*(b.X[i]-c.X[i])
		}
		ret := Elem{
			X: arr,
		}
		return ret, nil
	},
	Name: "best1",
}

// Best2 current_best + F(a-b) + F(c-d)
// a,b,c,d are random elements
var best2 VariantFn = VariantFn{
	fn: func(elems Elements, p varParams) (Elem, error) {
		ind := make([]int, 6)
		ind[0] = p.currPos
		ind[1] = 0 // best in elems
		err := generateIndices(2, len(elems), ind)
		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 4")
		}

		arr := make([]float64, p.DIM)
		a, b, c, d, e := elems[ind[1]], elems[ind[2]], elems[ind[3]], elems[ind[4]], elems[ind[5]]
		for i := 0; i < p.DIM; i++ {
			arr[i] = a.X[i] + p.F*(b.X[i]-c.X[i]) + p.F*(d.X[i]-e.X[i])
		}
		ret := Elem{
			X: arr,
		}
		return ret, nil
	},
	Name: "best2",
}

// TODO:
// link do artigo ->
var currToBest1 VariantFn = VariantFn{
	fn: func(elems Elements, p varParams) (Elem, error) {
		ind := make([]int, 5)
		ind[0] = p.currPos
		ind[1] = 0 // best in pop
		err := generateIndices(2, len(elems), ind)
		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 5")
		}
		arr := make([]float64, p.DIM)
		a, b, c, d, e := elems[ind[0]], elems[ind[1]], elems[ind[2]], elems[ind[3]], elems[ind[4]]
		for i := 0; i < p.DIM; i++ {
			arr[i] = a.X[i] + p.F*(b.X[i]-c.X[i]) + p.F*(d.X[i]-e.X[i])
		}
		ret := Elem{
			X: arr,
		}
		return ret, nil
	},
	Name: "current-to-best-1",
}

// TODO
// link do artigo ->
var currToBest2 VariantFn = VariantFn{
	fn: func(elems Elements, p varParams) (Elem, error) {
		ind := make([]int, 4)
		ind[0] = p.currPos
		ind[1] = 0 // best in population
		err := generateIndices(2, len(elems), ind)
		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 4")
		}
		arr := make([]float64, p.DIM)
		a, b, c, d := elems[ind[0]], elems[ind[1]], elems[ind[2]], elems[ind[3]]
		for i := 0; i < p.DIM; i++ {
			arr[i] = a.X[i] + p.F*(b.X[i]-a.X[i]) + p.F*(c.X[i]-d.X[i])
		}
		ret := Elem{
			X: arr,
		}
		return ret, nil
	},
	Name: "current-to-best-2",
}

// PBest implementation
var pbest VariantFn = VariantFn{
	fn: func(elems Elements, p varParams) (Elem, error) {
		popSz := float64(len(elems))
		ceilRand := int(math.Floor(popSz * p.P))
		var randPIndex int
		if ceilRand == 0 {
			randPIndex = ceilRand
		} else {
			randPIndex = rand.Int() % ceilRand
		}
		index := make([]int, 4)
		index[0] = p.currPos
		index[1] = randPIndex
		err := generateIndices(2, len(elems), index)
		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 4")
		}
		arr := make([]float64, p.DIM)
		curr, pB, a, b := index[0], index[1], index[2], index[3]
		for i := 0; i < p.DIM; i++ {
			arr[i] = elems[curr].X[i] + p.F*(elems[pB].X[i]-elems[curr].X[i]) + p.F*(elems[a].X[i]-elems[b].X[i])
		}
		ret := Elem{
			X: arr,
		}
		return ret, nil
	},
	Name: "pbest",
}
