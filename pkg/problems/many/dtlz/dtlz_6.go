package dtlz

import (
	"errors"
	"math"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type dtlz6 struct{}

// Dtlz6 returns the DTLZ6 test problem, similar to DTLZ5 with different g function.
// Domain: [0,1]^n, Objectives: m (configurable)
func Dtlz6() problems.Interface {
	return &dtlz6{}
}

func (v *dtlz6) Name() string {
	return "dtlz6"
}

func (v *dtlz6) Evaluate(e *models.Vector, m int) error {
	if len(e.Elements) <= m {
		return errors.New(
			"need to have an m lesser than the amount of variables",
		)
	}

	evalG := func(v []float64) float64 {
		g := 0.0
		for _, x := range v {
			g += math.Pow(x, 0.1) // consumes huge memory
		}
		return g
	}
	g := evalG(e.Elements[m-1:])
	t := math.Pi / (4.0 * (1.0 + g))

	objs := make([]float64, m)
	theta := make([]float64, m-1)

	theta[0] = e.Elements[0] * math.Pi / 2.0
	for i := 1; i < m-1; i++ {
		theta[i] = t * (1.0 + 2.0*g*e.Elements[i])
	}

	for i := 0; i < m; i++ {
		prod := (1 + g)
		for j := 0; j < m-(i+1); j++ {
			prod *= math.Cos(theta[j])
		}
		if i != 0 {
			aux := m - (i + 1)
			prod *= math.Sin(theta[aux])
		}

		objs[i] = prod
	}

	// puts new objectives into the elem
	e.Objectives = make([]float64, len(objs))
	copy(e.Objectives, objs)
	return nil
}
