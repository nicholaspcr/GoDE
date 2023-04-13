package multi

import (
	"errors"
	"math"

	"github.com/nicholaspcr/GoDE/pkg/api"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type vnt1 struct{}

// Vnt1 -> https://ti.arc.nasa.gov/m/pub-archive/archive/1163.pdf
func Vnt1() problems.Interface {
	return &vnt1{}
}

func (v *vnt1) Name() string {
	return "vnt1"
}

func (v *vnt1) Evaluate(e *api.Vector, M int) error {

	if len(e.Elements) != 2 {
		return errors.New("need at have only two variables/dimensions")
	}

	a, b := e.Elements[0], e.Elements[1]

	powSum := math.Pow(a, 2) + math.Pow(b, 2)
	f1 := 0.5*(powSum) + math.Sin(powSum)
	f2 := 15.0 + math.Pow(3*a-2*b+4, 2)/8.0 + math.Pow(a-b+1, 2)/27.0
	f3 := -1.1*math.Exp((-1)*powSum) + 1.0/(powSum+1)

	var newObjs []float64
	newObjs = append(newObjs, f1)
	newObjs = append(newObjs, f2)
	newObjs = append(newObjs, f3)

	// puts new objectives into the elem
	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)

	return nil
}
