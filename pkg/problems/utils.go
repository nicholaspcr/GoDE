package problems

import (
	"strings"

	"gitlab.com/nicholaspcr/go-de/pkg/problems/many"
	"gitlab.com/nicholaspcr/go-de/pkg/problems/models"
	"gitlab.com/nicholaspcr/go-de/pkg/problems/multi"
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

		many.DTLZ1.Name: many.DTLZ1,
		many.DTLZ2.Name: many.DTLZ2,
		many.DTLZ3.Name: many.DTLZ3,
		many.DTLZ4.Name: many.DTLZ4,
		many.DTLZ5.Name: many.DTLZ5,
		many.DTLZ6.Name: many.DTLZ6,
		many.DTLZ7.Name: many.DTLZ7,
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
