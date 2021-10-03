// Package many package has the implementation of the dtlz problems
///
/*
The implementations are translations of the python code made by pymoo
https://pymoo.org/problems/many/dtlz.html
*/
package dtlz

import (
	"errors"
	"math"

	"github.com/nicholaspcr/gde3/pkg/models"
)

// DTLZ1 multiObjective testcase
var DTLZ1 = models.Problem{
	Fn: func(e *models.Vector, M int) error {
		if len(e.X) <= M {
			return errors.New(
				"need to have an M lesser than the amount of variables",
			)
		}

		evalG := func(v []float64) float64 {
			g := 0.0
			for _, x := range v {
				g += (x-0.5)*(x-0.5) - math.Cos(20.0*math.Pi*(x-0.5))
			}
			k := float64(len(v))
			return 100.0 * (k + g)
		}
		g := evalG(e.X[M-1:])

		newObjs := make([]float64, M)
		for i := 0; i < M; i++ {
			prod := (1.0 + g) * 0.5
			for j := 0; j < M-(i+1); j++ {
				prod *= e.X[j]
			}
			if i != 0 {
				prod *= (1 - e.X[M-(i+1)])
			}
			newObjs[i] = prod
		}

		// puts new objectives into the elem
		e.Objs = make([]float64, len(newObjs))
		copy(e.Objs, newObjs)

		return nil
	},
	ProblemName: "dtlz1",
}
