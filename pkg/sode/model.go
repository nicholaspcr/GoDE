package so

import "math"

// Params of the sode
type Params struct {
	NP, DIM, GEN, EXECS   int
	FLOOR, CEIL, CR, F, P float64
}

// Input -> input used in main
type Input struct {
	Eq  Equation
	Var Variant
}

// Elem - Element with a point and it's fitness based on the function f(Point)
type Elem struct {
	X   []float64
	fit float64
}

// Copy the Elem struct
func (e *Elem) Copy() Elem {
	var ret Elem
	ret.X = make([]float64, len(e.X))
	copy(ret.X, e.X)
	ret.fit = e.fit
	return ret
}

// Elements is a slice of Elem
type Elements []Elem

// Copy of the slice of Elem
func (e Elements) Copy() Elements {
	arr := make(Elements, len(e))
	for i, v := range e {
		arr[i] = v.Copy()
	}
	return arr
}
func elemArrCopy(arr []Elem) []Elem {
	ret := make([]Elem, len(arr))
	for i := 0; i < len(arr); i++ {
		ret[i] = arr[i].Copy()
	}
	return ret
}

type byFit Elements

func (x byFit) Len() int           { return len(x) }
func (x byFit) Less(i, j int) bool { return x[i].fit < x[j].fit }
func (x byFit) Swap(i, j int)      { t := x[i].Copy(); x[i] = x[j].Copy(); x[i] = t }

// Equation - used as a parameter to determine fitness in a DE
type Equation struct {
	calcFunc func(x []float64) float64
	fileName string
}

// Ackley implementation
var Ackley Equation = Equation{
	calcFunc: func(x []float64) float64 {
		var sqrdSum, cosSum, invDim, ret float64 // squaredSum, sumOfCos, inverseDim
		dim := len(x)
		invDim = 1 / float64(dim)
		for i := 0; i < dim; i++ {
			sqrdSum += math.Pow(x[i], 2)
			cosSum += math.Cos(2 * math.Pi * x[i])
		}

		ret = (-20)*math.Exp((-0.2)*math.Sqrt(invDim*sqrdSum)) - math.Exp(invDim*cosSum) + 20 + math.E
		return ret
	},
	fileName: "ackley",
}

// Rastrigin implementation
var Rastrigin Equation = Equation{
	calcFunc: func(x []float64) float64 {
		var ret float64
		for i := 0; i < len(x); i++ {
			ret += math.Pow(x[i], 2) - 10*math.Cos(2*math.Pi*x[i]) + 10
		}
		return ret
	},
	fileName: "rastrigin",
}

// Schwefel implementation
var Schwefel Equation = Equation{
	calcFunc: func(x []float64) float64 {
		var sum float64 = 0.0
		dim := len(x)
		for i := 0; i < dim; i++ {
			sum += x[i] * math.Sin(math.Sqrt(math.Abs(x[i])))
		}
		ret := 418.9829*float64(dim) - sum
		return ret
	},
	fileName: "schwefel",
}

// Rosenbrock implementation
var Rosenbrock Equation = Equation{
	calcFunc: func(x []float64) float64 {
		var ret float64
		dim := len(x)
		for i := 0; i < dim-1; i++ {
			ret += 100*math.Pow(math.Pow(x[i], 2)-x[i+1], 2) + math.Pow(x[i]-1, 2)
		}
		return ret
	},
	fileName: "rosenbrock",
}
