package problems

import (
	"strings"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems/many/dtlz"
	"github.com/nicholaspcr/GoDE/pkg/problems/many/wfg"
	"github.com/nicholaspcr/GoDE/pkg/problems/multi"
)

var problems = map[string]models.Problem{
	multi.Zdt1().Name(): multi.Zdt1(),
	multi.Zdt2().Name(): multi.Zdt2(),
	multi.Zdt3().Name(): multi.Zdt3(),
	multi.Zdt4().Name(): multi.Zdt4(),
	multi.Zdt6().Name(): multi.Zdt6(),
	multi.Vnt1().Name(): multi.Vnt1(),

	dtlz.Dtlz1().Name(): dtlz.Dtlz1(),
	dtlz.Dtlz2().Name(): dtlz.Dtlz2(),
	dtlz.Dtlz3().Name(): dtlz.Dtlz3(),
	dtlz.Dtlz4().Name(): dtlz.Dtlz4(),
	dtlz.Dtlz5().Name(): dtlz.Dtlz5(),
	dtlz.Dtlz6().Name(): dtlz.Dtlz6(),
	dtlz.Dtlz7().Name(): dtlz.Dtlz7(),

	wfg.Wfg1().Name(): wfg.Wfg1(),
	wfg.Wfg2().Name(): wfg.Wfg2(),
	wfg.Wfg3().Name(): wfg.Wfg3(),
	wfg.Wfg4().Name(): wfg.Wfg4(),
	wfg.Wfg5().Name(): wfg.Wfg5(),
	wfg.Wfg6().Name(): wfg.Wfg6(),
	wfg.Wfg7().Name(): wfg.Wfg7(),
	wfg.Wfg8().Name(): wfg.Wfg8(),
	wfg.Wfg9().Name(): wfg.Wfg9(),
}

// GetProblemByProblemName -> returns the problem function
func GetProblemByName(name string) models.Problem {
	name = strings.ToLower(name)

	for k, v := range problems {
		if name == k {
			return v
		}
	}

	return nil
}
