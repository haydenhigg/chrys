package algo

import domain "github.com/haydenhigg/chrys/frame"

type Supertrend struct {
	ATR       *ATR
	Mult      float64
	LastUpper float64
	LastLower float64
	LastClose float64
	Value     float64
}

func NewSupertrend(period int, mult float64) *Supertrend {
	return &Supertrend{
		ATR:  NewATR(period),
		Mult: mult,
	}
}

func (supertrend *Supertrend) Apply(x float64) Machine {
	supertrend.Value = x
	return supertrend
}

func (supertrend *Supertrend) ApplyFrame(frame *domain.Frame) Machine {
	hl2 := (frame.High + frame.Low) / 2
	atr := supertrend.ATR.ApplyFrame(frame).Val()

	basicUpper := hl2 + supertrend.Mult*atr
	basicLower := hl2 - supertrend.Mult*atr

	upper := supertrend.LastUpper
	if basicUpper < upper || supertrend.LastClose > upper {
		upper = basicUpper
	}

	lower := supertrend.LastLower
	if basicLower > lower || supertrend.LastClose < lower {
		lower = basicLower
	}

	switch supertrend.Value {
	case upper:
		if frame.Close > upper {
			supertrend.Apply(lower)
		} else {
			supertrend.Apply(upper)
		}
	case lower:
		if frame.Close < lower {
			supertrend.Apply(upper)
		} else {
			supertrend.Apply(lower)
		}
	}

	return supertrend
}

func (supertrend *Supertrend) Val() float64 {
	return supertrend.Value
}
