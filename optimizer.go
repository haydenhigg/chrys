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

// OAT partial derivatives using the finite difference method
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

// OAT local perturbations, very much like a second derivative but not quite
func (opt *Optimizer) LocalSensitivity(f func([]float64) float64) []float64 {
	// epsilons := make([]float64, len(opt.inputs))
	// for i, x := range opt.inputs {
	// 	epsilons[i] = x * 1e-3
	// }

	h := 1e-3

	baseline := f(opt.X())
	sensitivities := make([]float64, len(opt.inputs))
	for i := range opt.inputs {
		x := opt.X()
		dx := 2 * h

		x[i] += h
		plus := f(x) - baseline

		x[i] -= dx
		minus := f(x) - baseline

		// (|f(x + h) - f(x)| + |f(x - h) - f(x)|) / 2h
		sensitivities[i] = (math.Abs(plus) + math.Abs(minus)) / math.Abs(dx)
	}

	return sensitivities
}
