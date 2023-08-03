package variants

import (
	"math/rand"

	"github.com/nicholaspcr/GoDE/pkg/models"
)

// Parameters are the necessary values that a variant uses
type Parameters struct {
	DIM     int     `json:"dim"      yaml:"dim"`
	CurrPos int     `json:"curr_pos" yaml:"curr_pos"`
	F       float64 `json:"f"        yaml:"f"`
	P       float64 `json:"p"        yaml:"p"`

	Random *rand.Rand `json:"-" yaml:"-"`
}

type Interface interface {
	Name() string
	Mutate(
		elems, rankZero []models.Vector,
		params Parameters,
	) (models.Vector, error)
}
