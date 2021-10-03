package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/nicholaspcr/gde3/pkg/models"
	"github.com/nicholaspcr/gde3/pkg/problems/many/dtlz"
)

func main() {
	nVars := 12
	numObjs := 3

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

	dtlz1(x, numObjs)
	dtlz2(x, numObjs)
	dtlz3(x, numObjs)
	dtlz4(x, numObjs)
	dtlz5(x, numObjs)
	dtlz6(x, numObjs)
	dtlz7(x, numObjs)
}

func dtlz1(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	dtlz.DTLZ1.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
func dtlz2(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	dtlz.DTLZ2.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
func dtlz3(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	dtlz.DTLZ3.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
func dtlz4(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	dtlz.DTLZ4.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
func dtlz5(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	dtlz.DTLZ5.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
func dtlz6(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	dtlz.DTLZ6.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
func dtlz7(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Vector{
		X: x,
	}
	dtlz.DTLZ7.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
