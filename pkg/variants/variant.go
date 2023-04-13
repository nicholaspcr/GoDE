package variants

import (
	"github.com/nicholaspcr/GoDE/pkg/api"
)

// Parameters are the necessary values that a variant uses
type Parameters struct {
	DIM     int     `json:"dim"      yaml:"dim"`
	CurrPos int     `json:"curr_pos" yaml:"curr_pos"`
	F       float64 `json:"f"        yaml:"f"`
	P       float64 `json:"p"        yaml:"p"`
}

type Interface interface {
	Name() string
	Mutate(
		elems, rankZero []api.Vector,
		params Parameters,
	) (*api.Vector, error)
}
