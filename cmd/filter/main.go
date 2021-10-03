package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/nicholaspcr/gde3/pkg/algorithms"
	"github.com/nicholaspcr/gde3/pkg/models"
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
			c := make(chan models.Population)

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
					tokens <- struct{}{}
					defer func() { <-tokens }()
					defer wg.Done()
					processFile(fileName, c)
				}()

			}
			// waits for all the routines to be done
			go func() {
				wg.Wait()
				close(c)
			}()

			var elems models.Population

			for v := range c {
				elems = append(elems, v...)
			}

			// reduce by crowdDistance
			_, elems = algorithms.ReduceByCrowdDistance(elems, len(elems))

			filePath := fmt.Sprintf(
				"%s/.gode/mode/multiExecutions/%s/%s/filteredPareto.csv",
				dirname,
				prob,
				variant,
			)

			f, err := os.Create(filePath)
			defer func() { f.Close() }()

			if err != nil {
				log.Fatalln(f)
			}

			writer := csv.NewWriter(f)
			writer.Comma = '\t'

			// header
			headerData := []string{"elems"}
			column := 'A'
			for range elems[0].Objs {
				headerData = append(headerData, string(column))
				column++
			}
			err = writer.Write(headerData)
			if err != nil {
				log.Fatal("Coudln't write file")
			}
			writer.Flush()

			// body
			bodyData := [][]string{}
			for i := range elems {
				tmpData := []string{}
				tmpData = append(
					tmpData,
					fmt.Sprintf(
						"elem[%d]",
						i,
					),
				)
				for _, p := range elems[i].Objs {
					tmpData = append(tmpData, fmt.Sprint(p))
				}
				bodyData = append(bodyData, tmpData)
			}
			err = writer.WriteAll(bodyData)
			if err != nil {
				log.Fatalln("failed to write body")
			}
			writer.Flush()
		}
	}
}

func processFile(fileName string, elemChan chan<- models.Population) {

	b, _ := os.ReadFile(fileName)
	reader := csv.NewReader(bytes.NewBuffer(b))
	reader.Comma = '\t'

	var elems models.Population

	// Filling the elems slice
	matStr, _ := reader.ReadAll()
	lines, columns := len(matStr), len(matStr[0])

	gens := lines / numObjs

	for i := 0; i < gens; i++ {
		for j := 0; j < columns; j++ {
			var e models.Vector
			e.Objs = make([]float64, numObjs)

			for k := 0; k < numObjs; k++ {
				value := matStr[i*numObjs+k+1][j]
				f, _ := strconv.ParseFloat(value, 64)
				e.Objs[k] = f
			}

			elems = append(elems, e)
		}
	}

	_, rankZero := algorithms.ReduceByCrowdDistance(elems, len(elems))
	elemChan <- rankZero
}
