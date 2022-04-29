package dtlz

import (
	"errors"
	"math"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

type dtlz5 struct{}

func Dtlz5() models.Problem {
	return &dtlz5{}
}

func (v *dtlz5) Name() string {
	return "dtlz5"
}

func (v *dtlz5) Evaluate(e *models.Vector, M int) error {

	if len(e.X) <= M {
		return errors.New(
			"need to have an M lesser than the amount of variables",
		)
	}

	varSz := len(e.X)
	k := varSz - M + 1
	evalG := func(x []float64) float64 {
		g := 0.0
		for _, v := range x {
			g += (v - 0.5) * (v - 0.5)
		}
		return g
	}
	g := evalG(e.X[varSz-k:])
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
	e.Objs = make([]float64, len(newObjs))
	copy(e.Objs, newObjs)

	return nil
}
