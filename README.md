# chrys

algorithmic trading framework for medium-frequency strategies

## principles

1. **Reliable**: It's robust enough to trust with serious sums of money.
2. **Simple**: Setting up an effective strategy is concise and intuitive.
3. **Flexible**: All functionality can be modified or composed.

## to-do

1. create backtest machinery.
    - annualize Volatility, Sharpe, Sortino properly using the Hurst exponent
    - parallelized monte carlo sub-period resampling
2. create optimizer.
    - random sampling
    - TPE
3. implement k-fold backtesting. more important with an auto-optimizer.
4. complete documentation.
5. add more algos.
    - RSI
    - MFI
7. create OrderStore.
    - internal `ledger`
7. create plug-ins (concise, reusable signals).
8. simplify driver interface.
    - replace `driver.FetchFramesSince` with `driver.FetchNFrames`
    - remove `frames.GetSince`
    - since you now have `now` when checking the cache, you can check whether the latest Frame is older than now-interval, in which case the cache is stale and you need to re-fetch anyway. this will make it more robust

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
