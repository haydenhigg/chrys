package optimizer

import "math"

func (opt *Optimizer) xPlus(k string, h float64) Input {
	x := opt.X()
	x[k] = opt.withConstraints(k, x[k]+h)

	return x
}

// OAT partial derivatives using the finite difference method
func (opt *Optimizer) Derivative(stepSize float64) Input {
	h := max(stepSize, 1e-7)

	derivatives := make(Input, len(opt.x))
	for k := range opt.x {
		forward := opt.F(opt.xPlus(k, h))
		backward := opt.F(opt.xPlus(k, -h))

		// (f(x + h) - f(x - h)) / |2h|
		derivatives[k] = (forward - backward) / math.Abs(2*h)
	}

	return derivatives
}

func (opt *Optimizer) SmoothedDerivative(stepSize float64, n int) Input {
	stepSize = max(stepSize, 1e-7)

	sums := make(Input, len(opt.x))
	for i := range n {
		for k, v := range opt.Derivative(stepSize * math.Pow(2, float64(i))) {
			sums[k] += v
		}
	}

	means := make(Input, len(opt.x))
	for k, v := range sums {
		means[k] = v / float64(n)
	}

	return means
}

// OAT partial second derivatives using the finite difference method
func (opt *Optimizer) SecondDerivative(stepSize float64) Input {
	h := max(stepSize, 1e-4)
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

func (opt *Optimizer) SmoothedSecondDerivative(stepSize float64, n int) Input {
	stepSize = max(stepSize, 1e-4)
	baseline := opt.F(opt.X())

	sums := make(Input, len(opt.x))
	for i := range n {
		h := stepSize * math.Pow(2, float64(i))

		for k := range opt.x {
			forward := opt.F(opt.xPlus(k, h)) - baseline
			backward := opt.F(opt.xPlus(k, -h)) - baseline

			// (f(x + h) - 2f(x) + f(x - h)) / h^2
			sums[k] += (forward + backward) / math.Pow(h, 2)
		}
	}

	means := make(Input, len(opt.x))
	for k, v := range sums {
		means[k] = v / float64(n)
	}

	return means
}
