package chrys

import (
	"math"
	"math/rand"
	"testing"
)

// helpers
func assertParametersEqual(a, b Parameters, t *testing.T) {
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

func assertDomainsEqual(a, b *Domain, t *testing.T) {
	if math.Abs(a.Lower-b.Lower) > 1e-6 {
		t.Errorf(`a.Lower != b.Lower: %f != %f`, a.Lower, b.Lower)
	}
	if math.Abs(a.Upper-b.Upper) > 1e-6 {
		t.Errorf(`a.Upper != b.Upper: %f != %f`, a.Upper, b.Upper)
	}
	if math.Abs(a.Resolution-b.Resolution) > 1e-6 {
		t.Errorf(
			`a.Resolution != b.Resolution: %f != %f`,
			a.Resolution, b.Resolution,
		)
	}
}

func assertDomainMapsEqual(a, b map[string]*Domain, t *testing.T) {
	for k, va := range a {
		if vb, ok := b[k]; !ok {
			t.Errorf(`b["%s"] does not exist`, k)
		} else {
			assertDomainsEqual(va, vb, t)
		}
	}

	for k := range b {
		if _, ok := a[k]; !ok {
			t.Errorf(`a["%s"] does not exist`, k)
		}
	}
}

// tests
// tests -> CreateX
func Test_CreateX(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer(func(x Parameters) float64 { return 0 })

	// CreateX()
	opt.CreateX("abc", -13.37, math.Inf(-1), 0, 1e-2)

	// assert
	assertParametersEqual(opt.x, Parameters{"abc": -13.37}, t)
	assertDomainMapsEqual(opt.xDomain, map[string]*Domain{
		"abc": {math.Inf(-1), 0, 1e-2},
	}, t)
}

// tests -> SetX
func Test_SetX(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer(func(x Parameters) float64 { return 0 })

	// SetX()
	opt.SetX("abc", -13.37)

	// assert
	assertParametersEqual(opt.x, Parameters{"abc": -13.37}, t)
	assertDomainMapsEqual(opt.xDomain, map[string]*Domain{}, t)
}

// tests -> X
func Test_X(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer(func(x Parameters) float64 { return 0 })

	opt.SetX("a", 4.1337)
	opt.SetX("b", 3.37)

	// X()
	x := opt.X()

	// assert
	assertParametersEqual(x, Parameters{"a": 4.1337, "b": 3.37}, t)
}

// tests -> Domain
func Test_Domain(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer(func(x Parameters) float64 { return 0 })

	opt.CreateX("abc", -13.37, math.Inf(-1), 0, 1e-2)

	// Domain()
	domain := opt.Domain("abc")

	// assert
	assertDomainsEqual(domain, &Domain{math.Inf(-1), 0, 1e-2}, t)
}

func Test_DomainDefault(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer(func(x Parameters) float64 { return 0 })

	opt.SetX("abc", -13.37)

	// X()
	domain := opt.Domain("abc")

	// assert
	assertDomainsEqual(domain, &Domain{math.Inf(-1), math.Inf(1), 1e-8}, t)
}

// tests -> perturb
func Test_perturb(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer(func(x Parameters) float64 { return 0 })

	opt.SetX("a", 4.1337)
	opt.CreateX("b", 3.337, 0, math.Inf(1), 1e-3)

	// perturb()
	x := opt.perturb(opt.X(), "b", 2e-3)

	// assert
	assertParametersEqual(x, Parameters{"a": 4.1337, "b": 3.339}, t)
}

func Test_perturbOutsideBounds(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer(func(x Parameters) float64 { return 0 })

	opt.SetX("a", 4.1337)
	opt.CreateX("b", 3.337, 0, 3.338, 1e-3)

	// perturb()
	x := opt.perturb(opt.X(), "b", 2e-3)

	// assert
	assertParametersEqual(x, Parameters{"a": 4.1337, "b": 3.338}, t)
}

func Test_perturbUnderResolution(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer(func(x Parameters) float64 { return 0 })

	opt.SetX("a", 4.1337)
	opt.CreateX("b", 3.337, 0, math.Inf(1), 1e-3)

	// perturb()
	x := opt.perturb(opt.X(), "b", 1e-4)

	// assert
	assertParametersEqual(x, Parameters{"a": 4.1337, "b": 3.338}, t)
}

// tests -> Derivative
func Test_Derivative(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer(func(x Parameters) float64 {
		return math.Pow(x["a"], 2) + x["b"]
	})

	opt.SetX("a", 4)
	opt.SetX("b", 3)

	// Derivative()
	derivative := opt.Derivative()

	// assert
	// f(a, b) = a^2 + b
	// f'(a) = 2a = 4, f'(b) = 1
	assertParametersEqual(derivative, Parameters{"a": 8, "b": 1}, t)
}

// test -> SecondDerivative
func Test_SecondDerivative(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer(func(x Parameters) float64 {
		return math.Pow(x["a"], 2) - 3*x["a"] + math.Sin(x["b"])
	})

	opt.SetX("a", 4)
	opt.SetX("b", 3)

	// SecondDerivative()
	sensitivity := opt.SecondDerivative()

	// assert
	// f(a, b) = a^2 - 3a + sin(b)
	// f'(a) = 2a - 3, f'(b) = cos(b)
	// f''(a) = 2, f''(b) = -sin(b)
	assertParametersEqual(sensitivity, Parameters{"a": 2, "b": -0.14112000806}, t)
}

// test -> SmoothedSecondDerivative
func Test_SmoothedSecondDerivative(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer(func(x Parameters) float64 {
		return math.Pow(x["a"], 3) - 3*x["a"] + math.Pow(x["b"], 2) + rand.Float64()*1e-7
	})

	opt.SetX("a", 4)
	opt.SetX("b", 3)

	// SecondDerivative()
	sensitivity := opt.SmoothedSecondDerivative()

	// assert
	// f(a, b) = a^3 - 3a + x^2
	// f'(a) = 3a^2 - 3, f'(b) = 2x
	// f''(a) = 6a, f''(b) = 2
	assertParametersEqual(sensitivity, Parameters{"a": 24, "b": 2}, t)
}

// tests -> GradientDescent
func Test_GradientDescent(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer(func(x Parameters) float64 {
		return math.Pow(x["a"], 2) - 3*x["a"] + math.Sin(x["b"])
	})

	opt.SetX("a", 2)
	opt.SetX("b", 4)

	// GradientDescent()
	optimized := opt.GradientDescent(.1, 1000)

	// assert
	assertParametersEqual(optimized, Parameters{"a": 1.5, "b": 4.7123890}, t)
}
