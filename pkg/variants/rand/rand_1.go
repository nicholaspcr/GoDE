package rand

import (
	"fmt"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// rand1
type rand1 struct{}

// Rand1 variant -> a + F(b - c)
func Rand1() variants.Interface {
	return &rand1{}
}

func (r *rand1) Name() string {
	return "rand1"
}

func (r *rand1) Mutate(
	elems, rankZero []models.Vector,
	p variants.Parameters,
) (models.Vector, error) {

	if len(elems) < 3 {
		return models.Vector{}, fmt.Errorf(
			"no sufficient amount of elements in the population, should be bigger than three",
		)
	}

	// generating random indices different from current pos
	inds := make([]int, 4)
	inds[0] = p.CurrPos
	err := variants.GenerateIndices(1, len(elems), inds)
	if err != nil {
		return models.Vector{}, err
	}

	result := models.Vector{}
	result.X = make([]float64, p.DIM)

	r1, r2, r3 := elems[inds[1]], elems[inds[2]], elems[inds[3]]
	for i := 0; i < p.DIM; i++ {
		result.X[i] = r1.X[i] + p.F*(r2.X[i]-r3.X[i])
	}
	return result, nil
}
