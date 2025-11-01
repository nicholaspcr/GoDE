package multi

import (
	"errors"
	"math"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type zdt4 struct{}

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

	var newObjs []float64
	newObjs = append(newObjs, e.Elements[0])
	newObjs = append(newObjs, g*h)

	// puts new objectives into the elem
	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)

	return nil
}
