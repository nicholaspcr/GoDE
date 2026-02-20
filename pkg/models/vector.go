package models

// Vector is the element of a population in a DE.
type Vector struct {
	Elements         []float64
	Objectives       []float64
	CrowdingDistance float64
}

// Copy all the content of a Vector into a new Vector.
func (v *Vector) Copy() Vector {
	vec := Vector{
		Elements:         make([]float64, len(v.Elements)),
		Objectives:       make([]float64, len(v.Objectives)),
		CrowdingDistance: v.CrowdingDistance,
	}
	copy(vec.Elements, v.Elements)
	copy(vec.Objectives, v.Objectives)
	return vec
}
