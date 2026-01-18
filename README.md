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

### algo

Indicators and small composable “machines” for turning sequences of prices/frames into signals.

#### Machine

A stateful transformer that can be fed raw values or `*frame.Frame` and produces a single float output.

- **Methods:**
  - `Apply(x float64) Machine` — Update internal state with a scalar
  - `ApplyFrame(frame *frame.Frame) Machine` — Update internal state with a frame (typically uses `Close`)
  - `Val() float64` — Current output value

#### Composer

A simple chain/combinator for `Machine`s. It applies the most recently-added machine first, then feeds its output backward through earlier machines.

- **Fields:**
  - `Machines []Machine`
  - `Value float64`

- **Methods:**
  - `Of(machine Machine) *Composer` — Append a machine to the chain
  - `Apply(x float64) Machine` — Apply chained machines to a scalar and update `Value`
  - `ApplyFrame(frame *frame.Frame) Machine` — Apply chained machines to a frame and update `Value`
  - `Val() float64` — Current composed value

Create with:
```go
algo.NewComposer(initial algo.Machine) *algo.Composer
```

#### MA

Exponential-style running moving average (updates in O(1) time per tick).

- **Fields:**
  - `Period float64`
  - `Value float64`

- **Methods:**
  - `Apply(x float64) Machine`
  - `ApplyFrame(frame *frame.Frame) Machine`
  - `Val() float64`

Create with:
```go
algo.NewMA(period int) *algo.MA
```

#### EMA

Exponential moving average with default `Alpha = 2 / (1 + period)`.

- **Fields:**
  - `Alpha float64`
  - `Value float64`

- **Methods:**
  - `Apply(x float64) Machine`
  - `ApplyFrame(frame *frame.Frame) Machine`
  - `Val() float64`

Create with:
```go
algo.NewEMA(period int) *algo.EMA
```

#### MACD

Moving average convergence/divergence. Internally computes a fast EMA minus slow EMA, then runs that through a signal EMA and returns `signal - line` (a histogram-like value).

- **Fields:**
  - `Fast *EMA`
  - `Slow *EMA`
  - `Signal *EMA`
  - `Value float64`

- **Methods:**
  - `Apply(x float64) Machine`
  - `ApplyFrame(frame *frame.Frame) Machine`
  - `Val() float64`

Create with:
```go
algo.NewMACD(fastPeriod, slowPeriod, signalPeriod int) *algo.MACD
```

#### ATR

Average True Range. Computes true range per frame using the previous close, then smooths it using `MA`.

- **Fields:**
  - `Average *MA`
  - `LastClose float64`

- **Methods:**
  - `Apply(x float64) Machine` — Feed a precomputed true range value
  - `ApplyFrame(frame *frame.Frame) Machine` — Compute true range from the frame + prior close
  - `Val() float64` — Current ATR

Create with:
```go
algo.NewATR(period int) *algo.ATR
```

#### ZScore

Compute the z-score of the most recent element in a series.

- **Function:**
  - `ZScore(xs []float64) float64`

#### Utilities

Basic math and frame helpers.

- **Math:**
  - `Mean(xs []float64) float64`
  - `Variance(xs []float64, mean float64) float64`
  - `StandardDeviation(xs []float64, mean float64) float64`

- **Frames:**
  - `MapFrames(frames []*frame.Frame, processor func(*frame.Frame) float64) []float64`
  - `Opens(frames []*frame.Frame) []float64`
  - `Highs(frames []*frame.Frame) []float64`
  - `Lows(frames []*frame.Frame) []float64`
  - `Closes(frames []*frame.Frame) []float64`

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
