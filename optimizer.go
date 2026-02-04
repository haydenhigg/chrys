package chrys

import (
	"math"
	// "fmt"
	"github.com/haydenhigg/chrys/algo"
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

	x       map[string]float64
	xDomain map[string]*Domain
}

// initializer
func NewOptimizer(objective ObjectiveFunc) *Optimizer {
	return &Optimizer{
		F:       objective,
		x:       map[string]float64{},
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
		return &Domain{math.Inf(-1), math.Inf(1), 1e-7}
	}
}

// sensitivity analysis
func (opt *Optimizer) perturb(x Parameters, k string, h float64) Parameters {
	domain := opt.Domain(k)

	x[k] += math.Copysign(max(math.Abs(h), domain.Resolution), h)

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
		h := max(opt.Domain(k).Resolution, 1e-7)
		forward := opt.F(opt.perturb(opt.X(), k, h))
		backward := opt.F(opt.perturb(opt.X(), k, -h))

		// (f(x + h) - f(x - h)) / |2h|
		derivatives[k] = (forward - backward) / math.Abs(2*h)
	}

	return derivatives
}

// OAT partial second derivatives using the finite difference method
func (opt *Optimizer) SecondDerivative() Parameters {
	baseline := opt.F(opt.X())
	derivatives := make(Parameters, len(opt.x))

	for k := range opt.x {
		h := max(opt.Domain(k).Resolution, 1e-4)
		forward := opt.F(opt.perturb(opt.X(), k, h))
		backward := opt.F(opt.perturb(opt.X(), k, -h))

		// (f(x + h) - 2f(x) + f(x - h)) / h^2
		derivatives[k] = (forward - 2*baseline + backward) / math.Pow(h, 2)
	}

	return derivatives
}

// OAT partial second derivatives using the finite difference method, averaged over multiple scales
func (opt *Optimizer) SmoothedSecondDerivative() Parameters {
	// baseline :=
	derivatives := make(map[string][]float64, len(opt.x))

	for k := range opt.x {
		n := 10
		derivatives[k] = make([]float64, n)

		for i := range n {
			h := math.Pow(2, float64(i))
			forward := opt.F(opt.perturb(opt.X(), k, h)) - opt.F(opt.X())
			backward := opt.F(opt.perturb(opt.X(), k, -h)) - opt.F(opt.X())

			// (f(x + h) - 2f(x) + f(x - h)) / h^2
			derivatives[k][i] = (forward + backward) / math.Pow(h, 2)
		}
	}

	smoothedDerivatives := make(Parameters, len(opt.x))

	for k, vs := range derivatives {
		smoothedDerivatives[k] = algo.Mean(vs)
	}

	return smoothedDerivatives
}

// optimization
func (opt *Optimizer) GradientDescent(learningRate float64, maxEpochs int) Parameters {
	for _ = range maxEpochs {
		shouldStop := true
		for k, partialGradient := range opt.Derivative() {
			if partialGradient == 0 {
				continue
			}

			// descend down the gradient
			opt.perturb(opt.x, k, -partialGradient*learningRate)
			shouldStop = false
		}

		// stop early if gradient == 0
		if shouldStop {
			break
		}
	}

	return opt.X()
}
