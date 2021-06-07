package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nicholaspcr/gde3/pkg/problems/many/wfg"
	"github.com/nicholaspcr/gde3/pkg/problems/models"
)

/*

X = [0.05551216167487592 0.6334452066673657 0.7298491448394264 0.24887026620666236 0.6989197977528141 0.2354619573369429 0.31296550317593386 0.33573415823218383 0.6774521155725167 0.44943802053528487 0.24228785695240124 0.05490745702743821]

WFG RESULT

[2.68753587 1.01263411 1.09117641]
[0.6536774  0.69107753 6.64691246]
[0.69505637 0.94269431 6.09517195]
[2.00023624 2.74047567 4.82544364]
[0.48167976 1.80617806 5.45197913]
[0.38439    1.21786666 6.13984593]
[2.71063813 1.9549086  3.15477817]
[0.36701558 0.91058783 6.26816222]
[0.97918769 2.74415375 5.0323668 ]

GODE RESULT

[2.6420374115796585 3.999591700862304 0.9983151305987765]
[1.115249504798948 0.8694292081405048 4.8546541121664415]
[1.4353225628452153 1.6445908151331106 6.653524093431408]
[2.4447317161386453 2.177282376483311 3.172733477769529]
[0.31403163744480583 1.199800092973553 5.902948588049306]
[2.205462025597935 2.375291183587542 3.2175545089200175]
[2.7925062995980117 1.6541584728239944 2.1234581662991228]
[0.7449442034656699 1.9604097937811058 5.638385279536934]
[1.276092043233046 2.4167479402675838 4.997047902946538]

*/

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
	e := models.Elem{
		X: x,
	}
	wfg.WFG1.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg2(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Elem{
		X: x,
	}
	wfg.WFG2.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg3(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Elem{
		X: x,
	}
	wfg.WFG3.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg4(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Elem{
		X: x,
	}
	wfg.WFG4.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg5(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Elem{
		X: x,
	}
	wfg.WFG5.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg6(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Elem{
		X: x,
	}
	wfg.WFG6.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg7(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Elem{
		X: x,
	}
	wfg.WFG7.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg8(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Elem{
		X: x,
	}
	wfg.WFG8.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
func wfg9(x []float64, numObjs int) {
	_x := make([]float64, len(x))
	copy(_x, x)
	e := models.Elem{
		X: x,
	}
	wfg.WFG9.Fn(&e, numObjs)

	fmt.Println(e.Objs)
}
