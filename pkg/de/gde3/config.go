package gde3

import "github.com/nicholaspcr/GoDE/pkg/de"

// Constants used for the gde3 algorithm.
type Constants struct {
	DE de.Constants

	CR float64 `json:"cr" yaml:"cr" name:"cr"`
	F  float64 `json:"f"  yaml:"f"  name:"f"`
	P  float64 `json:"p"  yaml:"p"  name:"p"`
}
