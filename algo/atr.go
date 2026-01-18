package algo

import (
	domain "github.com/haydenhigg/chrys/frame"
	"math"
)

type ATR struct {
	Average   *MA
	LastClose float64
}

func NewATR(period int) *ATR {
	return &ATR{Average: NewMA(period)}
}

func (atr *ATR) Apply(x float64) Machine {
	atr.Average.Apply(x)
	return atr
}

func (atr *ATR) ApplyFrame(frame *domain.Frame) Machine {
	atr.Apply(max(
		frame.High-frame.Low,
		math.Abs(frame.High-atr.LastClose),
		math.Abs(frame.Low-atr.LastClose),
	))
	atr.LastClose = frame.Close
	return atr
}

func (atr *ATR) Val() float64 {
	return atr.Average.Val()
}
