package chrys

import "time"

type Feed struct {
	Symbol   string
	Interval time.Duration
}

func NewFeed(symbol string, interval time.Duration) Feed {
	return Feed{
		Symbol:   symbol,
		Interval: interval,
	}
}
