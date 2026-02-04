package optimizer

import (
	"math"
	"math/rand"
	"testing"
)

// helpers
func assertParametersEqual(a, b Input, t *testing.T) {
	for k, va := range a {
		if vb, ok := b[k]; !ok {
			t.Errorf(`b["%s"] does not exist`, k)
		} else if math.Abs(va-vb) > 1e-6 {
			t.Errorf(`a["%s"] != b["%s"]: %v != %v`, k, k, va, vb)
		}
	}

	for k := range b {
		if _, ok := a[k]; !ok {
			t.Errorf(`a["%s"] does not exist`, k)
		}
	}
}

// tests
// tests -> constraints
func Test_Constrain(t *testing.T) {
	// create Optimizer
	opt := New(func(x Input) float64 { return 0 })

	// Constrain()
	opt.Constrain("abc", Min(0), Max(1))

	// assert
	if constraints, ok := opt.constraints["abc"]; !ok || len(constraints) < 2 {
		t.Errorf("constraints not added")
	}
}

func Test_withConstraints(t *testing.T) {
	// create Optimizer
	opt := New(func(x Input) float64 { return 0 })

	// Constrain()
	opt.Constrain("abc", Min(0), Max(1))

	// withConstraints()
	v := opt.withConstraints("abc", 1.5)

	// assert
	if v != 1.5 {
		t.Errorf("v != 1.5: %f", v)
	}
}

// tests -> X setter
func Test_SetX(t *testing.T) {
	// create Optimizer
	opt := New(func(x Input) float64 { return 0 })

	// SetX()
	x := Input{"abc": 13.37}
	opt.SetX(x)

	// assert
	assertParametersEqual(opt.x, x, t)
}

// tests -> X getter
func Test_X(t *testing.T) {
	// create Optimizer
	opt := New(func(x Input) float64 { return 0 })

	opt.SetX(Input{"a": 4.1337})
	opt.SetX(Input{"b": 3.37})

	// X()
	x := opt.X()

	// assert
	assertParametersEqual(x, Input{"a": 4.1337, "b": 3.37}, t)
}

func Test_XDeepCopy(t *testing.T) {
	// create Optimizer
	opt := New(func(x Input) float64 { return 0 })

	opt.SetX(Input{"a": 4.1337})

	// X()
	x := opt.X()
	x["a"] = 5

	// assert
	assertParametersEqual(x, Input{"a": 5}, t)
	assertParametersEqual(opt.x, Input{"a": 4.1337}, t)
}

// tests -> xPlus
func Test_xPlus(t *testing.T) {
	// create Optimizer
	opt := New(func(x Input) float64 { return 0 })

	opt.SetX(Input{"a": 4.1337, "b": 3.337})

	// perturb()
	x := opt.xPlus("b", 2e-3)

	// assert
	assertParametersEqual(x, Input{"a": 4.1337, "b": 3.339}, t)
}

// tests -> derivatives
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
		noise := (rand.Float64() - 0.5)*1e-7
		return math.Pow(x["a"], 3) - 3*x["a"] + math.Pow(x["b"], 2) + noise
	})

	opt.SetX(Input{"a": 4, "b": 3})

	// SecondDerivative()
	sensitivity := opt.SmoothedSecondDerivative(1, 20)

	// assert
	// f(a, b) = a^3 - 3a + x^2
	// f'(a) = 3a^2 - 3, f'(b) = 2x
	// f''(a) = 6a, f''(b) = 2
	assertParametersEqual(sensitivity, Input{"a": 24, "b": 2}, t)
}

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
