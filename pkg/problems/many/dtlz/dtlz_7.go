package dtlz

import (
	"errors"
	"math"

	"github.com/nicholaspcr/gde3/pkg/problems/models"
)

// DTLZ7 multiObjective testcase
var DTLZ7 = models.ProblemFn{
	Fn: func(e *models.Elem, M int) error {
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
		e.Objs = make([]float64, len(objs))
		copy(e.Objs, objs)

		return nil
	},
	Name: "dtlz7",
}
