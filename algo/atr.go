package algo

import (
	"github.com/haydenhigg/chrys"
	"math"
)

type ATR struct {
	Value     float64
	Period    int
	LastFrame *chrys.Frame
}

func NewATR(period int) *ATR {
	return &ATR{
		Period: period,
	}
}

func (atr *ATR) NextRaw(v float64) *ATR {
	atr.Value = (atr.Value + v) / float64(atr.Period)
	return atr
}

func (atr *ATR) Next(frame *chrys.Frame) *ATR {
	return atr.NextRaw(max(
		frame.High - frame.Low,
		math.Abs(frame.High - atr.LastFrame.Close),
		math.Abs(frame.Low - atr.LastFrame.Close),
	))
}
