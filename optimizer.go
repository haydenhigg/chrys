package chrys

import (
	"errors"
	"math"
)

// a single-objective optimizing object
type Optimizer struct {
	inputs []float64
}

func NewOptimizer(inputs []float64) *Optimizer {
	return &Optimizer{inputs: inputs}
}

// Copy
func (opt *Optimizer) X() []float64 {
	x := make([]float64, len(opt.inputs))
	copy(x, opt.inputs)
	return x
}

// one-at-a-time partial derivatives using the finite difference method
func (opt *Optimizer) Derivative(
	f func([]float64) float64,
	hs []float64,
) ([]float64, error) {
	if len(hs) != len(opt.inputs) {
		return []float64{}, errors.New("x and h have different magnitudes")
	}

	derivatives := make([]float64, len(hs))
	for i, h := range hs {
		x := opt.X()
		dx := 2 * h

		x[i] += h
		forward := f(x)

		x[i] -= dx
		backward := f(x)

		// (f(x + h) - f(x - h)) / 2h
		derivatives[i] = (forward - backward) / math.Abs(dx)
	}

	return derivatives, nil
}

func (opt *Optimizer) SecondDerivative(
	f func([]float64) float64,
	hs []float64,
) ([]float64, error) {
	if len(hs) != len(opt.inputs) {
		return []float64{}, errors.New("x and h have different magnitudes")
	}

	derivatives := make([]float64, len(hs))
	for i, h := range hs {
		x := opt.X()

		stationary := f(x)

		x[i] += h
		forward := f(x)

		x[i] -= 2 * h
		backward := f(x)

		// (f(x + h) - 2f(x) + f(x - h)) / h^2
		derivatives[i] = (forward - 2*stationary + backward) / math.Pow(h, 2)
	}

	return derivatives, nil
}
