package chrys

import (
	"fmt"
	"math"
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
	Code   string // for balances
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

// Order
type OrderType string

const (
	MARKET_BUY  OrderType = "buy"
	MARKET_SELL           = "sell"
)

type Order struct {
	Pair    *Pair
	Percent float64
	IsLive  bool
	Type    OrderType
}

func NewOrder(pair *Pair, percent float64) *Order {
	return &Order{
		Pair:    pair,
		Percent: math.Min(math.Max(percent, 0), 1), // clamped to [0,1]
	}
}

func (order *Order) SetIsLive(isLive bool) *Order {
	order.IsLive = isLive
	return order
}

func (order *Order) SetBuy() *Order {
	order.Type = MARKET_BUY
	return order
}

func (order *Order) SetSell() *Order {
	order.Type = MARKET_SELL
	return order
}
