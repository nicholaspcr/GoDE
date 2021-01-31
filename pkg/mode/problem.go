package mo

import (
	"errors"
	"math"
	"strings"
)

// ProblemFn definition of the test case functions
type ProblemFn struct {
	fn   func(e *Elem, M int) error
	Name string
}

// GetProblemByName -> returns the problem function
func GetProblemByName(Name string) ProblemFn {
	Name = strings.ToLower(Name)
	problems := map[string]ProblemFn{
		zdt1.Name:  zdt1,
		zdt2.Name:  zdt2,
		zdt3.Name:  zdt3,
		zdt4.Name:  zdt4,
		zdt6.Name:  zdt6,
		vnt1.Name:  vnt1,
		dtlz1.Name: dtlz1,
		dtlz2.Name: dtlz2,
		dtlz3.Name: dtlz3,
		dtlz4.Name: dtlz4,
		dtlz5.Name: dtlz5,
		dtlz6.Name: dtlz6,
		dtlz7.Name: dtlz7,
	}
	var problem ProblemFn
	for k, v := range problems {
		if Name == k {
			problem = v
			break
		}
	}
	return problem
}

// ZDT1 -> bi-objetive evaluation
var zdt1 = ProblemFn{
	fn: func(e *Elem, M int) error {
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
	},
	Name: "zdt1",
}

// ZDT2 -> bi-objetive evaluation
var zdt2 = ProblemFn{
	fn: func(e *Elem, M int) error {
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
	},
	Name: "zdt2",
}

// ZDT3 -> bi-objetive evaluation
var zdt3 = ProblemFn{
	fn: func(e *Elem, M int) error {
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
	},
	Name: "zdt3",
}

// ZDT4 -> bi-objetive evaluation
var zdt4 = ProblemFn{
	fn: func(e *Elem, M int) error {
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
	},
	Name: "zdt4",
}

// ZDT6 -> bi-objetive evaluation
var zdt6 = ProblemFn{
	fn: func(e *Elem, M int) error {
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
	},
	Name: "zdt6",
}

// VNT1 -> https://ti.arc.nasa.gov/m/pub-archive/archive/1163.pdf
// VNT1 -> recebe 2 variaveis e otimiza 3 funções objetivo
var vnt1 = ProblemFn{
	fn: func(e *Elem, M int) error {
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
	},
	Name: "vnt1",
}

/*
All the DTLZ implementations are translations of the python implementation made by pymoo
https://pymoo.org/problems/many/dtlz.html
*/

// DTLZ1 multiObjective testcase
var dtlz1 = ProblemFn{
	fn: func(e *Elem, M int) error {
		if len(e.X) <= M {
			return errors.New("need to have an M lesser than the amount of variables")
		}

		k := len(e.X) - M + 1
		g := 0.0
		varSz := len(e.X)
		for _, v := range e.X[varSz-k:] {
			g += (v-0.5)*(v-0.5) - math.Cos(20.0*math.Pi*(v-0.5))
		}
		g = 100 * (float64(k) + g)
		objs := make([]float64, M)
		for i := range objs {
			objs[i] = 0.5 * (1.0 + g)
			for j := 0; j < M-(i+1); j++ {
				objs[i] *= e.X[j]
			}
			if i != 0 {
				objs[i] *= 1 - e.X[M-(i+1)]
			}
		}
		e.objs = make([]float64, M)
		copy(e.objs, objs)
		return nil
	},
	Name: "dtlz1",
}

// DTLZ2  multiObjective testcase
var dtlz2 = ProblemFn{
	fn: func(e *Elem, M int) error {
		if len(e.X) <= M {
			return errors.New("need to have an M lesser than the amount of variables")
		}

		evalG := func(x []float64) float64 {
			g := 0.0
			for _, v := range x {
				g += (v - 0.5) * (v - 0.5)
			}
			return g
		}
		g := evalG(e.X[M:])

		newObjs := make([]float64, M)
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
	},
	Name: "dtlz2",
}

// DTLZ3 multiObjective testcase
var dtlz3 = ProblemFn{
	fn: func(e *Elem, M int) error {
		if len(e.X) <= M {
			return errors.New("need to have an M lesser than the amount of variables")
		}

		k := len(e.X) - M + 1
		g := 0.0
		varSz := len(e.X)
		for _, v := range e.X[varSz-k:] {
			g += (v-0.5)*(v-0.5) - math.Cos(20*math.Pi*(v-0.5))
		}
		g = 100 * (float64(k) + g)
		objs := make([]float64, M)
		for i := range objs {
			objs[i] = (1.0 + g)
			for j := 0; j < M-(i+1); j++ {
				objs[i] *= math.Cos(e.X[j] * 0.5 * math.Pi)
			}
			if i != 0 {
				objs[i] *= math.Sin(e.X[M-(i+1)] * 0.5 * math.Pi)
			}
		}
		// puts new objectives into the elem
		e.objs = make([]float64, len(objs))
		copy(e.objs, objs)
		return nil
	},
	Name: "dtlz3",
}

// DTLZ4 multiObjective testcase
var dtlz4 = ProblemFn{
	fn: func(e *Elem, M int) error {
		if len(e.X) <= M {
			return errors.New("need to have an M lesser than the amount of variables")
		}
		evalG := func(x []float64) float64 {
			g := 0.0
			for _, v := range x {
				g += (v - 0.5) * (v - 0.5)
			}
			return g
		}
		g := evalG(e.X[M:])

		newObjs := make([]float64, M)
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
	},
	Name: "dtlz4",
}

// DTLZ5 multiObjective testcase
var dtlz5 = ProblemFn{
	fn: func(e *Elem, M int) error {
		if len(e.X) <= M {
			return errors.New("need to have an M lesser than the amount of variables")
		}
		evalG := func(x []float64) float64 {
			g := 0.0
			for _, v := range x {
				g += (v - 0.5) * (v - 0.5)
			}
			return g
		}
		g := evalG(e.X[M:])
		t := math.Pi / (4.0 * (1 + g))

		newObjs := make([]float64, M)
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
	},
	Name: "dtlz5",
}

// DTLZ6 multiObjective testcase
var dtlz6 = ProblemFn{
	fn: func(e *Elem, M int) error {
		if len(e.X) <= M {
			return errors.New("need to have an M lesser than the amount of variables")
		}
		evalG := func(x []float64) float64 {
			g := 0.0
			for _, v := range x {
				g += math.Pow(v, 0.1)
			}
			return g
		}
		g := evalG(e.X[M:])
		t := math.Pi / (4.0 * (1 + g))

		theta := make([]float64, M-1)
		theta[0] = e.X[0] * math.Pi / 2.0
		for i := 1; i < M-1; i++ {
			theta[i] = t * (1.0 + 2.0*g*e.X[i])
		}

		newObjs := make([]float64, M)
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
	},
	Name: "dtlz7",
}

// DTLZ7 multiObjective testcase
var dtlz7 = ProblemFn{
	fn: func(e *Elem, M int) error {
		if len(e.X) <= M {
			return errors.New("need to have an M lesser than the amount of variables")
		}
		varSz := len(e.X)
		k := varSz - M + 1

		// calculating the value of the constant G
		g := 0.0
		for _, v := range e.X[varSz-k:] {
			g += v
		}
		g = 1.0 + (9.0*g)/float64(k)

		// calculating the value of the constant H
		h := 0.0
		for _, v := range e.X[:M-1] {
			h += (v / (1.0 + g)) * (1 + math.Sin(3.0*math.Pi*v))
		}
		h = float64(M) - h

		// calculating objs values
		objs := make([]float64, M)
		for i := range objs {
			objs[i] = e.X[i]
		}
		objs[M-1] = (1.0 + g) * h
		// puts new objectives into the elem
		e.objs = make([]float64, len(objs))
		copy(e.objs, objs)

		return nil
	},
	Name: "dtlz7",
}
