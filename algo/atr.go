package algo

import (
	"github.com/haydenhigg/chrys"
	"math"
)

type ATR struct {
	Value     float64
	Period    int
	LastClose float64
}

func NewATR(period int) *ATR {
	return &ATR{
		Period:    period,
		LastClose: 0,
	}
}

func (atr *ATR) NextRaw(v float64) *ATR {
	atr.Value = (atr.Value + v) / float64(atr.Period)
	return atr
}

func (atr *ATR) Next(frame *chrys.Frame) *ATR {
	atr.NextRaw(max(
		frame.High-frame.Low,
		math.Abs(frame.High-atr.LastClose),
		math.Abs(frame.Low-atr.LastClose),
	))
	atr.LastClose = frame.Close
	return atr
}
