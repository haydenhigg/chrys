package algo

func mean(xs []float64) float64 {
	sum := 0.

	for _, x := range xs {
		sum += x
	}

	return sum / float64(len(xs))
}
