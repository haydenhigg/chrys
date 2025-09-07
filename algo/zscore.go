package algo

func ZScore(xs []float64) float64 {
	mean := Mean(xs)
	stdev := StandardDeviation(xs, mean)

	return (xs[len(xs) - 1] - mean) / stdev
}
