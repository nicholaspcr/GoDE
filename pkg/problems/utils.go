package problems

import (
	"strings"

	"github.com/nicholaspcr/gde3/pkg/problems/many/dtlz"
	"github.com/nicholaspcr/gde3/pkg/problems/many/wfg"
	"github.com/nicholaspcr/gde3/pkg/problems/models"
	"github.com/nicholaspcr/gde3/pkg/problems/multi"
)

// GetProblemByName -> returns the problem function
func GetProblemByName(Name string) models.ProblemFn {
	Name = strings.ToLower(Name)
	problems := map[string]models.ProblemFn{
		multi.ZDT1.Name: multi.ZDT1,
		multi.ZDT2.Name: multi.ZDT2,
		multi.ZDT3.Name: multi.ZDT3,
		multi.ZDT4.Name: multi.ZDT4,
		multi.ZDT6.Name: multi.ZDT6,
		multi.VNT1.Name: multi.VNT1,

		dtlz.DTLZ1.Name: dtlz.DTLZ1,
		dtlz.DTLZ2.Name: dtlz.DTLZ2,
		dtlz.DTLZ3.Name: dtlz.DTLZ3,
		dtlz.DTLZ4.Name: dtlz.DTLZ4,
		dtlz.DTLZ5.Name: dtlz.DTLZ5,
		dtlz.DTLZ6.Name: dtlz.DTLZ6,
		dtlz.DTLZ7.Name: dtlz.DTLZ7,

		wfg.WFG1.Name: wfg.WFG1,
		wfg.WFG2.Name: wfg.WFG2,
		wfg.WFG3.Name: wfg.WFG3,
		wfg.WFG4.Name: wfg.WFG4,
		wfg.WFG5.Name: wfg.WFG5,
		wfg.WFG6.Name: wfg.WFG6,
		wfg.WFG7.Name: wfg.WFG7,
		wfg.WFG8.Name: wfg.WFG8,
		wfg.WFG9.Name: wfg.WFG9,
	}
	var problem models.ProblemFn
	for k, v := range problems {
		if Name == k {
			problem = v
			break
		}
	}
	return problem
}
