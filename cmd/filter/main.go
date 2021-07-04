package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/nicholaspcr/gde3/pkg/mode"
	"github.com/nicholaspcr/gde3/pkg/problems/models"
)

var (
	// number of objectives in the files being read
	numObjs = 3
	// quantity of executions done
	execs = 30

	// variants being read
	variants = []string{"rand1"}
	// problems being read
	problems = []string{"dtlz1"}

	// tokens is a counting semaphore use to
	// enforce  a limit of 3 concurrent requests
	tokens = make(chan struct{}, 3)
)

func main() {
	dirname, _ := os.UserHomeDir()

	for _, prob := range problems {
		for _, variant := range variants {
			wg := &sync.WaitGroup{}
			c := make(chan models.Elements)

			for exec := 1; exec <= execs; exec++ {
				fileName := fmt.Sprintf(
					"%s/.gode/mode/paretoFront/%s/%s/exec-%d.csv",
					dirname,
					prob,
					variant,
					exec,
				)
				wg.Add(1)

				go func() {
					defer wg.Done()
					processFile(fileName, c)
				}()

			}
			// waits for all the routines to be done
			go func() {
				wg.Wait()
				close(c)
			}()

			for v := range c {
				fmt.Println(len(v))
			}

		}
	}

}

func processFile(fileName string, elemChan chan<- models.Elements) {
	tokens <- struct{}{}
	defer func() { <-tokens }()

	b, _ := os.ReadFile(fileName)
	reader := csv.NewReader(bytes.NewBuffer(b))
	reader.Comma = '\t'

	var elems models.Elements

	// Filling the elems slice
	matStr, _ := reader.ReadAll()
	lines, columns := len(matStr), len(matStr[0])

	gens := lines / numObjs

	for i := 0; i < gens; i++ {
		for j := 0; j < columns; j++ {
			var e models.Elem
			e.Objs = make([]float64, numObjs)

			for k := 0; k < numObjs; k++ {
				value := matStr[i*numObjs+k+1][j]
				f, _ := strconv.ParseFloat(value, 64)
				e.Objs[k] = f
			}

			elems = append(elems, e)
		}
	}

	mode.CalculateCrwdDist(elems)
	_, rankZero := mode.ReduceByCrowdDistance(elems, len(elems))
	elemChan <- rankZero
}
