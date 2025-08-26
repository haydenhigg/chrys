package algo

import "chrys"

func Closes(frames []*chrys.Frame) []float64 {
	closes := make([]float64, len(frames))
	for i, frame := range frames {
		closes[i] = frame.Close
	}

	return closes
}

func mean(xs []float64) float64 {
	sum := 0.

	for _, x := range xs {
		sum += x
	}

	return sum / float64(len(xs))
}
