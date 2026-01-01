# chrys
algorithmic trading toolbox

## principles
* **Simplicity**: setting up a successful strategy is easy and concise.
* **Composability**: pieces can be combined in novel ways without writing new code.
* **Flexibility**: all trading parameters and dynamics can be modified.

## to-do
1. improve organization (and simplify domain modeling)
    - split off FrameStore and BalanceStore from Client and pass the Connector to them so they can fetch the data themselves
    - add `.Alias(...)` to BalanceStore to track mappings between asset symbols and exchange specific asset codes
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

	client := chrys.NewClient(c).SetFee(0.004).SetIsLive(true)

	// set up strategy data
	btc := chrys.NewAsset("BTC", "XBT.F")
	usd := chrys.NewAsset("USD", "ZUSD")

	pair := chrys.NewPair(btc, usd)

	// set up pipeline
	pipeline := chrys.NewPipeline().AddStage(func(now time.Time) error {
		frames, err := client.GetNFramesBefore(pair, time.Hour, 20, now)
		if err != nil {
			return err
		}

		zScore := algo.ZScore(algo.Closes(frames))
		fmt.Println("BB(20) =", zScore)

		err = nil
		if zScore < -2 {
			err = client.Buy(pair, 0.10, now)
		} else if zScore > 2 {
			err = client.Sell(pair, 0.10, now)
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

---

### Asset

A currency.

- **Fields:**
  - `Symbol string`
  - `Code string` (for balance tracking on exchanges like Kraken that use proprietary asset codes)

Create with:
```go
chrys.NewAsset(symbol, code string) *Asset
```

---

### Pair

A tradeable currency pair.

- **Fields:**
  - `Base *Asset`
  - `Quote *Asset`
  - `Name string`

Create with:
```go
chrys.NewPair(base, quote *Asset) *Pair
```

---

### Pipeline

A stateful function-chaining pipeline for building strategies.

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

Create with:
```go
chrys.NewPipeline() *Pipeline
```
