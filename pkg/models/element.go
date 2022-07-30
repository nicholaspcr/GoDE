package models

// Vector -> Element of population.
type Vector struct {
	X      []float64 `json:"x"      yaml:"x"`
	Objs   []float64 `json:"objs"   yaml:"objs"`
	Crwdst float64   `json:"crwdst" yaml:"crwdst"`
}

// Copy the entire struct.
func (v *Vector) Copy() Vector {
	var ret = Vector{}

	ret.X = make([]float64, len(v.X))
	copy(ret.X, v.X)

	ret.Objs = make([]float64, len(v.Objs))
	copy(ret.Objs, v.Objs)

	return ret
}

// Population is a slice of the type Vector.
type Population struct {
	Vectors        []Vector `json:"vectors" yaml:"vectors"`
	DimensionsSize int      `json:"dim_size" yaml:"dim_size"`
	ObjectivesSize int      `json:"obj_size" yaml:"obj_size"`
}

// Size of the current Vector slice.
func (p *Population) Size() int {
	if p == nil {
		return 0
	}
	return len(p.Vectors)
}

func (p *Population) DimSize() int {
	if p == nil {
		return 0
	}
	return p.ObjectivesSize
}

// ObjsSize of the current Vector slice.
func (p *Population) ObjSize() int {
	if p == nil {
		return 0
	}
	return p.DimensionsSize
}

// Copy of the Vector slice.
func (p *Population) Copy() Population {
	population := &Population{
		Vectors:        make([]Vector, p.Size()),
		DimensionsSize: p.DimSize(),
		ObjectivesSize: p.ObjSize(),
	}
	copy(population.Vectors, p.Vectors)
	return *population
}
