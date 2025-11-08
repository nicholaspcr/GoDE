package dtlz

import (
	"errors"
	"math"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type dtlz4 struct{}

// Dtlz4 returns the DTLZ4 test problem, similar to DTLZ2 with biased search space density.
// Domain: [0,1]^n, Objectives: m (configurable)
func Dtlz4() problems.Interface {
	return &dtlz4{}
}

func (v *dtlz4) Name() string {
	return "dtlz4"
}

func (v *dtlz4) Evaluate(e *models.Vector, m int) error {
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

	newObjs := make([]float64, m)
	for i := 0; i < m; i++ {
		prod := (1 + g)
		for j := 0; j < m-(i+1); j++ {
			prod *= math.Cos(math.Pow(e.Elements[j], 100) * math.Pi / 2.0)
		}
		if i != 0 {
			prod *= math.Sin(math.Pow(e.Elements[m-(i+1)], 100) * math.Pi / 2.0)
		}
		newObjs[i] = prod
	}

	// puts new objectives into the elem
	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)

	return nil
}
