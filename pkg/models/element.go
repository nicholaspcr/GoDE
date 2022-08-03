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
	FloorSlice     []float64
	CeilSlice      []float64
}

func (p *Population) Get(idx int) Vector {
	if p == nil {
		return Vector{}
	}
	return p.Vectors[idx]
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

// Floors contains the floor for each dimention in the vector slice.
func (p *Population) Floors() []float64 {
	if p == nil {
		return nil
	}
	return p.FloorSlice
}

// Ceils contains the ceil for each dimention in the vector slice.
func (p *Population) Ceils() []float64 {
	if p == nil {
		return nil
	}
	return p.CeilSlice
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
