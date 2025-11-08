package dtlz

import (
	"errors"
	"math"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
)

type dtlz7 struct{}

// Dtlz7 returns the DTLZ7 test problem, a many-objective benchmark with disconnected Pareto regions.
// Domain: [0,1]^n, Objectives: m (configurable)
func Dtlz7() problems.Interface {
	return &dtlz7{}
}

func (v *dtlz7) Name() string {
	return "dtlz7"
}

func (v *dtlz7) Evaluate(e *models.Vector, m int) error {
	if len(e.Elements) <= m {
		return errors.New(
			"need to have an m lesser than the amount of variables",
		)
	}
	varSz := len(e.Elements)
	k := varSz - m + 1

	// calculating the value of the constant G
	g := 0.0
	for _, v := range e.Elements[varSz-k:] {
		g += v
	}
	g = 1.0 + (9.0*g)/float64(k)

	// calculating the value of the constant H
	h := 0.0
	for _, v := range e.Elements[:m-1] {
		h += (v / (1.0 + g)) * (1 + math.Sin(3.0*math.Pi*v))
	}
	h = float64(m) - h

	// calculating objs values
	objs := make([]float64, m)
	for i := range objs {
		objs[i] = e.Elements[i]
	}
	objs[m-1] = (1.0 + g) * h
	// puts new objectives into the elem
	e.Objectives = make([]float64, len(objs))
	copy(e.Objectives, objs)

	return nil
}
