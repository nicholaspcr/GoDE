package variants

import "github.com/nicholaspcr/GoDE/pkg/models"

// Parameters are the necessary values that a variant uses
type Parameters struct {
	DIM     int     `json:"dim"      yaml:"dim"`
	CurrPos int     `json:"curr_pos" yaml:"curr_pos"`
	F       float64 `json:"f"        yaml:"f"`
	P       float64 `json:"p"        yaml:"p"`
}

type Interface interface {
	Name() string
	Mutate(elems, rankZero models.Population, params Parameters) (models.Vector, error)
}
