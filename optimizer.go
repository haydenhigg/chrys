package chrys

import "math"

// a single-objective optimizing object
type Optimizer struct {
	inputs []float64
}

func NewOptimizer(inputs []float64) *Optimizer {
	return &Optimizer{inputs: inputs}
}

// Copy and Perturb
func (opt *Optimizer) X(deltas ...float64) []float64 {
	x := make([]float64, len(opt.inputs))
	copy(x, opt.inputs)

	for i := range min(len(x), len(deltas)) {
		x[i] += deltas[i]
	}

	return x
}

// one-at-a-time partial derivatives using the finite difference method
func (opt *Optimizer) Derivative(f func([]float64) float64) []float64 {
	h := 1e-5

	derivatives := make([]float64, len(opt.inputs))
	for i := range opt.inputs {
		x := opt.X()
		dx := 2 * h

		x[i] += h
		forward := f(x)

		x[i] -= dx
		backward := f(x)

		// (f(x + h) - f(x - h)) / 2h
		derivatives[i] = (forward - backward) / math.Abs(dx)
	}

	return derivatives
}

// local perturbations
func (opt *Optimizer) LocalSensitivity(f func([]float64) float64) []float64 {
	x := opt.X()
	baseline := f(x)

	sensitivities := make([]float64, len(opt.inputs))
	for i := range opt.inputs {
		deltas := make([]float64, len(opt.inputs))

		deltas[i] += 0.01 * x[i]
		plus := f(opt.X(deltas...)) - baseline

		deltas[i] -= 0.02 * x[i]
		minus := f(opt.X(deltas...)) - baseline

		sensitivities[i] = (math.Abs(plus) + math.Abs(minus)) / 2
	}

	return sensitivities
}
