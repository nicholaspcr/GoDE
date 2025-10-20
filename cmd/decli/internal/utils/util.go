package utils

import (
	"errors"
	"strings"

	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/problems/many/dtlz"
	"github.com/nicholaspcr/GoDE/pkg/problems/many/wfg"
	"github.com/nicholaspcr/GoDE/pkg/problems/multi"
	"github.com/nicholaspcr/GoDE/pkg/variants"
	"github.com/nicholaspcr/GoDE/pkg/variants/best"
	currenttobest "github.com/nicholaspcr/GoDE/pkg/variants/current-to-best"
	"github.com/nicholaspcr/GoDE/pkg/variants/pbest"
	"github.com/nicholaspcr/GoDE/pkg/variants/rand"
)

var (
	problemSet = map[string]problems.Interface{
		"zdt1": multi.Zdt1(),
		"zdt2": multi.Zdt2(),
		"zdt3": multi.Zdt3(),
		"zdt4": multi.Zdt4(),
		"zdt6": multi.Zdt6(),
		"vnt1": multi.Vnt1(),

		"dtlz1": dtlz.Dtlz1(),
		"dtlz2": dtlz.Dtlz2(),
		"dtlz3": dtlz.Dtlz3(),
		"dtlz4": dtlz.Dtlz4(),
		"dtlz5": dtlz.Dtlz5(),
		"dtlz6": dtlz.Dtlz6(),
		"dtlz7": dtlz.Dtlz7(),

		"wfg1": wfg.Wfg1(),
		"wfg2": wfg.Wfg2(),
		"wfg3": wfg.Wfg3(),
		"wfg4": wfg.Wfg4(),
		"wfg5": wfg.Wfg5(),
		"wfg6": wfg.Wfg6(),
		"wfg7": wfg.Wfg7(),
		"wfg8": wfg.Wfg8(),
		"wfg9": wfg.Wfg9(),
	}

	variantSet = map[string]variants.Interface{
		"rand/1":                rand.Rand1(),
		"rand/2":                rand.Rand2(),
		"best/1":                best.Best1(),
		"best/2":                best.Best2(),
		"pbest/1":               pbest.Pbest(),
		"current-to-best/1": currenttobest.CurrToBest1(),
	}
)

// GetProblemByName -> returns the problem function of the given name
func GetProblemByName(name string) (problems.Interface, error) {
	name = strings.ToLower(name)
	if p, ok := problemSet[name]; ok {
		return p, nil
	}
	return nil, errors.New("problem not found")
}

// GetVariantByName -> returns the variant function of the given name
func GetVariantByName(name string) (variants.Interface, error) {
	name = strings.ToLower(name)
	if v, ok := variantSet[name]; ok {
		return v, nil
	}
	return nil, errors.New("variant not found")
}

// GetAllProblems -> returns all problems
func GetAllProblems() []problems.Interface {
	out := make([]problems.Interface, len(problemSet))
	index := 0
	for _, p := range problemSet {
		out[index] = p
		index++
	}
	return out
}

// GetAllVariants -> returns all variants
func GetAllVariants() []variants.Interface {
	out := make([]variants.Interface, len(variantSet))
	index := 0
	for _, v := range variantSet {
		out[index] = v
		index++
	}

	return out
}
