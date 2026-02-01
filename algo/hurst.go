package algo

import "math"

// estimate the Hurst exponent from returns using Variance of Differences
func Hurst(returns []float64) float64 {
	n := len(returns)
	if n < 4 {
		return 0.5
	}

	// too large increases noise; too small is insufficient for regression
	maxLag := min(n/2, 50)

	xs := make([]float64, 0, maxLag)
	ys := make([]float64, 0, maxLag)

	for lag := 1; lag <= maxLag; lag++ {
		m := n - lag
		diffs := make([]float64, m)
		for i := range m {
			diffs[i] = returns[i+lag] - returns[i]
		}

		variance := Variance(diffs, Mean(diffs))
		if variance <= 0 {
			continue
		}

		xs = append(xs, math.Log(float64(lag)))
		ys = append(ys, math.Log(variance))
	}

	// OLS slope = cov(x,y)/var(x)
	meanX := Mean(xs)
	meanY := Mean(ys)

	var num, den float64
	for i := range xs {
		dx := xs[i] - meanX

		num += dx * (ys[i] - meanY)
		den += dx * dx
	}

	if den == 0 {
		return 0.
	}

	slope := num / den

	return min(max(slope/2, 0), 1)
}
