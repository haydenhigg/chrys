package chrys

import (
	"fmt"
	"time"
)

// Frame
type Frame struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

// Asset
type Asset struct {
	Symbol string
	Code   string // exchange-specific, for balances
}

func NewAsset(symbol, code string) *Asset {
	return &Asset{
		Symbol: symbol,
		Code:   code,
	}
}

// Pair
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
