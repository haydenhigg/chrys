package algo

import "math"

func ZScore(prices []float64) float64 {
	mu := mean(prices)
	sumOfSquaredDifferences := 0.

	for _, price := range prices {
		sumOfSquaredDifferences += math.Pow(price-mu, 2)
	}

	n := len(prices)
	sigma := math.Sqrt(sumOfSquaredDifferences / float64(n))

	return (prices[n-1] - mu) / sigma
}
