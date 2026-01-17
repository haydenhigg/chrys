package algo

import domain "github.com/haydenhigg/chrys/frame"

type MACD struct {
	Fast   *EMA
	Slow   *EMA
	Signal *EMA
	Value  float64
}

func NewMACD(fastPeriod, slowPeriod, signalPeriod int) *MACD {
	return &MACD{
		Fast:   NewEMA(fastPeriod),
		Slow:   NewEMA(slowPeriod),
		Signal: NewEMA(signalPeriod),
	}
}

func (macd *MACD) Apply(x float64) Machine {
	line := macd.Fast.Apply(x).Val() - macd.Slow.Apply(x).Val()
	macd.Value = macd.Signal.Apply(line).Val() - line
	return macd
}

func (macd *MACD) ApplyFrame(frame *domain.Frame) Machine {
	return macd.Apply(frame.Close)
}

func (macd *MACD) Val() float64 {
	return macd.Value
}
