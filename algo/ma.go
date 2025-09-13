package algo

import "github.com/haydenhigg/chrys"

type MA struct {
	Period float64
	Value  float64
}

func NewMA(period int) *MA {
	return &MA{
		Period: float64(period),
	}
}

func (ma *MA) Apply(x float64) Machine {
	ma.Value = (x + ma.Value*(ma.Period-1)) / ma.Period
	return ma
}

func (ma *MA) ApplyFrame(frame *chrys.Frame) Machine {
	return ma.Apply(frame.Close)
}

func (ma *MA) Val() float64 {
	return ma.Value
}
