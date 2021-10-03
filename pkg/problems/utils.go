package problems

import (
	"strings"

	"github.com/nicholaspcr/gde3/pkg/models"
	"github.com/nicholaspcr/gde3/pkg/problems/many/dtlz"
	"github.com/nicholaspcr/gde3/pkg/problems/many/wfg"
	"github.com/nicholaspcr/gde3/pkg/problems/multi"
)

var problems = map[string]*models.Problem{
	multi.ZDT1.ProblemName: &multi.ZDT1,
	multi.ZDT2.ProblemName: &multi.ZDT2,
	multi.ZDT3.ProblemName: &multi.ZDT3,
	multi.ZDT4.ProblemName: &multi.ZDT4,
	multi.ZDT6.ProblemName: &multi.ZDT6,
	multi.VNT1.ProblemName: &multi.VNT1,

	dtlz.DTLZ1.ProblemName: &dtlz.DTLZ1,
	dtlz.DTLZ2.ProblemName: &dtlz.DTLZ2,
	dtlz.DTLZ3.ProblemName: &dtlz.DTLZ3,
	dtlz.DTLZ4.ProblemName: &dtlz.DTLZ4,
	dtlz.DTLZ5.ProblemName: &dtlz.DTLZ5,
	dtlz.DTLZ6.ProblemName: &dtlz.DTLZ6,
	dtlz.DTLZ7.ProblemName: &dtlz.DTLZ7,

	wfg.WFG1.ProblemName: &wfg.WFG1,
	wfg.WFG2.ProblemName: &wfg.WFG2,
	wfg.WFG3.ProblemName: &wfg.WFG3,
	wfg.WFG4.ProblemName: &wfg.WFG4,
	wfg.WFG5.ProblemName: &wfg.WFG5,
	wfg.WFG6.ProblemName: &wfg.WFG6,
	wfg.WFG7.ProblemName: &wfg.WFG7,
	wfg.WFG8.ProblemName: &wfg.WFG8,
	wfg.WFG9.ProblemName: &wfg.WFG9,
}

// GetProblemByProblemName -> returns the problem function
func GetProblemByName(name string) models.ProblemInterface {
	name = strings.ToLower(name)

	for k, v := range problems {
		if name == k {
			return v
		}
	}

	return &models.Problem{}
}
