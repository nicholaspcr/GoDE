package gde3

import (
	"errors"

	"github.com/nicholaspcr/GoDE/pkg/de"
)

// Constants used for the gde3 algorithm.
type Constants struct {
	DE de.Constants

	CR float64 `json:"cr" yaml:"cr" name:"cr"`
	F  float64 `json:"f"  yaml:"f"  name:"f"`
	P  float64 `json:"p"  yaml:"p"  name:"p"`
}

// Validate checks that all GDE3 Constants fields have valid values.
func (c *Constants) Validate() error {
	if err := c.DE.Validate(); err != nil {
		return err
	}
	if c.CR < 0 || c.CR > 1 {
		return errors.New("CR (crossover rate) must be in range [0, 1]")
	}
	if c.F <= 0 || c.F > 2 {
		return errors.New("F (scaling factor) must be in range (0, 2]")
	}
	if c.P <= 0 || c.P > 1 {
		return errors.New("P (selection parameter) must be in range (0, 1]")
	}
	return nil
}
