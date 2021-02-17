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
	fn   func(elems, rankZero Elements, p varParams) (Elem, error)
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
	fn: func(elems, rankZero Elements, p varParams) (Elem, error) {
		if len(elems) < 3 {
			return Elem{}, fmt.Errorf("no sufficient amount of elements in the population, should be bigger than three")
		}

		// generating random indices different from current pos
		inds := make([]int, 4)
		inds[0] = p.currPos
		err := generateIndices(1, len(elems), inds)
		if err != nil {
			return Elem{}, err
		}

		result := Elem{}
		result.X = make([]float64, p.DIM)

		r1, r2, r3 := elems[inds[1]], elems[inds[2]], elems[inds[3]]
		for i := 0; i < p.DIM; i++ {
			result.X[i] = r1.X[i] + p.F*(r2.X[i]-r3.X[i])
		}
		return result, nil
	},
	Name: "rand1",
}

// rand2 a + F(b-c) + F(d-e)
// a,b,c,d,e are random elements
var rand2 VariantFn = VariantFn{
	fn: func(elems, rankZero Elements, p varParams) (Elem, error) {
		// generating random indices different from current pos
		ind := make([]int, 6)
		ind[0] = p.currPos
		err := generateIndices(1, len(elems), ind)
		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 4")
		}

		arr := make([]float64, p.DIM)
		r1, r2, r3, r4, r5 := elems[ind[1]], elems[ind[2]], elems[ind[3]], elems[ind[4]], elems[ind[5]]
		for i := 0; i < p.DIM; i++ {
			arr[i] = r1.X[i] + p.F*(r2.X[i]-r3.X[i]) + p.F*(r4.X[i]-r5.X[i])
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
	fn: func(elems, rankZero Elements, p varParams) (Elem, error) {
		index := make([]int, 3)
		index[0] = p.currPos
		err := generateIndices(1, len(elems), index)

		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 4")
		}

		arr := make([]float64, p.DIM)

		best := rankZero[rand.Intn(len(rankZero))]
		r1, r2 := elems[index[1]], elems[index[2]]
		for i := 0; i < p.DIM; i++ {
			arr[i] = best.X[i] + p.F*(r1.X[i]-r2.X[i])
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
	fn: func(elems, rankZero Elements, p varParams) (Elem, error) {
		// indices of the
		ind := make([]int, 5)
		ind[0] = p.currPos
		err := generateIndices(1, len(elems), ind)

		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 4")
		}

		arr := make([]float64, p.DIM)

		// random element from rankZero
		best := rankZero[rand.Int()%len(rankZero)]
		r1, r2, r3, r4 := elems[ind[1]], elems[ind[2]], elems[ind[3]], elems[ind[4]]
		for i := 0; i < p.DIM; i++ {
			arr[i] = best.X[i] + p.F*(r1.X[i]-r2.X[i]) + p.F*(r3.X[i]-r4.X[i])
		}

		ret := Elem{
			X: arr,
		}
		return ret, nil
	},
	Name: "best2",
}

// currToBest1 -> variant defined in JADE paper
var currToBest1 VariantFn = VariantFn{
	fn: func(elems, rankZero Elements, p varParams) (Elem, error) {
		ind := make([]int, 4)
		ind[0] = p.currPos
		err := generateIndices(1, len(elems), ind)

		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 5")
		}

		arr := make([]float64, p.DIM)

		r1, r2, r3 := elems[ind[1]], elems[ind[2]], elems[ind[3]]
		curr := elems[p.currPos]
		best := rankZero[rand.Int()%len(rankZero)]

		for i := 0; i < p.DIM; i++ {
			arr[i] = curr.X[i] + p.F*(best.X[i]-r1.X[i]) + p.F*(r2.X[i]-r3.X[i])
		}

		ret := Elem{
			X: arr,
		}
		return ret, nil
	},
	Name: "currtobest1",
}

// PBest implementation
var pbest VariantFn = VariantFn{
	fn: func(elems, rankZero Elements, p varParams) (Elem, error) {
		ind := make([]int, 3)
		ind[0] = p.currPos

		err := generateIndices(1, len(elems), ind)
		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 5")
		}

		pThRankZeroSz := int(math.Floor(float64(len(rankZero)) * p.P))

		var bestIndex int
		if pThRankZeroSz == 0 {
			bestIndex = 0
		} else {
			bestIndex = rand.Int() % pThRankZeroSz
		}

		arr := make([]float64, p.DIM)

		r1, r2 := elems[ind[1]], elems[ind[2]]
		curr := elems[p.currPos]
		best := rankZero[bestIndex]

		for i := 0; i < p.DIM; i++ {
			arr[i] = curr.X[i] + p.F*(best.X[i]-curr.X[i]) + p.F*(r1.X[i]-r2.X[i])
		}

		ret := Elem{
			X: arr,
		}
		return ret, nil
	},
	Name: "pbest",
}
