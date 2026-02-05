package optimizer

import "math/rand"

func (opt *Optimizer) RandomSearch(
	changeRate float64,
	maxEpochs int,
) Input {
	baseline := opt.F(opt.X())

	for _ = range maxEpochs {
		x := opt.X()
		for k := range x {
			if rand.Float64() < .5 {
				coef := changeRate
				if rand.Float64() < .5 {
					coef = -changeRate
				}

				x[k] = opt.withConstraints(k, x[k]*(1+coef))
			}
		}

		score := opt.F(x)
		if opt.F(x) < baseline {
			baseline = score
			opt.SetX(x)
		}
	}

	return opt.X()
}

func (opt *Optimizer) GradientDescent(
	learningRate float64,
	maxEpochs int,
) Input {
	for _ = range maxEpochs {
		shouldStop := true
		for k, partialGradient := range opt.Derivative(0) {
			if partialGradient == 0 {
				continue
			}

			shouldStop = false

			// descend down the gradient
			delta := partialGradient * learningRate
			opt.x[k] = opt.withConstraints(k, opt.x[k]-delta)
		}

		// stop early if gradient == 0
		if shouldStop {
			break
		}
	}

	return opt.X()
}
