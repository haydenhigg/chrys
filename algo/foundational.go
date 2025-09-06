package algo

import "math"

func Mean(xs []float64) float64 {
	sum := 0.

	for _, x := range xs {
		sum += x
	}

	return sum / float64(len(xs))
}

func Variance(xs []float64, mean float64) float64 {
	sum := 0.

	for _, x := range xs {
		sum += math.Pow(x-mean, 2)
	}

	return sum / float64(len(xs) - 1)
}

func StandardDeviation(xs []float64, mean float64) float64 {
	return math.Sqrt(Variance(xs, mean))
}
