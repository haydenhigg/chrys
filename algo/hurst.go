package algo

import "math"

// estimate the Hurst exponent from returns using Variance of Differences
func Hurst(xs []float64) float64 {
	n := len(xs)
	if n < 4 {
		return 0.5
	}

	// too large increases noise; too small is insufficient for regression
	maxLag := min(n/2, 50)

	logLags := make([]float64, 0, maxLag)
	logVariances := make([]float64, 0, maxLag)

	for lag := 1; lag <= maxLag; lag++ {
		m := n - lag
		diffs := make([]float64, m)
		for i := range m {
			diffs[i] = xs[i+lag] - xs[i]
		}

		variance := Variance(diffs, Mean(diffs))
		if variance <= 0 {
			continue
		}

		logLags = append(logLags, math.Log(float64(lag)))
		logVariances = append(logVariances, math.Log(variance))
	}

	// OLS slope = cov(x,y)/var(x)
	meanLogLag := Mean(logLags)
	meanLogVariance := Mean(logVariances)

	var num, den float64
	for i := range logLags {
		dx := logLags[i] - meanLogLag

		num += dx * (logVariances[i] - meanLogVariance)
		den += dx * dx
	}

	if den == 0 {
		return 0.
	}

	slope := num / den

	return min(max(slope/2, 0), 1)
}
