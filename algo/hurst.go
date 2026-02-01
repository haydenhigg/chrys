package algo

import "math"

// estimates the Hurst exponent from prices using Variance-of-Differences
func Hurst(prices []float64) float64 {
	n := len(prices)
	if n < 32 {
		return 0.5
	}

	xs := []float64{}
	ys := []float64{}

	for lag := 2; lag <= min(n/2, 256); lag *= 2 {
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
