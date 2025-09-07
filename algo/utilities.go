package algo

import (
	"github.com/haydenhigg/chrys"
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

// frame utilities
func MapFrames(
	frames []*chrys.Frame,
	processor func(*chrys.Frame) float64,
) []float64 {
	processed := make([]float64, len(frames))
	for i, frame := range frames {
		processed[i] = processor(frame)
	}

	return processed
}

func Opens(frames []*chrys.Frame) []float64 {
	return MapFrames(frames, func(f *chrys.Frame) float64 { return f.Open })
}

func Highs(frames []*chrys.Frame) []float64 {
	return MapFrames(frames, func(f *chrys.Frame) float64 { return f.High })
}

func Lows(frames []*chrys.Frame) []float64 {
	return MapFrames(frames, func(f *chrys.Frame) float64 { return f.Low })
}

func Closes(frames []*chrys.Frame) []float64 {
	return MapFrames(frames, func(f *chrys.Frame) float64 { return f.Close })
}
