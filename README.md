# chrys

algorithmic trading toolbox for medium-frequency strategies

## principles

1. **Simplicity**: setting up a successful strategy is easy and concise.
2. **Composability**: pieces can be combined in novel ways without writing new code.
3. **Flexibility**: all trading parameters and dynamics can be modified.

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
	client, err := chrys.NewKrakenClient(
		os.Getenv("API_KEY"),
		os.Getenv("API_SECRET"),
	)
	if err != nil {
		panic(err)
	}

	client.Balances.Alias("BTC", "XBT.F").Alias("USD", "ZUSD")

	// set up pipeline
	pipeline := chrys.NewPipeline()
	pipeline.AddBlock(func(now time.Time) error {
		frames, err := client.Frames.GetNBefore("BTC/USD", time.Hour, 20, now)
		if err != nil {
			return err
		}

		zScore := algo.ZScore(algo.Closes(frames))
		fmt.Println("BB(20) =", zScore)

		err = nil
		if zScore < -2 {
			err = client.Buy("BTC/USD", 0.10, now)
		} else if zScore > 2 {
			err = client.Sell("BTC/USD", 0.10, now)
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

A single OHLCV candle.

- **Fields:**
  - `Time time.Time`
  - `Open float64`
  - `High float64`
  - `Low float64`
  - `Close float64`
  - `Volume float64`

### Pipeline

A stateful function-chaining pipeline for building strategies.

- **Fields:**
  - `Blocks []Block`

- **Types:**
  - `type Block = func(now time.Time) error`

- **Methods:**
  - `AddBlock(handler Block) *Pipeline` — Add a stage (function) to process
  - `Run(t time.Time) error` — Process all stages in order

Create with:
```go
chrys.NewPipeline() *Pipeline
```

### Client

A trading client that wraps an `API` and provides cached access to frames and balances, plus convenience order helpers.

- **Fields:**
  - `Frames *store.FrameStore` — cached OHLCV access
  - `Balances *store.BalanceStore` — cached balances + aliasing
  - `Fee float64` — applied to simulated balance updates
  - `IsLive bool` — if `true`, places real orders via the underlying `API`

- **Methods:**
  - `SetFee(fee float64) *Client` — Set fee rate used when updating balances
  - `SetIsLive(isLive bool) *Client` — Toggle live order placement
  - `Order(side OrderSide, pair string, percent float64, t time.Time) error` — Place a market order by fractional sizing
  - `Buy(pair string, percent float64, t time.Time) error` — Convenience wrapper for `Order(BUY, ...)`
  - `Sell(pair string, percent float64, t time.Time) error` — Convenience wrapper for `Order(SELL, ...)`

Create with:
```go
chrys.NewClient(api chrys.API) *chrys.Client
chrys.NewKrakenClient(key, secret string) (*chrys.Client, error) // default Fee is 0.4%
```

Notes:
- Order `percent` is clamped to `[0, 1]`.
- Order sizing is based on `percent * balances[base]`, using `Frames.GetPriceAt(pair, t)` to ensure you don’t overspend the quote balance on buys.
- If `IsLive == false`, no exchange order is placed; balances are still updated locally (useful for paper trading).

### API

Low-level interface a driver must implement. `Client` depends on this.

- **Embedded interfaces:**
  - `store.BalanceAPI`
  - `store.FrameAPI`

- **Methods:**
  - `MarketOrder(side, pair string, quantity float64) error`

### OrderSide

Side of an order.

- **Type:**
  - `type OrderSide string`

- **Constants:**
  - `BUY`
  - `SELL`

### FromJSONFile / ToJSONFile

Minimal helpers for persisting state to disk.

- **Functions:**
  - `FromJSONFile(name string, v any) error` — Read JSON from `name` into `v`
  - `ToJSONFile(name string, v any) error` — Write `v` as JSON to `name`

---

## to-do

1. backtest machinery + reporting
2. unit tests
    - [x] FrameStore
    - [x] BalanceStore
    - [ ] Client
    - [ ] Pipeline
    - [ ] algo
    - [ ] connector
