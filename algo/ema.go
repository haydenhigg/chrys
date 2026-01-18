package algo

import domain "github.com/haydenhigg/chrys/frame"

type EMA struct {
	Alpha float64 // 2 / (1 + period)
	Value float64
}

func NewEMA(period int) *EMA {
	return &EMA{Alpha: 2 / (1 + float64(period))}
}

func (ema *EMA) Apply(x float64) Machine {
	ema.Value = x*ema.Alpha + ema.Value*(1-ema.Alpha)
	return ema
}

func (ema *EMA) ApplyFrame(frame *domain.Frame) Machine {
	return ema.Apply(frame.Close)
}

func (ema *EMA) Val() float64 {
	return ema.Value
}
