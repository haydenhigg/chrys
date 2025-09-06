package chrys

import (
	"fmt"
	"time"
)

// domain-modeling utilities
type Asset struct {
	Symbol string
	Code   string // for trading
}

func NewAsset(symbol, code string) *Asset {
	return &Asset{
		Symbol: symbol,
		Code:   code,
	}
}

type Pair struct {
	Base  *Asset
	Quote *Asset
	Name  string
}

func NewPair(base, quote *Asset) *Pair {
	return &Pair{
		Base:  base,
		Quote: quote,
		Name:  fmt.Sprintf("%s/%s", base.Symbol, quote.Symbol),
	}
}

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
