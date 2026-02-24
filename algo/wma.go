package algo

import domain "github.com/haydenhigg/chrys/frame"

type WMA struct {
	Period float64
	Value  float64
}

func NewWMA(period int) *WMA {
	return &WMA{Period: float64(period)}
}

func (wma *WMA) Apply(x float64) Machine {
	wma.Value += (x - wma.Value) / wma.Period
	return wma
}

func (wma *WMA) ApplyFrame(frame *domain.Frame) Machine {
	return wma.Apply(frame.Close)
}

func (wma *WMA) Val() float64 {
	return wma.Value
}
