package rand

import (
	"github.com/nicholaspcr/GoDE/pkg/models"
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
	elems, rankZero []models.Vector,
	p variants.Parameters,
) (models.Vector, error) {
	// generating random indices different from current pos
	ind := make([]int, 6)
	ind[0] = p.CurrPos
	err := variants.GenerateIndices(1, len(elems), ind, p.Random)
	if err != nil {
		return models.Vector{}, variants.ErrInsufficientPopulation
	}

	arr := make([]float64, p.DIM)
	i1, i2, i3, i4, i5 := ind[1], ind[2], ind[3], ind[4], ind[5]
	for i := 0; i < p.DIM; i++ {
		arr[i] = elems[i1].Elements[i] +
			p.F*(elems[i2].Elements[i]-elems[i3].Elements[i]) +
			p.F*(elems[i4].Elements[i]-elems[i5].Elements[i])
	}
	ret := models.Vector{
		Elements: arr,
	}
	return ret, nil
}
