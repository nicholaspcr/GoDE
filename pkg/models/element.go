package models

type (
	// Vector -> Element of population
	Vector struct {
		X      []float64 `json:"x"      yaml:"x"`
		Objs   []float64 `json:"objs"   yaml:"objs"`
		Crwdst float64   `json:"crwdst" yaml:"crwdst"`
	}

	// Population is a slice of the type Vector
	Population []Vector
)

// Copy the entire struct
func (e *Vector) Copy() Vector {
	var ret = Vector{}

	ret.X = make([]float64, len(e.X))
	copy(ret.X, e.X)

	ret.Objs = make([]float64, len(e.Objs))
	copy(ret.Objs, e.Objs)

	return ret
}

// Copy of the []Vector slice
func (e Population) Copy() Population {
	arr := make(Population, len(e))
	for i, v := range e {
		arr[i] = v.Copy()
	}
	return arr
}
