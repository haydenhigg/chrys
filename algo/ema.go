package algo

import "github.com/haydenhigg/chrys"

type EMA struct {
	Value  float64
	Period int
	alpha  float64 // 2 / (1 + period)
}

func NewEMA(period int) *EMA {
	return &EMA{
		Period: period,
		alpha:  2 / (1 + float64(period)),
	}
}

func (ema *EMA) NextRaw(v float64) *EMA {
	ema.Value = v*ema.alpha + ema.Value*(1-ema.alpha)
	return ema
}

func (ema *EMA) Next(frame *chrys.Frame) *EMA {
	return ema.NextRaw(frame.Close)
}
