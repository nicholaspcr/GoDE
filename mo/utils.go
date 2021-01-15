package mo

import (
	"errors"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
)

// GetProblemByName -> returns the problem function
func GetProblemByName(name string) ProblemFn {
	name = strings.ToLower(name)
	problems := map[string]ProblemFn{
		"zdt1":  zdt1,
		"zdt2":  zdt2,
		"zdt3":  zdt3,
		"zdt4":  zdt4,
		"zdt6":  zdt6,
		"vnt1":  vnt1,
		"dtlz1": dtlz1,
		"dtlz2": dtlz2,
		"dtlz3": dtlz3,
		"dtlz4": dtlz4,
		"dtlz5": dtlz5,
		"dtlz6": dtlz6,
		"dtlz7": dtlz7,
	}
	var problem ProblemFn
	for k, v := range problems {
		if name == k {
			problem = v
			break
		}
	}
	return problem
}

// GetVariantByName -> Returns the variant function
func GetVariantByName(name string) VariantFn {
	name = strings.ToLower(name)
	variants := map[string]VariantFn{
		"rand1": rand1,
	}
	for k, v := range variants {
		if name == k {
			return v
		}
	}
	return VariantFn{}
}

func generatePopulation(p Params) Elements {
	ret := make(Elements, p.NP)
	constant := p.CEIL - p.FLOOR // range between floor and ceiling
	for i := 0; i < p.NP; i++ {
		ret[i].X = make([]float64, p.DIM)

		for j := 0; j < p.DIM; j++ {
			ret[i].X[j] = rand.Float64()*constant + p.FLOOR // value varies within [ceil,upper]
		}
	}
	return ret
}

// generates random indices in the int slice, r -> it's a pointer
func generateIndices(startInd, NP int, r []int) error {
	if len(r) > NP {
		return errors.New("insufficient elements in population to generate random indices")
	}
	for i := startInd; i < len(r); i++ {
		for done := false; !done; {
			r[i] = rand.Int() % NP
			done = true
			for j := 0; j < i; j++ {
				done = done && r[j] != r[i]
			}
		}
	}
	return nil
}

// checks existance of filePath
func checkFilePath(basePath, filePath string) {
	folders := strings.Split(filePath, "/")
	for _, folder := range folders {
		basePath += "/" + folder
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			err = os.Mkdir(basePath, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// todo: create a proper error handler
func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// returns NP elements filtered by rank and crwod distance
func reduceByCrowdDistance(elems, best *Elements, NP int) Elements {
	ranks := rankElements(*elems)
	*elems = make(Elements, 0)
	//sorting each rank by crowd distance
	for i := range ranks {
		calculateCrwdDist(ranks[i])
		sort.Sort(byCrwdst(ranks[i]))
	}
	// writes the best ranked into the pareto db
	*best = append(*best, ranks[0]...)

	for _, rank := range ranks {
		for _, elem := range rank {
			*elems = append(*elems, elem.Copy())
			if len(*elems) >= NP {
				return *elems
			}
		}
	}
	return *elems
}

// rankElements returna  map of dominating elements in ascending order
// destroys the slice provided
func rankElements(elems Elements) map[int]Elements {
	ranks := make(map[int]Elements)
	currentRank := 0
	for len(elems) > 0 {
		ranked, nonRanked := filterDominated(elems)
		ranks[currentRank] = append(ranks[currentRank], ranked...)
		currentRank++
		elems = make(Elements, 0)
		elems = append(elems, nonRanked...)
	}
	return ranks
}

// assumes that the slice is composed of non dominated elements
func calculateCrwdDist(elems Elements) {
	if len(elems) <= 3 {
		return
	}
	objs := len(elems[0].objs)
	maxes := make([]float64, len(elems))
	minis := make([]float64, len(elems))
	for i := range elems {
		// resets the crwdst
		elems[i].crwdst = 0
		// gets the max/min values of each objective
		for j := 0; j < objs; j++ {
			maxes[j] = math.Max(maxes[j], elems[i].objs[j])
			minis[j] = math.Min(minis[j], elems[i].objs[j])
		}
	}
	for m := 0; m < objs; m++ {
		// sort by current objective
		sort.SliceStable(elems, func(i, j int) bool {
			return elems[i].objs[m] < elems[j].objs[m]
		})
		// adds an advantage to the points in the extreme
		elems[0].crwdst = maxes[m]
		elems[len(elems)-1].crwdst = maxes[m]
		for i := 1; i < len(elems)-1; i++ {
			elems[i].crwdst = elems[i].crwdst + (elems[i+1].objs[m]-elems[i-1].objs[m])/(maxes[m]-minis[m])
		}
	}
}

// filterDominated -> returns elements that are not dominated in the set
func filterDominated(elems Elements) (nonDominated, dominated Elements) {
	sort.Sort(byFirstObj(elems))
	nonDom := make(Elements, 0)
	dom := make(Elements, 0)
	for i := len(elems) - 1; i >= 0; i-- {
		flag := true
		for j, second := range elems {
			if i == j {
				continue
			}
			if second.dominates(elems[i]) {
				flag = false
				break
			}
		}
		if flag {
			nonDom = append(nonDom, elems[i])
		} else {
			dom = append(dom, elems[i])
		}
	}
	return nonDom, dom
}
