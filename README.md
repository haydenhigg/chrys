# chrys
lightweight algorithmic trading framework

## to-do
1. algo state management (through `Pipeline`? or its own component?)
2. backtest machinery
    - add `(client *Client) CalculateEquity(out *Asset, t time.Time) (float64, error)`
    - add `(pipeline *Pipeline) RunBacktest`
    - add more backtesting metrics (volatility, Sharpe ratio)
3. add/test more algos
    - ADI
    - MFI
    - make ZScore Incremental (`interface { NextRaw(v float64) Incremental; Next(frame *chrys.Frame) Incremental }`)
4. expand MLP implementation
5. plug-ins

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

	// set up strategy data
	pair := chrys.NewPair(
		chrys.NewAsset("BTC", "XBT.F"),
		chrys.NewAsset("USD", "ZUSD"),
	)

	series := chrys.NewSeries(pair, time.Hour)
	order := chrys.NewOrder(pair, 0.10).SetIsLive(true) // ±10%

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

## Frame
a frame of TOHLCV data (a "candle")
- `Time time.Time`
- `Open float64`
- `High float64`
- `Low float64`
- `Close float64`
- `Volume float64`

## Pipeline
`chrys.NewPipeline() *Pipeline`

a stateful function pipeline
- `Data map[string]float64`
- `Stages []Stage`

#### functions
- `Get(k string) float64`
- `Set(k string, v float64) *Pipeline`
- `AddStage(handler Stage) *Pipeline`
- `Run(t time.Time) error` (processes stages in order)

#### types
- `type Stage = func(now time.Time) error`

## Asset
`chrys.NewAsset(symbol, code string) *Asset`

an asset with a human-readable symbol and an exchange-specific code
- `Symbol string`
- `Code string`

## Pair
`chrys.NewPair(base, quote string) *Pair`

a pair with a human-readable name
- `Base *Asset`
- `Quote *Asset`
- `Name string`

## Series
`chrys.NewSeries(pair *Pair, interval time.Duration) *Series`

a `Pair` and an interval to get a chartable series
- `Pair *Pair`
- `Interval time.Duration`

## Order
`chrys.NewOrder(pair *Pair, percent float64) *Order`

a `Pair` and order configuration details
- `Pair *Pair`
- `Percent float64`
- `IsLive bool`
- `Type OrderType` (e.g., `chrys.BUY`, `chrys.SELL`)

#### functions
- `SetIsLive(isLive bool) *Order`
- `SetBuy() *Order`
- `SetSell() *Order`

## Client
`chrys.NewClient(connector Connector) *Client`

a caching client for connectors
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
