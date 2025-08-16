package candle

import "time"

type Candle struct {
	Time time.Time
	Open,
	High,
	Low,
	Close,
	Volume float64
}
