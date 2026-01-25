package multi

import (
	"errors"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type zdt2 struct{}

// Zdt2 returns the ZDT2 test problem, a bi-objective benchmark with a non-convex Pareto front.
// Domain: [0,1]^n, Objectives: 2
func Zdt2() problems.Interface {
	return &zdt2{}
}

func (v *zdt2) Name() string {
	return "zdt2"
}

func (v *zdt2) Evaluate(e *models.Vector, M int) error {
	if len(e.Elements) < 2 {
		return errors.New("need at least two variables/dimensions")
	}
	evalG := func(x []float64) float64 {
		g := 0.0
		for i := 1; i < len(x); i++ {
			g += x[i]
		}
		constant := (9.0 / (float64(len(x)) - 1.0))

		return 1.0 + constant*g
	}
	evalH := func(f, g float64) float64 {
		return 1.0 - f/g*(f/g)
	}
	g := evalG(e.Elements)
	h := evalH(e.Elements[0], g)

	e.Objectives = []float64{e.Elements[0], g * h}

	return nil
}
