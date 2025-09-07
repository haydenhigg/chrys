package algo

import "github.com/haydenhigg/chrys"

type Composable interface {
	Apply(x float64) Composable
	ApplyFrame(frame *chrys.Frame) Composable
	Val() float64
}

type Composed struct {
	A     Composable
	B     Composable
	Value float64
}

func (composed *Composed) Apply(x float64) Composable {
	composed.Value = composed.A.Apply(composed.B.Apply(x).Val()).Val()
	return composed
}

func (composed *Composed) ApplyFrame(frame *chrys.Frame) Composable {
	composed.Value = composed.A.Apply(composed.B.ApplyFrame(frame).Val()).Val()
	return composed
}

func (composed *Composed) Val() float64 {
	return composed.Value
}

func Compose(a, b Composable) Composable {
	return &Composed{A: a, B: b}
}
