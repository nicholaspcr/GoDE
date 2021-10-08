package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/nicholaspcr/gde3/pkg/models"
	"github.com/nicholaspcr/gde3/pkg/problems/many/wfg"
)

func main() {
	nVars := 24
	numObjs := 3

	// random elements
	rand.Seed(time.Now().UnixNano())
	var x []float64
	for i := 0; i < nVars; i++ {
		x = append(x, rand.Float64())
	}
	// printing random elements
	var str string
	for _, v := range x {
		str += fmt.Sprintf("%v,", v)
	}
	str = strings.TrimSuffix(str, ",")
	fmt.Println("Array:")
	fmt.Printf("[%v]\n", str)
	fmt.Println("Results:")

	// calling objective function
	wfg1(x, numObjs)

	// calling objective function
	wfg2(x, numObjs)

	// calling objective function
	wfg3(x, numObjs)

	// calling objective function
	wfg4(x, numObjs)

	// calling objective function
	wfg5(x, numObjs)

	// calling objective function
	wfg6(x, numObjs)

	// calling objective function
	wfg7(x, numObjs)

	// calling objective function
	wfg8(x, numObjs)

	// calling objective function
	wfg9(x, numObjs)
}

func wfg1(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	wfg.Wfg1().Evaluate(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg2(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	wfg.Wfg2().Evaluate(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg3(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	wfg.Wfg3().Evaluate(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg4(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	wfg.Wfg4().Evaluate(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg5(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	wfg.Wfg5().Evaluate(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg6(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	wfg.Wfg6().Evaluate(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg7(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	wfg.Wfg7().Evaluate(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg8(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	wfg.Wfg8().Evaluate(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg9(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	wfg.Wfg9().Evaluate(&e, numObjs)

	fmt.Println(e.Objs)
}
