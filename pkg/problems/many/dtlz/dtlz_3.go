package dtlz

import (
	"errors"
	"math"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type dtlz3 struct{}

// Dtlz3 returns the DTLZ3 test problem, similar to DTLZ2 but with many local Pareto fronts.
// Domain: [0,1]^n, Objectives: m (configurable)
func Dtlz3() problems.Interface {
	return &dtlz3{}
}

func (v *dtlz3) Name() string {
	return "dtlz3"
}

func (v *dtlz3) Evaluate(e *models.Vector, m int) error {
	if len(e.Elements) <= m {
		return errors.New(
			"need to have an m lesser than the amount of variables",
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

	g := evalG(e.Elements[m-1:])
	objs := make([]float64, m)

	for i := 0; i < m; i++ {
		prod := (1.0 + g)
		for j := 0; j < m-(i+1); j++ {
			prod *= math.Cos(e.Elements[j] * 0.5 * math.Pi)
		}
		if i != 0 {
			prod *= math.Sin(e.Elements[m-(i+1)] * 0.5 * math.Pi)
		}
		objs[i] = prod
	}

	// puts new objectives into the elem
	e.Objectives = make([]float64, len(objs))
	copy(e.Objectives, objs)
	return nil
}
