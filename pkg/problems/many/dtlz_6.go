package many

import (
	"errors"
	"math"

	"gitlab.com/nicholaspcr/go-de/pkg/problems/models"
)

// DTLZ6 multiObjective testcase
var DTLZ6 = models.ProblemFn{
	Fn: func(e *models.Elem, M int) error {
		if len(e.X) <= M {
			return errors.New("need to have an M lesser than the amount of variables")
		}

		var dtlzG float64
		evalG := func(v []float64) float64 {
			dtlzG = 0.0
			for _, x := range v {
				dtlzG += math.Pow(x, 1.0/10.0)
			}
			return dtlzG
		}
		g := evalG(e.X[M-1:])
		t := math.Pi / (4.0 * (1 + g))

		theta := make([]float64, M-1)
		theta[0] = e.X[0] * math.Pi / 2.0
		for i := 1; i < M-1; i++ {
			theta[i] = t * (1.0 + 2.0*g*e.X[i])
		}

		newObjs := make([]float64, M)
		for i := 0; i < M; i++ {
			newObjs[i] = (1 + g)
			for j := 0; j < M-(i+1); j++ {
				newObjs[i] *= math.Cos(theta[j])
			}
			if i != 0 {
				newObjs[i] *= math.Sin(theta[M-(i+1)])
			}
		}

		// puts new objectives into the elem
		e.Objs = make([]float64, len(newObjs))
		copy(e.Objs, newObjs)

		return nil
	},
	Name: "dtlz6",
}
