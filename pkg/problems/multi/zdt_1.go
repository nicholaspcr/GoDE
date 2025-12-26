// Package multi implements multi-objective test problems including ZDT and VNT benchmark suites.
package multi

import (
	"errors"
	"math"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type zdt1 struct{}

// Zdt1 returns the ZDT1 test problem, a bi-objective benchmark with a convex Pareto front.
// Domain: [0,1]^n, Objectives: 2
func Zdt1() problems.Interface {
	return &zdt1{}
}

func (v *zdt1) Name() string {
	return "zdt1"
}

func (v *zdt1) Evaluate(e *models.Vector, M int) error {
	if len(e.Elements) < 2 {
		return errors.New("need at least two variables/dimensions")
	}
	evalG := func(x []float64) float64 {
		g := 0.0
		for i := 1; i < len(x); i++ {
			g += x[i]
		}
		constant := 9.0 / (float64(len(x)) - 1.0)

		return 1.0 + constant*g
	}
	evalH := func(f, g float64) float64 {
		return 1.0 - math.Sqrt(f/g)
	}
	g := evalG(e.Elements)
	h := evalH(e.Elements[0], g)

	if math.IsNaN(h) {
		return errors.New("sqrt of a negative number")
	}

	e.Objectives = []float64{e.Elements[0], g*h}

	return nil
}
