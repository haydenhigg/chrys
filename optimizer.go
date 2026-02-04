package chrys

import (
	"math"
	// "math/rand"
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
	x[k] += max(h, domain.Resolution)

	// clamp to bounds
	if x[k] > domain.Upper {
		x[k] = domain.Upper
	} else if x[k] < domain.Lower {
		x[k] = domain.Lower
	}

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

// OAT local perturbations, very much like a derivative but not quite
func (opt *Optimizer) LocalSensitivity() Parameters {
	// isLowerBoundInf := math.IsInf(lowerBound, 0)
	// isUpperBoundInf := math.IsInf(upperBound, 0)

	// h := resolution
	// if isLowerBoundInf && isUpperBoundInf {

	// } else if isLowerBoundInf {
	// } else if isUpperBoundInf {
	// } else {
	// }

	// math.IsInf(lowerBound, -1)
	// math.IsInf(domain[0], 1)
	// h := rand.Float64()

	// // // //

	// baseline := f(opt.X())
	sensitivities := make(Parameters, len(opt.x))
	// for i := range opt.inputs {
	// 	x := opt.X()
	// 	dx := 2 * h

	// 	x[i] += h
	// 	plus := f(x) - baseline

	// 	x[i] -= dx
	// 	minus := f(x) - baseline

	// 	// (|f(x + h) - f(x)| + |f(x - h) - f(x)|) / 2h
	// 	sensitivities[i] = (math.Abs(plus) + math.Abs(minus)) / math.Abs(dx)
	// }

	return sensitivities
}
