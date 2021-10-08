package multi

import (
	"errors"
	"math"

	"github.com/nicholaspcr/gde3/pkg/models"
)

type zdt3 struct{}

func Zdt3() models.Problem {
	return &zdt3{}
}

func (v *zdt3) Name() string {
	return "zdt3"
}

func (v *zdt3) Evaluate(e *models.Vector, M int) error {

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
		return 1.0 - math.Sqrt(f/g) - (f/g)*math.Sin(10.0*f*math.Pi)
	}
	g := evalG(e.X)
	h := evalH(e.X[0], g)
	if math.IsNaN(h) {
		return errors.New("sqrt of a negative number")
	}
	var newObjs []float64
	newObjs = append(newObjs, e.X[0])
	newObjs = append(newObjs, g*h)

	// puts new objectives into the elem
	e.Objs = make([]float64, len(newObjs))
	copy(e.Objs, newObjs)

	return nil
}
