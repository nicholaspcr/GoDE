package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nicholaspcr/gde3/pkg/problems/many/wfg"
	"github.com/nicholaspcr/gde3/pkg/problems/models"
)

func main() {
	nVars := 12
	numObjs := 3

	// random elements
	rand.Seed(time.Now().Unix())
	var x []float64
	for i := 0; i < nVars; i++ {
		x = append(x, rand.Float64())
	}
	// printing random elements
	fmt.Println(x)

	var elem models.Elem = models.Elem{
		X: x,
	}

	// calling objective function
	wfg1(&elem, numObjs)
	// printing the new objectives
	fmt.Println(elem.Objs)

	// calling objective function
	wfg2(&elem, numObjs)
	// printing the new objectives
	fmt.Println(elem.Objs)

	// calling objective function
	wfg3(&elem, numObjs)
	// printing the new objectives
	fmt.Println(elem.Objs)

	// calling objective function
	wfg4(&elem, numObjs)
	// printing the new objectives
	fmt.Println(elem.Objs)

	// calling objective function
	wfg5(&elem, numObjs)
	// printing the new objectives
	fmt.Println(elem.Objs)

	// calling objective function
	wfg6(&elem, numObjs)
	// printing the new objectives
	fmt.Println(elem.Objs)

	// calling objective function
	wfg7(&elem, numObjs)
	// printing the new objectives
	fmt.Println(elem.Objs)

	// calling objective function
	wfg8(&elem, numObjs)
	// printing the new objectives
	fmt.Println(elem.Objs)

	// calling objective function
	wfg9(&elem, numObjs)
	// printing the new objectives
	fmt.Println(elem.Objs)
}

func wfg1(e *models.Elem, numObjs int) {
	wfg.WFG1.Fn(e, numObjs)
}
func wfg2(e *models.Elem, numObjs int) {
	wfg.WFG2.Fn(e, numObjs)
}
func wfg3(e *models.Elem, numObjs int) {
	wfg.WFG3.Fn(e, numObjs)
}
func wfg4(e *models.Elem, numObjs int) {
	wfg.WFG4.Fn(e, numObjs)
}
func wfg5(e *models.Elem, numObjs int) {
	wfg.WFG5.Fn(e, numObjs)
}
func wfg6(e *models.Elem, numObjs int) {
	wfg.WFG6.Fn(e, numObjs)
}
func wfg7(e *models.Elem, numObjs int) {
	wfg.WFG7.Fn(e, numObjs)
}
func wfg8(e *models.Elem, numObjs int) {
	wfg.WFG8.Fn(e, numObjs)
}
func wfg9(e *models.Elem, numObjs int) {
	wfg.WFG9.Fn(e, numObjs)
}
