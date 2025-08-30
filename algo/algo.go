package algo

import "github.com/haydenhigg/chrys"

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
