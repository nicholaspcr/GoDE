package algorithms

import (
	"math/rand"
	"os"

	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/writer"
)

// type that contains the definition of the GDE3 algorithm
type gde3 struct{}

// GDE3 Returns an instance of an object that implements the GDE3 algorithm. It
// is compliant with the Mode
func GDE3() models.Mode {
	return &gde3{}
}

// Execute is responsible for receiving the standard parameters defined in the
// Mode and executing the gde3 algorithm
func (g *gde3) Execute(
	rankedCh chan<- models.Population,
	maximumObjs chan<- []float64,
	p models.AlgorithmParams,
	problem models.Problem,
	variant models.Variant,
	population models.Population,
	f *os.File,
) {
	defer f.Close()

	// var writer *csv.Writer
	w := writer.NewWriter(f)
	w.Comma = ';'

	// maximum objs found
	maxObjs := make([]float64, p.M)

	// calculates the objs of the inital population
	for i := range population {
		err := problem.Evaluate(&population[i], p.M)
		checkError(err)

		for j, obj := range population[i].Objs {
			if obj > maxObjs[j] {
				maxObjs[j] = obj
			}
		}
	}

	// writes the header in this execution's file
	w.WriteHeader(p.M)
	// writes the inital generation
	w.ElementsObjs(population)

	// stores the rank[0] of each generation
	bestElems := make(models.Population, 0)

	var genRankZero models.Population
	var bestInGen models.Population
	var trial models.Vector

	for g := 0; g < p.GEN; g++ {
		// gets non dominated of the current population
		genRankZero, _ = FilterDominated(population)

		for i := 0; i < len(population); i++ {

			// generates the mutatated vector
			vr, err := variant.Mutate(
				population,
				genRankZero,
				models.VariantParams{
					DIM:     p.DIM,
					F:       p.F,
					CurrPos: i,
					P:       p.P,
				})
			checkError(err)

			// trial element
			trial = population[i].Copy()

			// CROSS OVER
			currInd := rand.Int() % p.DIM
			luckyIndex := rand.Int() % p.DIM

			for j := 0; j < p.DIM; j++ {
				changeProb := rand.Float64()
				if changeProb < p.CR || currInd == luckyIndex {
					trial.X[currInd] = vr.X[currInd]
				}

				if trial.X[currInd] < p.FLOOR[currInd] {
					trial.X[currInd] = p.FLOOR[currInd]
				}
				if trial.X[currInd] > p.CEIL[currInd] {
					trial.X[currInd] = p.CEIL[currInd]
				}
				currInd = (currInd + 1) % p.DIM
			}

			evalErr := problem.Evaluate(&trial, p.M)
			checkError(evalErr)

			// SELECTION
			comp := DominanceTest(population[i].Objs, trial.Objs)
			if comp == 1 {
				population[i] = trial.Copy()
			} else if comp == 0 && len(population) <= 2*p.NP {
				population = append(population, trial.Copy())
			}
		}

		population, bestInGen = ReduceByCrowdDistance(population, p.NP)
		bestElems = append(bestElems, bestInGen...)

		// writes the objectives of the population
		w.ElementsObjs(population)

		// checks for the biggest objective
		for _, p := range population {
			for j, obj := range p.Objs {
				if obj > maxObjs[j] {
					maxObjs[j] = obj
				}
			}
		}
	}

	// sending via channel the data
	rankedCh <- bestElems
	maximumObjs <- maxObjs
}
