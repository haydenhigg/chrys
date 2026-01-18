# chrys
algorithmic trading toolbox

## principles
* **Simplicity**: setting up a successful strategy is easy and concise.
* **Composability**: pieces can be combined in novel ways without writing new code.
* **Flexibility**: all trading parameters and dynamics can be modified.

## to-do
1. improve Frames
    - [x] ~~if more than `interval` time has passed since the last frame in the cache, assume that the cache is stale and refetch~~ **It's not possible to do this because the cache miss at the end will need to use `time.Now()` as a reference, which will not work for the HistoricalDriver (since its data often does not go up to the present moment). So the cache will never miss at the end; in other words, the cache should only be used within the interval that it was created. chrys is not for long-running programs.**
    - [ ] ~~in the above case, only retrieve what's new and needed (`since = frames[len(frames)-1].Time + interval`) instead of retrieving all the overlapping frames~~
    - [x] use a binary search to find where to chop off older cached frames
    - [ ] unit test
2. unit tests
    - [ ] FrameStore
    - [ ] BalanceStore
    - [ ] Client
    - [ ] Pipeline
    - [ ] algo
    - [ ] connector
3. backtest machinery + reporting

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
