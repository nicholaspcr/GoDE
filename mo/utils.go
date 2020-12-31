package mo

import (
	"errors"
	"fmt"
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
	var variant VariantFn
	for k, v := range variants {
		if name == k {
			variant = v
			break
		}
	}
	return variant
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

// checks if
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

// todo: maybe remove this and do a separate subcommand to write the result in a .csv file!
func writeHeader(pop []Elem, f *os.File) {
	for i := range pop {
		fmt.Fprintf(f, "pop[%d]\t", i)
	}
	fmt.Fprintf(f, "\n")
}

// todo: maybe remove this and do a separate subcommand to write the result in a .csv file!
func writeGeneration(pop Elements, f *os.File) {
	qtdObjs := len(pop[0].objs)
	for i := 0; i < qtdObjs; i++ {
		for _, p := range pop {
			fmt.Fprintf(f, "%10.3f\t", p.objs[i])
		}
		fmt.Fprintf(f, "\n")
	}
}

// returns NP elements filtered by rank and crwod distance
func reduceByCrowdDistance(elems Elements, NP int) Elements {
	ranks := rankElements(elems)
	elems = make(Elements, 0)
	//sorting each rank by crowd distance
	for i := range ranks {
		calcCrwdst(ranks[i])
		sort.Sort(byCrwdst(ranks[i]))
	}
	for _, rank := range ranks {
		for _, elem := range rank {
			elems = append(elems, elem.Copy())
			if len(elems) >= NP {
				return elems
			}
		}
	}
	return elems
}

// rankElements returna  map of dominating elements in ascending order
// destroys the slice provided
func rankElements(elems Elements) map[int]Elements {
	ranks := make(map[int]Elements)
	currentRank := 0
	for len(elems) > 0 {
		inRank := make(Elements, 0)
		notInRank := make(Elements, 0)
		for i := range elems {
			flag := true
			for j := range elems {
				if i == j {
					continue
				}
				if elems[j].dominates(elems[i]) {
					flag = false
					break
				}
			}
			if flag == true {
				inRank = append(inRank, elems[i].Copy())
			} else {
				notInRank = append(notInRank, elems[i].Copy())
			}
		}
		ranks[currentRank] = append(ranks[currentRank], inRank...)
		currentRank++
		elems = make(Elements, 0)
		elems = append(elems, notInRank...)
	}
	return ranks
}

// calcCrwdst -> calculates the crowd distance between elements
func calcCrwdst(elems Elements) {
	if len(elems) <= 3 {
		return
	}
	sort.Sort(byFirstObj(elems))
	for i := range elems {
		for j := range elems[i].objs {
			//end of tail
			if i == 0 {
				elems[i].crwdst += (elems[i+1].objs[j] - elems[i].objs[j]) * (elems[i+1].objs[j] - elems[i].objs[j])
			} else if i == len(elems)-1 {
				elems[i].crwdst += (elems[i-1].objs[j] - elems[i].objs[j]) * (elems[i-1].objs[j] - elems[i].objs[j])
			} else {
				elems[i].crwdst += (elems[i-1].objs[j] - elems[i+1].objs[j]) * (elems[i-1].objs[j] - elems[i+1].objs[j])
			}
			elems[i].crwdst = math.Sqrt(elems[i].crwdst)
		}
	}
}

// filterDominated -> returns elements that are not dominated in the set
func filterDominated(elems Elements) Elements {
	result := make(Elements, 0)
	for i, first := range elems {
		flag := true
		for j, second := range elems {
			if i == j {
				continue
			}
			if second.dominates(first) {
				flag = false
				break
			}
		}
		if flag {
			result = append(result, first)
		}
	}
	return result
}
