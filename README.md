# chrys
algorithmic trading toolbox

## principles
* **Simplicity**: setting up a successful strategy is easy and concise.
* **Composability**: pieces can be combined in novel ways without writing new code.
* **Flexibility**: all trading parameters and dynamics can be modified.

## to-do
1. improve organization (and simplify domain modeling)
2. unit tests
3. backtest machinery
    - write `(pipeline *Pipeline) RunBetween(start, end time.Time) error`
4. backtest metrics
    - volatility
    - Sharpe ratio
5. add/test more algos
    - ROC
    - ADI
    - MFI
    - make ZScore a Machine?
    - make TrueRange a Machine?

## example
This trades on **BOLL(20, 2)** signals for **1h BTC/USD** using a **10%** fractional trade amount.

```go
package main

import (
	"fmt"
	"github.com/haydenhigg/chrys"
	"github.com/haydenhigg/chrys/algo"
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
	order := chrys.NewOrder(pair, 0.10).SetIsLive(true) // ±10%

	// set up pipeline
	pipeline := chrys.NewPipeline().AddStage(func(now time.Time) error {
		frames, err := client.GetFrames(pair, time.Hour, now, 20)
		if err != nil {
			return err
		}

		zScore := algo.ZScore(algo.Closes(frames))
		fmt.Println("BB(20) =", zScore)

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

## API

### Frame

Represents a single candle of OHLCV data.

- **Fields:**
  - `Time time.Time`
  - `Open float64`
  - `High float64`
  - `Low float64`
  - `Close float64`
  - `Volume float64`

---

### Asset

Create with:
```go
chrys.NewAsset(symbol, code string) *Asset
```
Represents an asset with a human-readable symbol and an exchange-specific code.

- **Fields:**
  - `Symbol string`
  - `Code string`

---

### Pair

Create with:
```go
chrys.NewPair(base, quote *Asset) *Pair
```
Represents a trading pair.

- **Fields:**
  - `Base *Asset`
  - `Quote *Asset`
  - `Name string`

---

### Order

Create with:
```go
chrys.NewOrder(pair *Pair, percent float64) *Order
```
Describes an order configuration on a pair.

- **Fields:**
  - `Pair *Pair`
  - `Percent float64` — Fraction of portfolio to buy/sell, e.g., 0.10 for 10%
  - `IsLive bool` — Execute order live or in simulation
  - `Type OrderType` — One of `chrys.MARKET_BUY`, `chrys.MARKET_SELL`, etc.

- **Methods:**
  - `SetIsLive(isLive bool) *Order` — Enable/disable live mode
  - `SetBuy() *Order` — Set as buy order
  - `SetSell() *Order` — Set as sell order

---

### Client

Create with:
```go
chrys.NewClient(connector Connector) *Client
```
Manages caching, balances, and calling trading connector.

- **Fields:**
  - `Connector Connector`
  - `FrameCache map[string]map[time.Duration][]*Frame`
  - `Balances map[string]float64`
  - `Fee float64` — Trading fee as a decimal

- **Methods:**
  - `SetFee(fee float64) *Client` — Set per-trade fee
  - `GetFramesSince(pair *Pair, interval time.Duration, t time.Time) ([]*Frame, error)` — Retrieve frames before a timestamp
  - `GetFrames(pair *Pair, interval time.Duration, t time.Time, n int) ([]*Frame, error)` — Retrieve `n` frames before timestamp
  - `GetBalances() (map[string]float64, error)` — Get asset balances
  - `PlaceOrder(order *Order, t time.Time) error` — Place order at specified time

---

### Pipeline

Create with:
```go
chrys.NewPipeline() *Pipeline
```
Stateful function-chaining pipeline for building strategy evaluations.

- **Fields:**
  - `Data map[string]float64`
  - `Stages []Stage`

- **Types:**
  - `type Stage = func(now time.Time) error`

- **Methods:**
  - `Get(k string) float64` — Retrieve value from pipeline data store
  - `Set(k string, v float64) *Pipeline` — Set value in data store
  - `AddStage(handler Stage) *Pipeline` — Add a stage (function) to process
  - `Run(t time.Time) error` — Process all stages in order
