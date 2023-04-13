package multi

import (
	"errors"
	"math"

	"github.com/nicholaspcr/GoDE/pkg/api"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type zdt6 struct{}

func Zdt6() problems.Interface {
	return &zdt6{}
}

func (v *zdt6) Name() string {
	return "zdt6"
}

func (v *zdt6) Evaluate(e *api.Vector, M int) error {

	if len(e.Elements) < 2 {
		return errors.New("need at least two variables/dimensions")
	}
	evalF := func(x float64) float64 {
		f := math.Exp(-4.0 * x)
		f = f * math.Pow(math.Sin(6*math.Pi*x), 6)
		f = 1 - f
		return f
	}
	evalG := func(x []float64) float64 {
		g := 0.0
		for i := 1; i < len(x); i++ {
			g += x[i]
		}
		g = g / float64(len(x)-1)
		g = math.Pow(g, 1.0/4)
		g = g*9 + 1.0
		return g
	}
	evalH := func(f, g float64) float64 {
		return 1.0 - math.Pow(f/g, 2)
	}
	F := evalF(e.Elements[0])
	G := evalG(e.Elements)
	H := evalH(F, G)

	var newObjs []float64
	newObjs = append(newObjs, F)
	newObjs = append(newObjs, G*H)

	// puts new objectives into the elem
	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)

	return nil
}
