package chrys

import (
	"fmt"
	"time"
)

type Series struct {
	Pair     *Pair
	Interval time.Duration
}

func NewSeries(pair *Pair, interval time.Duration) *Series {
	return &Series{
		Pair:     pair,
		Interval: interval,
	}
}

func (series *Series) String() string {
	return fmt.Sprintf("%s:%s", series.Pair, series.Interval)
}
