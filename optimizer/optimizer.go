package optimizer

import (
	"math"
	"github.com/haydenhigg/chrys/algo"
	"maps"
)

type Input = map[string]float64
type Constraints = map[string][]Constraint
type ObjectiveFunc = func(Input) float64

type Optimizer struct {
	F ObjectiveFunc
	x Input
	constraints Constraints
}

// initializer
func New(objective ObjectiveFunc) *Optimizer {
	return &Optimizer{
		F: objective,
		x: Input{},
		constraints: Constraints{},
	}
}

// constraints
func (opt *Optimizer) Constrain(
	k string,
	newConstraints ...Constraint,
) *Optimizer {
	if constraints, ok := opt.constraints[k]; !ok {
		opt.constraints[k] = newConstraints
	} else {
		opt.constraints[k] = append(constraints, newConstraints...)
	}

	return opt
}

func (opt *Optimizer) withConstraints(k string, v float64) float64 {
	if constraints, ok := opt.constraints[k]; !ok {
		return applyConstraints(v, constraints)
	}

	return v
}

// X setter
func (opt *Optimizer) SetX(x Input) *Optimizer {
	for k, v := range x {
		opt.x[k] = opt.withConstraints(k, v)
	}

	return opt
}

// X getter
func (opt *Optimizer) X() Input {
	x := make(Input, len(opt.x))
	maps.Copy(x, opt.x)

	return x
}

// sensitivity analysis
func (opt *Optimizer) xPlus(k string, h float64) Input {
	x := opt.X()
	x[k] = opt.withConstraints(k, x[k] + h)

	return x
}

// OAT partial derivatives using the finite difference method
func (opt *Optimizer) Derivative(stepSize float64) Input {
	h := max(1e-7, stepSize)

	derivatives := make(Input, len(opt.x))

	for k := range opt.x {
		forward := opt.F(opt.xPlus(k, h))
		backward := opt.F(opt.xPlus(k, -h))

		// (f(x + h) - f(x - h)) / |2h|
		derivatives[k] = (forward - backward) / math.Abs(2*h)
	}

	return derivatives
}

// OAT partial second derivatives using the finite difference method
func (opt *Optimizer) SecondDerivative(stepSize float64) Input {
	h := max(1e-4, stepSize)
	baseline := opt.F(opt.X())

	derivatives := make(Input, len(opt.x))

	for k := range opt.x {
		forward := opt.F(opt.xPlus(k, h))
		backward := opt.F(opt.xPlus(k, -h))

		// (f(x + h) - 2f(x) + f(x - h)) / h^2
		derivatives[k] = (forward - 2*baseline + backward) / math.Pow(h, 2)
	}

	return derivatives
}

// OAT partial second derivatives using the finite difference method, averaged over multiple scales
func (opt *Optimizer) SmoothedSecondDerivative(stepSize float64, n int) Input {
	baseline := opt.F(opt.X())
	derivatives := make(map[string][]float64, len(opt.x))

	for k := range opt.x {
		derivatives[k] = make([]float64, n)

		h := max(1e-4, stepSize)


		for i := range n {
			forward := opt.F(opt.xPlus(k, h)) - baseline
			backward := opt.F(opt.xPlus(k, -h)) - baseline

			// (f(x + h) - 2f(x) + f(x - h)) / h^2
			derivatives[k][i] = (forward + backward) / math.Pow(h, 2)

			h *= 2
		}
	}

	smoothedDerivatives := make(Input, len(opt.x))

	for k, vs := range derivatives {
		smoothedDerivatives[k] = algo.Mean(vs)
	}

	return smoothedDerivatives
}

// optimization
func (opt *Optimizer) GradientDescent(learningRate float64, maxEpochs int) Input {
	for _ = range maxEpochs {
		shouldStop := true
		for k, partialGradient := range opt.Derivative(0) {
			if partialGradient == 0 {
				continue
			}

			shouldStop = false

			// descend down the gradient
			delta := partialGradient * learningRate
			opt.x[k] = opt.withConstraints(k, opt.x[k] - delta)
		}

		// stop early if gradient == 0
		if shouldStop {
			break
		}
	}

	return opt.X()
}
