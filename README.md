# chrys

algorithmic trading toolbox for medium-frequency strategies

## principles

1. **Simplicity**: setting up a successful strategy is concise and intuitive.
2. **Flexibility**: all trading parameters and dynamics can be modified.
3. **Composability**: pieces can be combined in novel ways without writing new code.

## to-do

1. create backtest machinery.
    - interval argument in `pipeline.AddBlock`
    - OrderStore with an internal ledger
    - a Backtester subpackage
2. write key unit test suites.
    - [x] FrameStore
    - [x] BalanceStore
    - [ ] OrderStore
    - [x] Client
    - [ ] Pipeline
3. optimize. backtests over 1hr interval for 1yr are slow -- why?
4. change driver interface.
    - replace `driver.FetchFramesSince` with `driver.FetchNFrames`
    - remove `client.GetFramesSince`
5. clean up and complete documentation.

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

	client.Balances.Alias("BTC", "XXBT").Alias("USD", "ZUSD")

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
	now := time.Now()
	if err := pipeline.Run(now); err != nil {
		panic(err)
	}

	// print portfolio value
	value, err := client.TotalValue([]string{"USD", "BTC"}, now)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Portfolio Value: $%.2f\n", value)
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
