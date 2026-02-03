package chrys

import (
	"math"
	"testing"
)

func Test_Derivative(t *testing.T) {
	// mock objective
	f := func(x []float64) float64 {
		return math.Pow(x[0], 2) + x[1]
	}

	// create Optimizer
	opt := NewOptimizer([]float64{4, 3})

	// Derivative()
	fPrime, err := opt.Derivative(f, []float64{1e-5, 1e-5})
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	// f(a, b) = a^2 + b
	// f'(a) = 2a = 4, f'(b) = 1
	assertSlicesEqual(fPrime, []float64{8, 1}, t)
}

func Test_SecondDerivative(t *testing.T) {
	// mock objective
	f := func(x []float64) float64 {
		return math.Pow(x[0], 2) + x[1]
	}

	// create Optimizer
	opt := NewOptimizer([]float64{4, 3})

	// SecondDerivative()
	fPrimePrime, err := opt.SecondDerivative(f, []float64{1e-3, 1e-3})
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	// f(a, b) = a^2 + b
	// f'(a) = 2a, f'(b) = 1
	// f''(a) = 2, f''(b) = 0
	assertSlicesEqual(fPrimePrime, []float64{2, 0}, t)
}
