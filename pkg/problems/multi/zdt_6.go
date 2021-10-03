package multi

import (
	"errors"
	"math"

	"github.com/nicholaspcr/gde3/pkg/models"
)

// ZDT6 -> bi-objetive evaluation
var ZDT6 = models.Problem{
	Fn: func(e *models.Vector, M int) error {
		if len(e.X) < 2 {
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
		F := evalF(e.X[0])
		G := evalG(e.X)
		H := evalH(F, G)

		var newObjs []float64
		newObjs = append(newObjs, F)
		newObjs = append(newObjs, G*H)

		// puts new objectives into the elem
		e.Objs = make([]float64, len(newObjs))
		copy(e.Objs, newObjs)

		return nil
	},
	ProblemName: "zdt6",
}
