# chrys
lightweight algorithmic trading framework

## composable building blocks
- **chrys.Frame** - a frame of TOHLCV data (a "candle")
- **chrys.Pair** - a tradeable pair with customizable asset codes
- **chrys.Series** - a `chrys.Pair` and an interval
- **chrys.Order** - a `chrys.Pair` and order configuration details
- **chrys.Client** - a caching client for connectors
- **chrys.Pipeline** - a stateful function pipeline

### Frame
#### fields
- `Time   time.Time`
- `Open   float64`
- `High   float64`
- `Low    float64`
- `Close  float64`
- `Volume float64`

### Pair
*constructor*: `chrys.NewPair(base, quote string) &Pair`

#### fields
- `Symbol string`
- `BaseCode string` (defaults to `base`)
- `QuoteCode string` (defaults to `quote`)

#### functions
- `Base() string`
- `Quote() string`
- `SetBaseCode(baseCode string) *Pair`
- `SetQuoteCode(quoteCode string) *Pair`
- `SetCodes(baseCode, quoteCode string) *Pair`

### Series
*constructor*: `chrys.NewSeries(pair *Pair, interval time.Duration) *Series`

#### fields
- `Pair *Pair`
- `Interval time.Duration`

### Order
*constructor*: `chrys.NewOrder(pair *Pair, percent float64, isLive bool) *Order`

#### fields
- `Pair *Pair`
- `Percent float64`
- `IsLive bool`
- `Type OrderType` (e.g., `chrys.BUY`, `chrys.SELL`)

#### functions
- `SetBuy() *Order`
- `SetSell() *Order`

### Client
*constructor*: `chrys.NewClient(connector Connector) *Client`

#### fields
- `Connector Connector`
- `FrameCache map[string]map[time.Duration][]*Frame`
- `Balances map[string]float64`
- `Fee float64`

#### functions
- `SetFee(fee float64) *Client`
- `GetFramesSince(series *Series, t time.Time) ([]*Frame, error)`
- `GetFrames(series *Series, t time.Time, n int) ([]*Frame, error)`
- `GetBalances() (map[string]float64, error)`
- `PlaceOrder(order *Order, t time.Time) error`

### Pipeline
*constructor*: `chrys.NewPipeline() *Pipeline`

#### fields
- `Stages []func(now time.Time) error`
- `Data map[string]float64`

#### functions
- `AddStage(handler func(now time.Time) error) *Pipeline`
- `Get(k string) float64`
- `Set(k string, v float64) *Pipeline`
- `Run(t time.Time) error` (processes stages in order)

## to-do
- save time on `Client` so you can do `client.AtTime(t).`...

1. tidying and API improvements
2. backtesting components
3. algo state management components
5. add built-in logging to client
6. expand MLP implementation
6. plug-ins

## example
This trades on **BOLL(20, 2)** signals for **1h BTC/USD** using a **10%** fractional trade amount.

```go
package main

import (
	"fmt"
	"github.com/haydenhigg/chrys"
	"github.com/haydenhigg/chrys/algo"
	"github.com/haydenhigg/chrys/connector"
	"os"
	"time"
)

func main() {
	// set up client
	c, err := connector.NewKraken(os.Getenv("API_KEY"), os.Getenv("API_SECRET"))
	if err != nil {
		panic(err)
	}
	client := chrys.NewClient(c).SetFee(0.004)

	// set up strategy-specific data objects
	pair := chrys.NewPair("BTC", "USD").SetCodes("XBT.F", "ZUSD")

	series := chrys.NewSeries(pair, time.Hour)
	order := chrys.NewOrder(pair, 0.10).SetLive(true) // ±10% live

	// set up pipeline
	pipeline := chrys.NewPipeline().AddStage(func(now time.Time) error {
		frames, err := client.GetFrames(series, now, 20)
		if err != nil {
			return err
		}

		zScore := algo.ZScore(algo.Closes(frames))
		fmt.Println("BB =", zScore)

		err = nil
		if zScore < -2 {
			err = client.PlaceOrder(order.SetBuy(), now)
		} else if zScore > 2 {
			err = client.PlaceOrder(order.SetSell(), now)
		}

		return err
	})

	// run
	if err := pipeline.Run(time.Now()); err != nil {
		panic(err)
	}
}
```
