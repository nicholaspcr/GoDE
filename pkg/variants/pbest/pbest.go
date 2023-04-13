package pbest

import (
	"errors"
	"math"
	"math/rand"

	"github.com/nicholaspcr/GoDE/pkg/api"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// pbest
type pbest struct{}

func Pbest() variants.Interface {
	return &pbest{}
}

func (p *pbest) Name() string {
	return "pbest"
}

func (p *pbest) Mutate(
	elems, rankZero []api.Vector,
	params variants.Parameters,
) (*api.Vector, error) {

	ind := make([]int, 3)
	ind[0] = params.CurrPos

	err := variants.GenerateIndices(1, len(elems), ind)
	if err != nil {
		return nil, errors.New(
			"insufficient size for the population, must me equal or greater than 5",
		)
	}

	indexLimit := int(math.Ceil(float64(len(rankZero)) * params.P))
	bestIndex := rand.Int() % indexLimit

	arr := make([]float64, params.DIM)
	for i := 0; i < params.DIM; i++ {
		arr[i] = elems[params.CurrPos].Elements[i] +
			params.F*(elems[bestIndex].Elements[i]-elems[params.CurrPos].Elements[i]) +
			params.F*(elems[ind[1]].Elements[i]-elems[ind[2]].Elements[i])
	}

	ret := &api.Vector{
		Elements: arr,
	}
	return ret, nil
}
