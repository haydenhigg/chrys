package optimizer

import (
	"math"
	"math/rand"
	"testing"
)

func Test_xPlus(t *testing.T) {
	// create Optimizer
	opt := New(func(x Input) float64 { return 0 })

	opt.SetX(Input{"a": 4.1337, "b": 3.337})

	// perturb()
	x := opt.xPlus("b", 2e-3)

	// assert
	assertParametersEqual(x, Input{"a": 4.1337, "b": 3.339}, t)
}

func Test_xPlusConstrained(t *testing.T) {
	// create Optimizer
	opt := New(func(x Input) float64 { return 0 })

	opt.SetX(Input{"a": 4.1337, "b": 3.337})
	opt.Constrain("b", Max(3.338))

	// perturb()
	x := opt.xPlus("b", 2e-3)

	// assert
	assertParametersEqual(x, Input{"a": 4.1337, "b": 3.338}, t)
}

func Test_Derivative(t *testing.T) {
	// create Optimizer
	opt := New(func(x Input) float64 {
		return math.Pow(x["a"], 2) + x["b"]
	})

	opt.SetX(Input{"a": 4, "b": 3})

	// Derivative()
	derivative := opt.Derivative(0)

	// assert
	// f(a, b) = a^2 + b
	// f'(a) = 2a = 4, f'(b) = 1
	assertParametersEqual(derivative, Input{"a": 8, "b": 1}, t)
}

func Test_SmoothedDerivative(t *testing.T) {
	// create Optimizer
	opt := New(func(x Input) float64 {
		noise := (rand.Float64() - 0.5) * 1e-7
		return math.Pow(x["a"], 2) + x["b"] + noise
	})

	opt.SetX(Input{"a": 4, "b": 3})

	// Derivative()
	derivative := opt.SmoothedDerivative(1e-2, 5)

	// assert
	// f(a, b) = a^2 + b
	// f'(a) = 2a = 4, f'(b) = 1
	assertParametersEqual(derivative, Input{"a": 8, "b": 1}, t)
}

func Test_SecondDerivative(t *testing.T) {
	// create Optimizer
	opt := New(func(x Input) float64 {
		return math.Pow(x["a"], 2) - 3*x["a"] + math.Sin(x["b"])
	})

	opt.SetX(Input{"a": 4, "b": 3})

	// SecondDerivative()
	sensitivity := opt.SecondDerivative(0)

	// assert
	// f(a, b) = a^2 - 3a + sin(b)
	// f'(a) = 2a - 3, f'(b) = cos(b)
	// f''(a) = 2, f''(b) = -sin(b)
	assertParametersEqual(sensitivity, Input{"a": 2, "b": -0.14112000806}, t)
}

func Test_SmoothedSecondDerivative(t *testing.T) {
	// create Optimizer
	opt := New(func(x Input) float64 {
		noise := (rand.Float64() - 0.5) * 1e-7
		return math.Pow(x["a"], 3) - 3*x["a"] + math.Pow(x["b"], 2) + noise
	})

	opt.SetX(Input{"a": 4, "b": 3})

	// SecondDerivative()
	sensitivity := opt.SmoothedSecondDerivative(1, 5)

	// assert
	// f(a, b) = a^3 - 3a + x^2
	// f'(a) = 3a^2 - 3, f'(b) = 2x
	// f''(a) = 6a, f''(b) = 2
	assertParametersEqual(sensitivity, Input{"a": 24, "b": 2}, t)
}
