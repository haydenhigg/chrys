# chrys

algorithmic trading framework for medium-frequency (>=1min) strategies

## principles

1. **Reliable**: It's robust enough to handle serious sums of money.
2. **Simple**: Setting up an effective strategy is concise and intuitive.
3. **Flexible**: All functionality can be modified or composed.

## to-do

1. create backtest machinery.
    - other functions like `backtest.TotalReturn()` etc.
    - accept multiple evaluator functions?
2. explore optimizations. 30s backtest for simple 1hr strategy over 1yr -- why?
    - is it that `algo.Closes` uses `algo.MapFrames`? that's a lot of closures to create and then garbage-collect.
    - is it that non-Machines like ZScore are fetching and processing much of the same data every time?
    - is it that `HistoricalDriver` is inefficient and slow?
3. change driver interface.
    - replace `driver.FetchFramesSince` with `driver.FetchNFrames`
    - remove `frames.GetSince`
    - since you now have `now` when checking the cache, you can check whether the latest Frame is older than now-interval, in which case the cache is stale and you need to re-fetch anyway. this will make it more robust
4. complete documentation.
5. create OrderStore.
    - internal `ledger`
6. create plug-ins (concise, reusable signals).

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

	// set up scheduler
	scheduler := chrys.NewScheduler()
	scheduler.Add(time.Minute, func(now time.Time) error {
		// print portfolio value
		value, err := client.Value([]string{"USD", "BTC"}, now)
		if err == nil {
			fmt.Printf("Portfolio value: $%.2f\n", value)
		}

		// get frames
		frames, err := client.Frames.GetNBefore("BTC/USD", time.Minute, 20, now)
		if err != nil {
			return err
		}

		// calculate signal and place order if necessary
		zScore := algo.ZScore(algo.Closes(frames))
		if zScore < -2 {
			err = client.Buy("BTC/USD", 0.10, now)
		} else if zScore > 2 {
			err = client.Sell("BTC/USD", 0.10, now)
		}

		return err
	})

	// run
	if err := scheduler.Run(time.Now()); err != nil {
		panic(err)
	}
}
```

## documentation

### Frame

A single OHLC frame.

- **Fields:**
  - `Time time.Time`
  - `Open float64`
  - `High float64`
  - `Low float64`
  - `Close float64`
  - `Volume float64`
