package algo

import "github.com/haydenhigg/chrys"

type MACD struct {
	Value  float64
	Fast   *EMA
	Slow   *EMA
	Signal *EMA
}

func NewMACD(fastPeriod, slowPeriod, signalPeriod int) *MACD {
	return &MACD{
		Fast:   NewEMA(fastPeriod),
		Slow:   NewEMA(slowPeriod),
		Signal: NewEMA(signalPeriod),
	}
}

func (macd *MACD) NextRaw(v float64) *MACD {
	line := macd.Fast.NextRaw(v).Value - macd.Slow.NextRaw(v).Value
	macd.Value = macd.Signal.NextRaw(line).Value - line
	return macd
}

func (macd *MACD) Next(frame *chrys.Frame) *MACD {
	return macd.NextRaw(frame.Close)
}
