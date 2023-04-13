package rand

import (
	"fmt"

	"github.com/nicholaspcr/GoDE/pkg/api"
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
	elems, rankZero []api.Vector,
	p variants.Parameters,
) (*api.Vector, error) {

	if len(elems) < 3 {
		return nil, fmt.Errorf(
			"no sufficient amount of elements in the population, should be bigger than three",
		)
	}

	// generating random indices different from current pos
	inds := make([]int, 4)
	inds[0] = p.CurrPos
	err := variants.GenerateIndices(1, len(elems), inds)
	if err != nil {
		return nil, err
	}

	result := &api.Vector{
		Elements: make([]float64, p.DIM),
	}
	for i := 0; i < p.DIM; i++ {
		result.Elements[i] = elems[inds[1]].Elements[i] +
			p.F*(elems[inds[2]].Elements[i]-elems[inds[3]].Elements[i])
	}
	return result, nil
}
