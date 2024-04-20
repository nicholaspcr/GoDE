package models

import "github.com/nicholaspcr/GoDE/pkg/api/v1"

// Vector is the element of a population in a DE.
type Vector struct {
	Elements         []float64
	Objectives       []float64
	CrowdingDistance float64
}

// ToPB converts a Vector to a protobuf Vector.
func (v *Vector) ToPB() *api.Vector {
	vec := &api.Vector{
		Elements:         make([]float64, len(v.Elements)),
		Objectives:       make([]float64, len(v.Objectives)),
		CrowdingDistance: v.CrowdingDistance,
	}
	copy(vec.Elements, v.Elements)
	copy(vec.Objectives, v.Objectives)
	return vec
}

// VectorFromPB converts a protobuf Vector to a Vector.
func VectorFromPB(v *api.Vector) Vector {
	vec := Vector{
		Elements:         make([]float64, len(v.Elements)),
		Objectives:       make([]float64, len(v.Objectives)),
		CrowdingDistance: v.CrowdingDistance,
	}
	copy(vec.Elements, v.Elements)
	copy(vec.Objectives, v.Objectives)
	return vec
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
