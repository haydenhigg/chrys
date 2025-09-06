package algo

import "github.com/haydenhigg/chrys"

type EMA struct {
	Value  float64
	Period int
	Alpha  float64 // 2 / (1 + period)
}

func NewEMA(period int) *EMA {
	return &EMA{
		Period: period,
		Alpha:  2 / (1 + float64(period)),
	}
}

func (ema *EMA) NextRaw(v float64) *EMA {
	ema.Value = v*ema.Alpha + ema.Value*(1-ema.Alpha)
	return ema
}

func (ema *EMA) Next(frame *chrys.Frame) *EMA {
	return ema.NextRaw(frame.Close)
}
