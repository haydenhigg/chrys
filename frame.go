package chrys

import (
	"fmt"
	"time"
)

type Frame struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

func (frame *Frame) String() string {
	return fmt.Sprintf(
		"Frame{T: %s, O: %f, H: %f, L: %f, C: %f, V: %f}",
		frame.Time,
		frame.Open,
		frame.High,
		frame.Low,
		frame.Close,
		frame.Volume,
	)
}
