package gde3

import (
	"context"
	"math/rand"

	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/store"
	"github.com/nicholaspcr/GoDE/pkg/variants"
)

// gde3 type that contains the definition of the GDE3 algorithm.
type gde3 struct{}

// GDE3 Returns an instance of an object that implements the GDE3
// algorithm. It is compliant with the Mode
func New() de.Algorithm {
	return &gde3{}
}

// Execute is responsible for receiving the standard parameters defined
// in the Mode and executing the gde3 algorithm
func (g *gde3) Execute(
	ctx context.Context,
	population models.Population,
	problem problems.Interface,
	variant variants.Interface,
	store store.Store,
	pareto chan<- []models.Vector,
) error {
	GEN := de.FetchGenerations(ctx)
	dimSize := population.ObjSize()
	// maximum objs found
	maxObjs := make([]float64, dimSize)
	// calculates the objs of the inital population
	for i := range population.Vectors {
		err := problem.Evaluate(&population.Vectors[i], dimSize)
		if err != nil {
			return err
		}
		for j, obj := range population.Vectors[i].Objs {
			if obj > maxObjs[j] {
				maxObjs[j] = obj
			}
		}
	}

	//// writes the header in this execution's file
	//if err := store.Header(); err != nil {
	//	// TODO: Add header contents to methods
	//	return err
	//}
	// writes the inital generation
	if err := store.Population(population); err != nil {
		panic(err)
	}

	// stores the rank[0] of each generation
	bestElems := make([]models.Vector, population.DimSize())

	var genRankZero []models.Vector
	var bestInGen []models.Vector
	var trial models.Vector

	for g := 0; g < GEN; g++ {
		// gets non dominated of the current population
		genRankZero, _ = de.FilterDominated(population.Vectors)

		for i := 0; i < len(population.Vectors); i++ {
			// generates the mutatated vector
			vr, err := variant.Mutate(
				population.Vectors,
				genRankZero,
				variants.Parameters{
					DIM:     population.DimSize(),
					F:       de.FetchFConst(ctx),
					CurrPos: i,
					P:       de.FetchPConst(ctx),
				})
			if err != nil {
				return err
			}

			// trial element
			trial = population.Vectors[i].Copy()

			// CROSS OVER
			currInd := rand.Int() % population.DimSize()
			luckyIndex := rand.Int() % population.DimSize()

			for j := 0; j < population.DimSize(); j++ {
				changeProb := rand.Float64()
				if changeProb < de.FetchCRConst(ctx) || currInd == luckyIndex {
					trial.X[currInd] = vr.X[currInd]
				}

				if trial.X[currInd] < population.Floors()[currInd] {
					trial.X[currInd] = population.Floors()[currInd]
				}
				if trial.X[currInd] > population.Ceils()[currInd] {
					trial.X[currInd] = population.Ceils()[currInd]
				}
				currInd = (currInd + 1) % population.DimSize()
			}

			if err := problem.Evaluate(&trial, dimSize); err != nil {
				return err
			}

			// SELECTION
			comp := de.DominanceTest(population.Vectors[i].Objs, trial.Objs)
			if comp == 1 {
				population.Vectors[i] = trial.Copy()
			} else if comp == 0 && len(population.Vectors) <= 2*population.Size() {
				population.Vectors = append(population.Vectors, trial.Copy())
			}
		}

		population.Vectors, bestInGen = de.ReduceByCrowdDistance(
			population.Vectors,
			population.Size(),
		)
		bestElems = append(bestElems, bestInGen...)

		// writes the objectives of the population
		if err := store.Population(population); err != nil {
			return err
		}

		//// checks for the biggest objective
		//for _, vector := range population.Vectors {
		//	for j, obj := range vector.Objs {
		//		if obj > maxObjs[j] {
		//			maxObjs[j] = obj
		//		}
		//	}
		//}
	}

	// sending via channel the data
	// maximumObjs <- maxObjs
	pareto <- bestElems
	return nil
}
