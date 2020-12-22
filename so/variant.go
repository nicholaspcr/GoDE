package so

import (
	"errors"
	"math"
	"math/rand"
)

// Variant type for changing the variant applied on the ED
type Variant struct {
	funcName   string
	makeMutant func(pop []Elem, F, P float64, currPos, dim int) (Elem, error)
}

// Rand1 implementation
var Rand1 Variant = Variant{
	makeMutant: func(pop []Elem, F, P float64, currPos, dim int) (Elem, error) {
		index := make([]int, 4)
		index[0] = currPos
		err := generateIndices(1, len(pop), index)
		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 4")
		}
		arr := make([]float64, dim)
		a, b, c := pop[index[1]], pop[index[2]], pop[index[3]]
		for i := 0; i < dim; i++ {
			arr[i] = a.X[i] + F*(b.X[i]-c.X[i])
		}
		ret := Elem{
			X:   arr,
			fit: 0.0,
		}
		return ret, nil
	},
	funcName: "rand1",
}

// Rand2 implementation
var Rand2 Variant = Variant{
	makeMutant: func(pop []Elem, F, P float64, currPos, dim int) (Elem, error) {
		index := make([]int, 7)
		index[0] = currPos
		err := generateIndices(1, len(pop), index)
		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 4")
		}

		arr := make([]float64, dim)
		a, b, c, d, e := pop[index[1]], pop[index[2]], pop[index[3]], pop[index[4]], pop[index[5]]
		for i := 0; i < dim; i++ {
			arr[i] = a.X[i] + F*(b.X[i]-c.X[i]) + F*(d.X[i]-e.X[i])
		}
		ret := Elem{
			X:   arr,
			fit: 0.0,
		}
		return ret, nil
	},
	funcName: "rand2",
}

// Best1 implementation
var Best1 Variant = Variant{
	makeMutant: func(pop []Elem, F, P float64, currPos, dim int) (Elem, error) {
		index := make([]int, 4)
		index[0] = currPos
		index[1] = 0 // best in pop
		err := generateIndices(2, len(pop), index)
		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 4")
		}

		arr := make([]float64, dim)
		a, b, c := pop[index[1]], pop[index[2]], pop[index[3]]
		for i := 0; i < dim; i++ {
			arr[i] = a.X[i] + F*(b.X[i]-c.X[i])
		}
		ret := Elem{
			X:   arr,
			fit: 0.0,
		}
		return ret, nil
	},
	funcName: "best1",
}

// Best2 implementation
var best2 Variant = Variant{
	makeMutant: func(pop []Elem, F, P float64, currPos, dim int) (Elem, error) {
		index := make([]int, 6)
		index[0] = currPos
		index[1] = 0 // best in pop
		err := generateIndices(2, len(pop), index)
		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 4")
		}

		arr := make([]float64, dim)
		a, b, c, d, e := pop[index[1]], pop[index[2]], pop[index[3]], pop[index[4]], pop[index[5]]
		for i := 0; i < dim; i++ {
			arr[i] = a.X[i] + F*(b.X[i]-c.X[i]) + F*(d.X[i]-e.X[i])
		}
		ret := Elem{
			X:   arr,
			fit: 0.0,
		}
		return ret, nil
	},
	funcName: "best2",
}

// TODO:
// link do artigo ->
var currToBestv1 Variant = Variant{
	makeMutant: func(pop []Elem, F, P float64, currPos, dim int) (Elem, error) {
		index := make([]int, 5)
		index[0] = currPos
		index[1] = 0 // best in pop
		err := generateIndices(2, len(pop), index)
		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 5")
		}
		arr := make([]float64, dim)
		a, b, c, d, e := pop[index[0]], pop[index[1]], pop[index[2]], pop[index[3]], pop[index[4]]
		for i := 0; i < dim; i++ {
			arr[i] = a.X[i] + F*(b.X[i]-c.X[i]) + F*(d.X[i]-e.X[i])
		}
		ret := Elem{
			X:   arr,
			fit: 0.0,
		}
		return ret, nil
	},
	funcName: "current-to-best-1",
}

// TODO
// link do artigo ->
var currToBestv2 Variant = Variant{
	makeMutant: func(pop []Elem, F, P float64, currPos, dim int) (Elem, error) {
		index := make([]int, 4)
		index[0] = currPos
		index[1] = 0 // best in pop
		err := generateIndices(2, len(pop), index)
		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 4")
		}
		arr := make([]float64, dim)
		a, b, c, d := pop[index[0]], pop[index[1]], pop[index[2]], pop[index[3]]
		for i := 0; i < dim; i++ {
			arr[i] = a.X[i] + F*(b.X[i]-a.X[i]) + F*(c.X[i]-d.X[i])
		}
		ret := Elem{
			X:   arr,
			fit: 0.0,
		}
		return ret, nil
	},
	funcName: "current-to-best-2",
}

// PBest implementation
var PBest Variant = Variant{
	makeMutant: func(pop []Elem, F, P float64, currPos, dim int) (Elem, error) {
		popSz := float64(len(pop))
		ceilRand := int(math.Floor(popSz * P))
		var randPIndex int
		if ceilRand == 0 {
			randPIndex = ceilRand
		} else {
			randPIndex = rand.Int() % ceilRand
		}
		index := make([]int, 4)
		index[0] = currPos
		index[1] = randPIndex
		err := generateIndices(2, len(pop), index)
		if err != nil {
			return Elem{}, errors.New("insufficient size for the population, must me equal or greater than 4")
		}
		arr := make([]float64, dim)
		curr, pB, a, b := index[0], index[1], index[2], index[3]
		for i := 0; i < dim; i++ {
			arr[i] = pop[curr].X[i] + F*(pop[pB].X[i]-pop[curr].X[i]) + F*(pop[a].X[i]-pop[b].X[i])
		}
		ret := Elem{
			X:   arr,
			fit: 0.0,
		}
		return ret, nil
	},
	funcName: "pbest",
}

// generates random indices in the int slice, r -> it's a pointer
func generateIndices(startInd, popSz int, r []int) error {
	if len(r) > popSz {
		return errors.New("insufficient elements in population to generate random indices")
	}
	for i := startInd; i < len(r); i++ {
		for done := false; !done; {
			r[i] = rand.Int() % popSz
			done = true
			for j := 0; j < i; j++ {
				done = done && r[j] != r[i]
			}
		}
	}
	return nil
}
