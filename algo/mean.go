package algo

func Mean(prices []float64) float64 {
	sum := 0.

	for _, x := range prices {
		sum += x
	}

	return sum / float64(len(prices))
}
