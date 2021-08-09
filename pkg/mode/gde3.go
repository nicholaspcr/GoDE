package mode

import (
	"math/rand"
	"os"

	"github.com/nicholaspcr/gde3/pkg/problems/models"
	"github.com/nicholaspcr/gde3/pkg/variants"
	"github.com/nicholaspcr/gde3/pkg/writer"
)

// GD3 -> runs a simple multiObjective DE in the ZDT1 case
func GD3(
	rankedCh chan<- models.Elements,
	maximumObjs chan<- []float64,
	p models.Params,
	evaluate func(e *models.Elem, M int) error,
	variant variants.VariantFn,
	population models.Elements,
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
		err := evaluate(&population[i], p.M)
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
	bestElems := make(models.Elements, 0)

	var genRankZero models.Elements
	var bestInGen models.Elements
	var trial models.Elem

	for g := 0; g < p.GEN; g++ {
		// gets non dominated of the current population
		genRankZero, _ = FilterDominated(population)

		for i := 0; i < len(population); i++ {

			// generates the mutatated vector
			vr, err := variant.Fn(
				population,
				genRankZero,
				variants.Params{
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

			evalErr := evaluate(&trial, p.M)
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
