# chrys
algorithmic trading toolbox

## principles
* **Simplicity**: setting up a successful strategy is easy and concise.
* **Composability**: pieces can be combined in novel ways without writing new code.
* **Flexibility**: all trading parameters and dynamics can be modified.

## notes
FrameStore caches frames, but it does not try to check if new frames are available from the data source if frames already exist in the cache. So, you should not trust the output of the cache if your program runs longer than the interval you use.

## to-do
1. backtest machinery + reporting
2. unit tests
    - [x] FrameStore
    - [x] BalanceStore
    - [ ] Client
    - [ ] Pipeline
    - [ ] algo
    - [ ] connector

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

---

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
