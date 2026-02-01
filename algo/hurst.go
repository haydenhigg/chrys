package algo

import "math"

// estimate the Hurst exponent from prices using Variance-of-Differences
func Hurst(prices []float64) float64 {
	n := len(prices)
	if n < 40 {
		return 0.5
	}

	xs := []float64{}
	ys := []float64{}

	for lag := 2; lag <= max(n / 10, 32); lag *= 2 {
		m := n - lag
		diffs := make([]float64, m)
		for i := range m {
			diffs[i] = prices[i+lag] - prices[i]
		}

		variance := Variance(diffs, Mean(diffs))
		if variance <= 0 {
			continue
		}

		xs = append(xs, math.Log(float64(lag)))
		ys = append(ys, math.Log(variance))
	}

	// OLS slope = cov(x,y)/var(x)
	meanLogLag := Mean(xs)
	meanLogVariance := Mean(ys)

	var num, den float64
	for i := range xs {
		dx := xs[i] - meanLogLag

		num += dx * (ys[i] - meanLogVariance)
		den += dx * dx
	}

	if den == 0 {
		return 0.
	}

	slope := num / den

	return min(max(slope/2, 0), 1)
}
