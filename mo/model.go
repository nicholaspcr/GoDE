package mo

import (
	"errors"
	"math"
)

// Elem -> Element of population
type Elem struct {
	X      []float64
	objs   []float64
	crwdst float64
}

func (e *Elem) makeCpy() Elem {
	var ret Elem
	ret.X = make([]float64, len(e.X))
	ret.objs = make([]float64, len(e.objs))
	copy(ret.X, e.X)
	copy(ret.objs, e.objs)
	ret.crwdst = e.crwdst
	return ret
}

func (e *Elem) dominates(other Elem) bool {
	if len(e.objs) != len(other.objs) {
		return false
	}
	dominates := false
	for i := range e.objs {
		if e.objs[i] > other.objs[i] {
			return false
		}
		if e.objs[i] < other.objs[i] {
			dominates = true
		}
	}
	return dominates
}

type byCrwdst []Elem

func (x byCrwdst) Len() int           { return len(x) }
func (x byCrwdst) Less(i, j int) bool { return x[i].crwdst > x[j].crwdst }
func (x byCrwdst) Swap(i, j int)      { t := x[i].makeCpy(); x[i] = x[j].makeCpy(); x[j] = t }

type byFirstObj []Elem

func (x byFirstObj) Len() int           { return len(x) }
func (x byFirstObj) Less(i, j int) bool { return x[i].objs[0] < x[j].objs[0] }
func (x byFirstObj) Swap(i, j int)      { t := x[i].makeCpy(); x[i] = x[j].makeCpy(); x[j] = t }

// ZDT1 -> bi-objetive evaluation
func ZDT1(e *Elem) error {
	if len(e.X) < 2 {
		return errors.New("Need at least two variables/dimensions")
	}
	evalG := func(x []float64) float64 {
		g := 0.0
		for i := 1; i < len(x); i++ {
			g += x[i]
		}
		constant := 9.0 / (float64(len(x)) - 1.0)

		return 1.0 + constant*g
	}
	evalH := func(f, g float64) float64 {
		return 1.0 - math.Sqrt(f/g)
	}
	g := evalG(e.X)
	h := evalH(e.X[0], g)

	if math.IsNaN(h) == true {
		return errors.New("Sqrt of a negative number")
	}

	var newObjs []float64
	newObjs = append(newObjs, e.X[0])
	newObjs = append(newObjs, g*h)

	// puts new objectives into the elem
	e.objs = make([]float64, len(newObjs))
	copy(e.objs, newObjs)

	return nil
}

// ZDT2 -> bi-objetive evaluation
func ZDT2(e *Elem) error {
	if len(e.X) < 2 {
		return errors.New("Need at least two variables/dimensions")
	}
	evalG := func(x []float64) float64 {
		g := 0.0
		for i := 1; i < len(x); i++ {
			g += x[i]
		}
		constant := (9.0 / (float64(len(x)) - 1.0))

		return 1.0 + constant*g
	}
	evalH := func(f, g float64) float64 {
		return 1.0 - math.Pow(f/g, 2)
	}
	g := evalG(e.X)
	h := evalH(e.X[0], g)

	var newObjs []float64
	newObjs = append(newObjs, e.X[0])
	newObjs = append(newObjs, g*h)

	// puts new objectives into the elem
	e.objs = make([]float64, len(newObjs))
	copy(e.objs, newObjs)

	return nil
}

// ZDT3 -> bi-objetive evaluation
func ZDT3(e *Elem) error {
	if len(e.X) < 2 {
		return errors.New("Need at least two variables/dimensions")
	}
	evalG := func(x []float64) float64 {
		g := 0.0
		for i := 1; i < len(x); i++ {
			g += x[i]
		}
		constant := (9.0 / (float64(len(x)) - 1.0))

		return 1.0 + constant*g
	}
	evalH := func(f, g float64) float64 {
		return 1.0 - math.Sqrt(f/g) - (f/g)*math.Sin(10.0*f*math.Pi)
	}
	g := evalG(e.X)
	h := evalH(e.X[0], g)
	if math.IsNaN(h) {
		return errors.New("Sqrt of a negative number")
	}
	var newObjs []float64
	newObjs = append(newObjs, e.X[0])
	newObjs = append(newObjs, g*h)

	// puts new objectives into the elem
	e.objs = make([]float64, len(newObjs))
	copy(e.objs, newObjs)

	return nil
}

// ZDT4 -> bi-objetive evaluation
func ZDT4(e *Elem) error {
	if len(e.X) < 2 {
		return errors.New("Need at least two variables/dimensions")
	}
	evalG := func(x []float64) float64 {
		g := 0.0
		for i := 1; i < len(x); i++ {
			g += math.Pow(x[i], 2) - 10*math.Cos(4*math.Pi*x[i])
		}
		sz := float64(len(x) - 1)
		return 1.0 + 10.0*sz + g
	}
	evalH := func(f, g float64) float64 {
		return 1.0 - math.Sqrt(f/g)
	}

	g := evalG(e.X)
	h := evalH(e.X[0], g)

	var newObjs []float64
	newObjs = append(newObjs, e.X[0])
	newObjs = append(newObjs, g*h)

	// puts new objectives into the elem
	e.objs = make([]float64, len(newObjs))
	copy(e.objs, newObjs)

	return nil
}

// ZDT6 -> bi-objetive evaluation
func ZDT6(e *Elem) error {
	if len(e.X) < 2 {
		return errors.New("Need at least two variables/dimensions")
	}
	evalF := func(x float64) float64 {
		f := math.Exp(-4.0 * x)
		f = f * math.Pow(math.Sin(6*math.Pi*x), 6)
		f = 1 - f
		return f
	}
	evalG := func(x []float64) float64 {
		g := 0.0
		for i := 1; i < len(x); i++ {
			g += x[i]
		}
		g = g / float64(len(x)-1)
		g = math.Pow(g, 1.0/4)
		g = g*9 + 1.0
		return g
	}
	evalH := func(f, g float64) float64 {
		return 1.0 - math.Pow(f/g, 2)
	}
	F := evalF(e.X[0])
	G := evalG(e.X)
	H := evalH(F, G)

	var newObjs []float64
	newObjs = append(newObjs, F)
	newObjs = append(newObjs, G*H)

	// puts new objectives into the elem
	e.objs = make([]float64, len(newObjs))
	copy(e.objs, newObjs)

	return nil
}

// VNT1 -> https://ti.arc.nasa.gov/m/pub-archive/archive/1163.pdf
// VNT1 -> recebe 2 variaveis e otimiza 3 funções objetivo
func VNT1(e *Elem) error {
	if len(e.X) != 2 {
		return errors.New("Need at have only two variables/dimensions")
	}

	a, b := e.X[0], e.X[1]

	powSum := math.Pow(a, 2) + math.Pow(b, 2)
	f1 := 0.5*(powSum) + math.Sin(powSum)
	f2 := 15.0 + math.Pow(3*a-2*b+4, 2)/8.0 + math.Pow(a-b+1, 2)/27.0
	f3 := -1.1*math.Exp((-1)*powSum) + 1.0/(powSum+1)

	var newObjs []float64
	newObjs = append(newObjs, f1)
	newObjs = append(newObjs, f2)
	newObjs = append(newObjs, f3)

	// puts new objectives into the elem
	e.objs = make([]float64, len(newObjs))
	copy(e.objs, newObjs)

	return nil
}

/*
All the DTLZ implementations are translations of the python implementation made by pymoo
https://pymoo.org/problems/many/dtlz.html
*/

// DTLZ1 multiObjective testcase
func DTLZ1(e *Elem, M int) error {
	if len(e.X) <= M {
		return errors.New("need to have an M lesser than the amount of variables")
	}

	newObjs := make([]float64, M)

	evalG := func(x []float64) float64 {
		k := len(x) - M + 1
		g := 0.0
		for i := len(x) - k; i < len(x); i++ {
			g += (x[i]-0.5)*(x[i]-0.5) - math.Cos(20*math.Pi*(x[i]-0.5))
		}
		return 100 * (float64(k) + g)
	}
	g := evalG(e.X)

	for i := 0; i < M; i++ {
		prod := 0.5 * (1 + g)
		for j := 0; j < M-(i+1); j++ {
			prod *= e.X[j]
		}
		if i != 0 {
			prod *= (1 - e.X[M-(i+1)])
		}
		newObjs[i] = prod
	}

	// puts new objectives into the elem
	e.objs = make([]float64, len(newObjs))
	copy(e.objs, newObjs)

	return nil
}

// DTLZ2  multiObjective testcase
func DTLZ2(e *Elem, M int) error {
	if len(e.X) <= M {
		return errors.New("need to have an M lesser than the amount of variables")
	}

	newObjs := make([]float64, M)

	evalG := func(x []float64) float64 {
		k := len(x) - M + 1
		g := 0.0
		for i := len(x) - k; i < len(x); i++ {
			g += (x[i] - 0.5) * (x[i] - 0.5)
		}
		return g
	}
	g := evalG(e.X)

	for i := 0; i < M; i++ {
		prod := (1 + g)
		for j := 0; j < M-(i+1); j++ {
			prod *= math.Cos(e.X[j] * 0.5 * math.Pi)
		}
		if i != 0 {
			prod *= math.Sin(0.5 * math.Pi * e.X[M-(i+1)])
		}
		newObjs[i] = prod
	}

	// puts new objectives into the elem
	e.objs = make([]float64, len(newObjs))
	copy(e.objs, newObjs)

	return nil
}

// DTLZ3 multiObjective testcase
func DTLZ3(e *Elem, M int) error {
	if len(e.X) <= M {
		return errors.New("need to have an M lesser than the amount of variables")
	}

	newObjs := make([]float64, M)

	evalG := func(x []float64) float64 {
		k := len(x) - M + 1
		g := 0.0
		for i := len(x) - k; i < len(x); i++ {
			g += (x[i]-0.5)*(x[i]-0.5) - math.Cos(20.0*(x[i]-0.5)*math.Pi)
		}
		return 100 * (float64(k) + g)
	}
	g := evalG(e.X)

	for i := 0; i < M; i++ {
		prod := (1.0 + g)
		for j := 0; j < M-(i+1); j++ {
			prod *= math.Cos(e.X[j] * 0.5 * math.Pi)
		}
		if i != 0 {
			prod *= math.Sin(e.X[M-(i+1)] * 0.5 * math.Pi)
		}
		newObjs[i] = prod
	}

	// puts new objectives into the elem
	e.objs = make([]float64, len(newObjs))
	copy(e.objs, newObjs)

	return nil
}

// DTLZ4 multiObjective testcase
func DTLZ4(e *Elem, M int) error {
	if len(e.X) <= M {
		return errors.New("need to have an M lesser than the amount of variables")
	}

	newObjs := make([]float64, M)

	evalG := func(x []float64) float64 {
		k := len(x) - M + 1
		g := 0.0
		for i := len(x) - k; i < len(x); i++ {
			g += (x[i] - 0.5) * (x[i] - 0.5)
		}
		return g
	}
	g := evalG(e.X)

	for i := 0; i < M; i++ {
		prod := (1 + g)
		for j := 0; j < M-(i+1); j++ {
			prod *= math.Cos(math.Pow(e.X[j], 100) * math.Pi / 2.0)
		}
		if i != 0 {
			prod *= math.Sin(math.Pow(e.X[M-(i+1)], 100) * math.Pi / 2.0)
		}
		newObjs[i] = prod
	}

	// puts new objectives into the elem
	e.objs = make([]float64, len(newObjs))
	copy(e.objs, newObjs)

	return nil
}

// DTLZ5 multiObjective testcase
func DTLZ5(e *Elem, M int) error {
	if len(e.X) <= M {
		return errors.New("need to have an M lesser than the amount of variables")
	}

	newObjs := make([]float64, M)

	evalG := func(x []float64) float64 {
		k := len(x) - M + 1
		g := 0.0
		for i := len(x) - k; i < len(x); i++ {
			g += (x[i] - 0.5) * (x[i] - 0.5)
		}
		return g
	}
	g := evalG(e.X)
	t := math.Pi / (4.0 * (1 + g))

	theta := make([]float64, M-1)
	theta[0] = e.X[0] * math.Pi / 2.0
	for i := 1; i < M-1; i++ {
		theta[i] = t * (1.0 + 2.0*g*e.X[i])
	}

	for i := 0; i < M; i++ {
		prod := (1 + g)
		for j := 0; j < M-(i+1); j++ {
			prod *= math.Cos(theta[j])
		}
		if i != 0 {
			prod *= math.Sin(theta[M-(i+1)])
		}
		newObjs[i] = prod
	}

	// puts new objectives into the elem
	e.objs = make([]float64, len(newObjs))
	copy(e.objs, newObjs)

	return nil
}

// DTLZ6 multiObjective testcase
func DTLZ6(e *Elem, M int) error {
	if len(e.X) <= M {
		return errors.New("need to have an M lesser than the amount of variables")
	}

	newObjs := make([]float64, M)

	evalG := func(x []float64) float64 {
		k := len(x) - M + 1
		g := 0.0
		for i := len(x) - k; i < len(x); i++ {
			g += math.Pow(x[i], 0.1)
		}
		return g
	}
	g := evalG(e.X)
	t := math.Pi / (4.0 * (1 + g))

	theta := make([]float64, M-1)
	theta[0] = e.X[0] * math.Pi / 2.0
	for i := 1; i < M-1; i++ {
		theta[i] = t * (1.0 + 2.0*g*e.X[i])
	}

	for i := 0; i < M; i++ {
		prod := (1 + g)
		for j := 0; j < M-(i+1); j++ {
			prod *= math.Cos(theta[j])
		}
		if i != 0 {
			prod *= math.Sin(theta[M-(i+1)])
		}
		newObjs[i] = prod
	}

	// puts new objectives into the elem
	e.objs = make([]float64, len(newObjs))
	copy(e.objs, newObjs)

	return nil
}

// DTLZ7 multiObjective testcase
func DTLZ7(e *Elem, M int) error {
	if len(e.X) <= M {
		return errors.New("need to have an M lesser than the amount of variables")
	}

	newObjs := make([]float64, M)

	evalG := func(x []float64) float64 {
		k := len(x) - M + 1
		g := 0.0
		for i := len(x) - k; i < len(x); i++ {
			g += math.Pow(x[i], 0.1)
		}
		return g
	}
	evalH := func(x []float64, g float64) float64 {
		h := 0.0
		for i := 0; i < M-1; i++ {
			h += (x[i] / (1 + g) * (1 + math.Sin(3*math.Pi*x[i])))
		}
		return float64(M) - h
	}
	g := evalG(e.X)
	h := evalH(e.X, g)

	for i := 0; i < M-1; i++ {
		newObjs[i] = e.X[i]
	}
	newObjs[M-1] = (1 + g) * h

	// puts new objectives into the elem
	e.objs = make([]float64, len(newObjs))
	copy(e.objs, newObjs)

	return nil
}
