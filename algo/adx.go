package algo

import (
	domain "github.com/haydenhigg/chrys/frame"
	"math"
)

type ADX struct {
	ATR      *ATR
	PosDM    *WMA
	NegDM    *WMA
	Average  *MA
	LastLow  float64
	LastHigh float64
}

func NewADX(period int) *ADX {
	return &ADX{
		ATR:     NewATR(period),
		Average: NewMA(period),
	}
}

func (adx *ADX) Apply(x float64) Machine {
	adx.Average.Apply(x)
	return adx
}

func (adx *ADX) ApplyFrame(frame *domain.Frame) Machine {
	pos := frame.High - adx.LastHigh
	neg := adx.LastLow - frame.Low

	if pos > neg {
		adx.PosDM.Apply(max(pos, 0))
	} else if neg > pos {
		adx.NegDM.Apply(max(neg, 0))
	}

	atr := adx.ATR.ApplyFrame(frame).Val()
	posDI := adx.PosDM.Val() / atr
	negDI := adx.NegDM.Val() / atr

	adx.Apply(100 * math.Abs((posDI-negDI)/(posDI+negDI)))
	adx.LastLow = frame.Low
	adx.LastHigh = frame.High

	return adx
}

func (adx *ADX) Val() float64 {
	return adx.Average.Val()
}
