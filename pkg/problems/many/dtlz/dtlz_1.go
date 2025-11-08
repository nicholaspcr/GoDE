// Package dtlz implements the DTLZ many-objective test problem suite.
//
// The implementations are translations of the python code made by pymoo
// https://pymoo.org/problems/many/dtlz.html
package dtlz

import (
	"errors"
	"math"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type dtlz1 struct{}

// Dtlz1 returns the DTLZ1 test problem, a many-objective benchmark with a linear Pareto front.
// Domain: [0,1]^n, Objectives: m (configurable)
func Dtlz1() problems.Interface {
	return &dtlz1{}
}

func (v *dtlz1) Name() string {
	return "dtlz1"
}

func (v *dtlz1) Evaluate(e *models.Vector, m int) error {
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

	newObjs := make([]float64, m)
	for i := 0; i < m; i++ {
		prod := (1.0 + g) * 0.5
		for j := 0; j < m-(i+1); j++ {
			prod *= e.Elements[j]
		}
		if i != 0 {
			prod *= (1 - e.Elements[m-(i+1)])
		}
		newObjs[i] = prod
	}

	// puts new objectives into the elem
	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)

	return nil
}
