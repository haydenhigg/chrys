package chrys

import "time"

type Series struct {
	Pair     *Pair
	Interval time.Duration
}

func NewSeries(pair *Pair, interval time.Duration) *Series {
	return &Series{
		Pair:        pair,
		Interval:    interval,
	}
}
