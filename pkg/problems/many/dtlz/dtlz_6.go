package dtlz

import (
	"errors"
	"math"

	"github.com/nicholaspcr/gde3/pkg/problems/models"
)

// DTLZ6 multiObjective testcase
var DTLZ6 = models.ProblemFn{
	Fn: func(e *models.Elem, M int) error {
		if len(e.X) <= M {
			return errors.New("need to have an M lesser than the amount of variables")
		}

		evalG := func(v []float64) float64 {
			g := 0.0
			for _, x := range v {
				g += math.Pow(x, 0.1) // consumes huge memory
			}
			return g
		}
		g := evalG(e.X[M-1:])
		t := math.Pi / (4.0 * (1.0 + g))

		objs := make([]float64, M)
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
				aux := M - (i + 1)
				prod *= math.Sin(theta[aux])
			}

			objs[i] = prod
		}

		// puts new objectives into the elem
		e.Objs = make([]float64, len(objs))
		copy(e.Objs, objs)
		return nil
	},
	Name: "dtlz6",
}
