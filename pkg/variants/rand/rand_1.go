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
	err := variants.GenerateIndices(1, len(elems), inds, p.Random)
	if err != nil {
		return models.Vector{}, err
	}

	result := models.Vector{
		Elements: make([]float64, p.DIM),
	}
	for i := 0; i < p.DIM; i++ {
		result.Elements[i] = elems[inds[1]].Elements[i] +
			p.F*(elems[inds[2]].Elements[i]-elems[inds[3]].Elements[i])
	}
	return result, nil
}
