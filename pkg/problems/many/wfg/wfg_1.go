package wfg

import (
	"gitlab.com/nicholaspcr/go-de/pkg/problems/models"
)

var WFG1 = models.ProblemFn{
	Fn: func(e *models.Elem, M int) error {
		return nil
	},
	Name: "WFG1",
}
