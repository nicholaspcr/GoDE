package multi

import (
	"errors"
	"math"

	"github.com/nicholaspcr/gde3/pkg/models"
)

// ZDT4 -> bi-objetive evaluation
var ZDT4 = models.Problem{
	Fn: func(e *models.Vector, M int) error {
		if len(e.X) < 2 {
			return errors.New("need at least two variables/dimensions")
		}
		evalG := func(x []float64) float64 {
			g := 0.0
			for i := 1; i < len(x); i++ {
				g += math.Pow(x[i], 2) - 10*math.Cos(4*math.Pi*x[i])
			}
			sz := float64(len(x) - 1)
			return 1.0 + 10.0*sz + g
		}
		evalH := func(f, g float64) float64 {
			return 1.0 - math.Sqrt(f/g)
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
	ProblemName: "zdt4",
}
