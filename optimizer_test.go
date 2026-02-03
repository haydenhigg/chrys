package chrys

import (
	"math"
	"testing"
)

// tests -> X
func Test_X(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer([]float64{4.1337, 3.37})

	// X()
	x := opt.X()

	// assert
	assertSlicesEqual(x, []float64{4.1337, 3.37}, t)
}

func Test_XPerturbed(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer([]float64{4.1337, 3.37})

	// X()
	x := opt.X(-.1337, .63)

	// assert
	assertSlicesEqual(x, []float64{4, 4}, t)
}

func Test_XPerturbedTooFew(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer([]float64{4.1337, 3.37})

	// X()
	x := opt.X(-.1337)

	// assert
	assertSlicesEqual(x, []float64{4, 3.37}, t)
}

func Test_XPerturbedTooMany(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer([]float64{4.1337, 3.37})

	// X()
	x := opt.X(-.1337, .63, 2.5)

	// assert
	assertSlicesEqual(x, []float64{4, 4}, t)
}

// tests -> Derivative
func Test_Derivative(t *testing.T) {
	// mock objective
	f := func(x []float64) float64 {
		return math.Pow(x[0], 2) + x[1]
	}

	// create Optimizer
	opt := NewOptimizer([]float64{4, 3})

	// Derivative()
	fPrime := opt.Derivative(f)

	// assert
	// f(a, b) = a^2 + b
	// f'(a) = 2a = 4, f'(b) = 1
	assertSlicesEqual(fPrime, []float64{8, 1}, t)
}

// tests ->
func Test_LocalSensitivity(t *testing.T) {
	// mock objective
	f := func(x []float64) float64 {
		return math.Pow(x[0], 2) + x[1]
	}

	// create Optimizer
	opt := NewOptimizer([]float64{4, 3})

	// Derivative()
	sens := opt.LocalSensitivity(f)

	// assert
	assertSlicesEqual(sens, []float64{.32, .03}, t)
}
