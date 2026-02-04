package chrys

import (
	"math"
	"fmt"
	"maps"
)

type Parameters = map[string]float64
type ObjectiveFunc = func(Parameters) float64

type Domain struct {
	Lower,
	Upper,
	Resolution float64
}

type Optimizer struct {
	F ObjectiveFunc

	x map[string]float64
	xDomain map[string]*Domain
}

// initializer
func NewOptimizer(objective ObjectiveFunc) *Optimizer {
	return &Optimizer{
		F: objective,
		x: map[string]float64{},
		xDomain: map[string]*Domain{},
	}
}

// X setters
func (opt *Optimizer) CreateX(
	k string, v float64,
	lowerBound float64,
	upperBound float64,
	resolution float64,
) *Optimizer {
	opt.SetX(k, v)
	opt.xDomain[k] = &Domain{lowerBound, upperBound, math.Abs(resolution)}

	return opt
}

func (opt *Optimizer) SetX(k string, v float64) *Optimizer {
	opt.x[k] = v
	return opt
}

// X getter
func (opt *Optimizer) X() Parameters {
	x := make(Parameters, len(opt.x))
	maps.Copy(x, opt.x)

	return x
}

func (opt *Optimizer) Domain(k string) *Domain {
	if domain, ok := opt.xDomain[k]; ok {
		return domain
	} else {
		return &Domain{math.Inf(-1), math.Inf(1), 1e-8}
	}
}

func (opt *Optimizer) XPerturb(k string, h float64) Parameters {
	domain := opt.Domain(k)

	x := opt.X()
	x[k] += math.Copysign(max(math.Abs(h), domain.Resolution), h)

	// clamp to bounds
	if x[k] > domain.Upper {
		x[k] = domain.Upper
	} else if x[k] < domain.Lower {
		x[k] = domain.Lower
	}

	fmt.Println(x)

	return x
}

// OAT partial derivatives using the finite difference method
func (opt *Optimizer) Derivative() Parameters {
	derivatives := make(Parameters, len(opt.x))
	for k := range opt.x {
		h := opt.Domain(k).Resolution
		forward := opt.F(opt.XPerturb(k, h))
		backward := opt.F(opt.XPerturb(k, -h))

		// (f(x + h) - f(x - h)) / 2h
		derivatives[k] = (forward - backward) / math.Abs(2 * h)
	}

	return derivatives
}

// OAT local perturbations, similar to a derivative
const epsilon float64 = .1

func (opt *Optimizer) LocalSensitivity() Parameters {
	baseline := opt.F(opt.X())
	sensitivities := make(Parameters, len(opt.x))
	for k, v := range opt.x {
		domain := opt.Domain(k)

		delta := epsilon
		if math.IsInf(domain.Upper, 0) && math.IsInf(domain.Lower, 0) {
			delta *= v
		} else if math.IsInf(domain.Upper, 0) {
			delta *= v - domain.Lower
		} else if math.IsInf(domain.Lower, 0) {
			delta *= domain.Upper - v
		} else {
			delta *= domain.Upper - domain.Lower
		}

		fmt.Println(opt.F(opt.XPerturb(k, delta)) / baseline - 1)
		fmt.Println(opt.F(opt.XPerturb(k, -delta)) / baseline - 1)

		plus := math.Abs(opt.F(opt.XPerturb(k, delta)) / baseline - 1)
		minus := math.Abs(opt.F(opt.XPerturb(k, -delta)) / baseline -1)

		// percentage-wise, how much of a change in f(x) results from wiggling?
		sensitivities[k] = (plus + minus) / (2 * epsilon)
	}

	return sensitivities
}
