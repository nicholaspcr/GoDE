package multi

import (
	"errors"
	"math"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type zdt4 struct{}

// Zdt4 returns the ZDT4 test problem, a bi-objective benchmark with many local Pareto fronts.
// Domain: x1 in [0,1], xi in [-5,5], Objectives: 2
func Zdt4() problems.Interface {
	return &zdt4{}
}

func (v *zdt4) Name() string {
	return "zdt4"
}

func (v *zdt4) Evaluate(e *models.Vector, M int) error {
	if len(e.Elements) < 2 {
		return errors.New("need at least two variables/dimensions")
	}
	evalG := func(x []float64) float64 {
		g := 0.0
		for i := 1; i < len(x); i++ {
			g += x[i]*x[i] - 10*math.Cos(4*math.Pi*x[i])
		}
		sz := float64(len(x) - 1)
		return 1.0 + 10.0*sz + g
	}
	evalH := func(f, g float64) float64 {
		return 1.0 - math.Sqrt(f/g)
	}

	g := evalG(e.Elements)
	h := evalH(e.Elements[0], g)

	e.Objectives = []float64{e.Elements[0], g*h}

	return nil
}
