package multi

import (
	"errors"
	"math"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

type zdt2 struct{}

func Zdt2() models.Problem {
	return &zdt2{}
}

func (v *zdt2) Name() string {
	return "zdt2"
}

func (v *zdt2) Evaluate(e *models.Vector, M int) error {

	if len(e.X) < 2 {
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
		return 1.0 - math.Pow(f/g, 2)
	}
	g := evalG(e.X)
	h := evalH(e.X[0], g)

	var newObjs []float64
	newObjs = append(newObjs, e.X[0])
	newObjs = append(newObjs, g*h)

	// puts new objectives into the elem
	e.Objs = make([]float64, len(newObjs))
	copy(e.Objs, newObjs)

	return nil
}
