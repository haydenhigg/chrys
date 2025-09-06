package algo

func ZScore(prices []float64) float64 {
	mean := Mean(prices)
	standardDeviation := StandardDeviation(prices, mean)

	return (prices[len(prices)-1] - mean) / standardDeviation
}
