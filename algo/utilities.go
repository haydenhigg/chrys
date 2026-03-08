package algo

import (
	domain "github.com/haydenhigg/chrys/frame"
	"math"
)

// mathematical utilities
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

	return sum / float64(len(xs)-1)
}

func StandardDeviation(xs []float64, mean float64) float64 {
	return math.Sqrt(Variance(xs, mean))
}

func Covariance(xs, ys []float64, xMean, yMean float64) float64 {
	sum := 0.

	for i, x := range xs {
		sum += (x - xMean) * (ys[i] - yMean)
	}

	return sum / float64(len(xs)-1)
}

func Correlation(xs, ys []float64, xMean, yMean float64) float64 {
	numerator := Covariance(xs, ys, xMean, yMean)
	denominator := StandardDeviation(xs, xMean) * StandardDeviation(ys, yMean)

	return numerator / denominator
}

// frame utilities
func MapFrames(
	frames []*domain.Frame,
	processor func(*domain.Frame) float64,
) []float64 {
	processed := make([]float64, len(frames))
	for i, frame := range frames {
		processed[i] = processor(frame)
	}

	return processed
}

func Opens(frames []*domain.Frame) []float64 {
	return MapFrames(frames, func(f *domain.Frame) float64 { return f.Open })
}

func Highs(frames []*domain.Frame) []float64 {
	return MapFrames(frames, func(f *domain.Frame) float64 { return f.High })
}

func Lows(frames []*domain.Frame) []float64 {
	return MapFrames(frames, func(f *domain.Frame) float64 { return f.Low })
}

func Closes(frames []*domain.Frame) []float64 {
	return MapFrames(frames, func(f *domain.Frame) float64 { return f.Close })
}

func Volumes(frames []*domain.Frame) []float64 {
	return MapFrames(frames, func(f *domain.Frame) float64 { return f.Volume })
}
