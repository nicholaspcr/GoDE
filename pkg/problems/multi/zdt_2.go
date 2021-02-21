package multi

import (
	"errors"
	"math"

	"gitlab.com/nicholaspcr/go-de/pkg/problems/models"
)

// ZDT2 -> bi-objetive evaluation
var ZDT2 = models.ProblemFn{
	Fn: func(e *models.Elem, M int) error {
		if len(e.X) < 2 {
			return errors.New("Need at least two variables/dimensions")
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
	},
	Name: "zdt2",
}
