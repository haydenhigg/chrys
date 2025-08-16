package algo

import "math"

func ZScore(prices []float64) float64 {
	m := mean(prices)
	sumOfSquaredDifferences := 0.

	for _, price := range prices {
		sumOfSquaredDifferences += math.Pow(price-m, 2)
	}

	n := len(prices)
	standardDeviation := math.Sqrt(sumOfSquaredDifferences / float64(n))

	return (prices[n-1] - m) / standardDeviation
}
