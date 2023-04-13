package rand

import (
	"errors"

	"github.com/nicholaspcr/GoDE/pkg/api"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// rand2 a + F(b-c) + F(d-e)
type rand2 struct{}

func Rand2() variants.Interface {
	return &rand2{}
}

func (r *rand2) Name() string {
	return "rand2"
}

func (r *rand2) Mutate(
	elems, rankZero []api.Vector,
	p variants.Parameters,
) (*api.Vector, error) {
	// generating random indices different from current pos
	ind := make([]int, 6)
	ind[0] = p.CurrPos
	err := variants.GenerateIndices(1, len(elems), ind)
	if err != nil {
		return nil, errors.New(
			"insufficient size for the population, must me equal or greater than 4",
		)
	}

	arr := make([]float64, p.DIM)
	i1, i2, i3, i4, i5 := ind[1], ind[2], ind[3], ind[4], ind[5]
	for i := 0; i < p.DIM; i++ {
		arr[i] = elems[i1].Elements[i] +
			p.F*(elems[i2].Elements[i]-elems[i3].Elements[i]) +
			p.F*(elems[i4].Elements[i]-elems[i5].Elements[i])
	}
	ret := &api.Vector{
		Elements: arr,
	}
	return ret, nil
}
