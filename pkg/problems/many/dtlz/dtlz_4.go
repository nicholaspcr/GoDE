package dtlz

import (
	"errors"
	"math"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type dtlz4 struct{}

func Dtlz4() problems.Interface {
	return &dtlz4{}
}

func (v *dtlz4) Name() string {
	return "dtlz4"
}

func (v *dtlz4) Evaluate(e *models.Vector, M int) error {

	if len(e.X) <= M {
		return errors.New(
			"need to have an M lesser than the amount of variables",
		)
	}
	varSz := len(e.X)
	k := varSz - M + 1
	evalG := func(x []float64) float64 {
		g := 0.0
		for _, v := range x {
			g += (v - 0.5) * (v - 0.5)
		}
		return g
	}
	g := evalG(e.X[varSz-k:])

	newObjs := make([]float64, M)
	for i := 0; i < M; i++ {
		prod := (1 + g)
		for j := 0; j < M-(i+1); j++ {
			prod *= math.Cos(math.Pow(e.X[j], 100) * math.Pi / 2.0)
		}
		if i != 0 {
			prod *= math.Sin(math.Pow(e.X[M-(i+1)], 100) * math.Pi / 2.0)
		}
		newObjs[i] = prod
	}

	// puts new objectives into the elem
	e.Objs = make([]float64, len(newObjs))
	copy(e.Objs, newObjs)

	return nil
}
