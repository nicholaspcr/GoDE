package utils

import (
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

	variantSet = map[string]variants.Interface{
		rand.Rand1().Name():                rand.Rand1(),
		rand.Rand2().Name():                rand.Rand2(),
		best.Best1().Name():                best.Best1(),
		best.Best2().Name():                best.Best2(),
		pbest.Pbest().Name():               pbest.Pbest(),
		currenttobest.CurrToBest1().Name(): currenttobest.CurrToBest1(),
	}
)

// GetProblemByProblemName -> returns the problem function of the given name
func GetProblemByName(name string) problems.Interface {
	name = strings.ToLower(name)
	for k, v := range problemSet {
		if name == k {
			return v
		}
	}
	return nil
}

// GetVariantByVariantName -> returns the variant function of the given name
func GetVariantByName(name string) variants.Interface {
	name = strings.ToLower(name)
	for k, v := range variantSet {
		if name == k {
			return v
		}
	}
	return nil
}

// GetAllProblems -> returns all variants
func GetAllVariants() []variants.Interface {
	out := make([]variants.Interface, len(variantSet))
	index := 0
	for _, v := range variantSet {
		out[index] = v
		index++
	}

	return out
}
