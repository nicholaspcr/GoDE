package dtlz

import (
	"errors"
	"math"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type dtlz2 struct{}

func Dtlz2() problems.Interface {
	return &dtlz2{}
}

func (v *dtlz2) Name() string {
	return "dtlz2"
}

func (v *dtlz2) Evaluate(e *models.Vector, M int) error {
	if len(e.Elements) <= M {
		return errors.New(
			"need to have an M lesser than the amount of variables",
		)
	}

	varSz := len(e.Elements)
	k := varSz - M + 1
	evalG := func(x []float64) float64 {
		g := 0.0
		for _, v := range x {
			g += (v - 0.5) * (v - 0.5)
		}
		return g
	}
	g := evalG(e.Elements[varSz-k:])

	newObjs := make([]float64, M)
	for i := 0; i < M; i++ {
		prod := (1 + g)
		for j := 0; j < M-(i+1); j++ {
			prod *= math.Cos(e.Elements[j] * 0.5 * math.Pi)
		}
		if i != 0 {
			prod *= math.Sin(0.5 * math.Pi * e.Elements[M-(i+1)])
		}
		newObjs[i] = prod
	}

	// puts new objectives into the elem
	e.Objectives = make([]float64, len(newObjs))
	copy(e.Objectives, newObjs)

	return nil
}
