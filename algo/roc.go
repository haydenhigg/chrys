package algo

import domain "github.com/haydenhigg/chrys/frame"

type ROC struct {
	IsInitial bool
	LastPrice float64
	Value     float64
}

func NewROC() *ROC {
	return &ROC{IsInitial: true}
}

func (roc *ROC) Apply(x float64) Machine {
	if !roc.IsInitial {
		roc.Value = x/roc.LastPrice - 1
	} else {
		roc.IsInitial = false
	}

	roc.LastPrice = x

	return roc
}

func (roc *ROC) ApplyFrame(frame *domain.Frame) Machine {
	return roc.Apply(frame.Close)
}

func (roc *ROC) Val() float64 {
	return roc.Value
}
