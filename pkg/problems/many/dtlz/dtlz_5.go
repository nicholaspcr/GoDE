package dtlz

import (
	"errors"
	"math"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type dtlz5 struct{}

// Dtlz5 returns the DTLZ5 test problem, a degenerate Pareto front problem.
// Domain: [0,1]^n, Objectives: m (configurable)
func Dtlz5() problems.Interface {
	return &dtlz5{}
}

func (v *dtlz5) Name() string {
	return "dtlz5"
}

func (v *dtlz5) Evaluate(e *models.Vector, m int) error {
	if len(e.Elements) <= m {
		return errors.New(
			"need to have an m lesser than the amount of variables",
		)
	}

	varSz := len(e.Elements)
	k := varSz - m + 1
	evalG := func(x []float64) float64 {
		g := 0.0
		for _, v := range x {
			g += (v - 0.5) * (v - 0.5)
		}
		return g
	}
	g := evalG(e.Elements[varSz-k:])
	t := math.Pi / (4.0 * (1 + g))

	newObjs := make([]float64, m)
	theta := make([]float64, m-1)
	theta[0] = e.Elements[0] * math.Pi / 2.0
	for i := 1; i < m-1; i++ {
		theta[i] = t * (1.0 + 2.0*g*e.Elements[i])
	}

	for i := range m {
		prod := (1 + g)
		for j := 0; j < m-(i+1); j++ {
			prod *= math.Cos(theta[j])
		}
		if i != 0 {
			prod *= math.Sin(theta[m-(i+1)])
		}
		newObjs[i] = prod
	}

	// puts new objectives into the elem
	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)

	return nil
}
