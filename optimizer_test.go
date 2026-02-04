package chrys

import (
	"math"
	"testing"
	"fmt"
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
	if math.Abs(a.Lower - b.Lower) > 1e-6 {
		t.Errorf(`a.Lower != b.Lower: %f != %f`, a.Lower, b.Lower)
	}
	if math.Abs(a.Upper - b.Upper) > 1e-6 {
		t.Errorf(`a.Upper != b.Upper: %f != %f`, a.Upper, b.Upper)
	}
	if math.Abs(a.Resolution - b.Resolution) > 1e-6 {
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

// tests -> XPerturb
func Test_XPerturb(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer(func(x Parameters) float64 { return 0 })

	opt.SetX("a", 4.1337)
	opt.CreateX("b", 3.337, 0, math.Inf(1), 1e-3)

	// XPerturb()
	x := opt.XPerturb("b", 2e-3)

	// assert
	assertParametersEqual(x, Parameters{"a": 4.1337, "b": 3.339}, t)
}

func Test_XPerturbOutsideBounds(t *testing.T) {
	// create Optimizer
	opt := NewOptimizer(func(x Parameters) float64 { return 0 })

	opt.SetX("a", 4.1337)
	opt.CreateX("b", 3.337, 0, 3.338, 1e-3)

	// XPerturb()
	x := opt.XPerturb("b", 2e-3)

	// assert
	assertParametersEqual(x, Parameters{"a": 4.1337, "b": 3.338}, t)
}

// tests -> Derivative
func Test_Derivative(t *testing.T) {
	// mock objective
	f := func(x Parameters) float64 {
		return math.Pow(x["a"], 2) + x["b"]
	}

	// create Optimizer
	opt := NewOptimizer(f)

	opt.SetX("a", 4)
	opt.SetX("b", 3)

	// Derivative()
	derivative := opt.Derivative()

	fmt.Println(derivative)

	// assert
	// f(a, b) = a^2 + b
	// f'(a) = 2a = 4, f'(b) = 1
	assertParametersEqual(derivative, Parameters{"a": 8, "b": 1}, t)
}

// tests -> LocalSensitivity
// func Test_LocalSensitivity(t *testing.T) {
// 	// mock objective
// 	f := func(x []float64) float64 {
// 		return math.Pow(x[0], 2) + x[1]
// 	}

// 	// create Optimizer
// 	opt := NewOptimizer([]float64{4, 3})

// 	// Derivative()
// 	sens := opt.LocalSensitivity(f)

// 	// assert
// 	assertSlicesEqual(sens, []float64{.32, .03}, t)
// }
