package algo

type EMA struct {
	Value  float64
	Period int
	alpha  float64 // 2 / (1 + period)
}

func NewEMA(period int) *EMA {
	return &EMA{
		Period: period,
		alpha: 2 / (1 + float64(period)),
	}
}

func (ema *EMA) Next(v float64) *EMA {
	ema.Value = v * ema.alpha + v * (1 - ema.alpha)
	return ema
}
