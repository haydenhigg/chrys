package optimizer

import (
	"math"
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
	if v != 1 {
		t.Errorf("v != 1: %f", v)
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
