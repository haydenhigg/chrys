package optimizer

import (
	"math"
	"testing"
)

// tests -> GradientDescent
func Test_GradientDescent(t *testing.T) {
	// create Optimizer
	opt := New(func(x Input) float64 {
		return math.Pow(x["a"], 2) - 3*x["a"] + math.Sin(x["b"])
	})

	opt.SetX(Input{"a": 4, "b": 3})

	// GradientDescent()
	optimized := opt.GradientDescent(.1, 1000)

	// assert
	assertParametersEqual(optimized, Input{"a": 1.5, "b": 4.7123890}, t)
}
