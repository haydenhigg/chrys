# chrys
lightweight algorithmic trading framework

## concepts
1. *chrys.Frame*: a frame of TOHLCV data
2. *chrys.Pair*: an asset pair
2. *chrys.Pipeline*: a stateful function pipeline

## /client
Provides *client.Client*, a stateful cache layer over third-party APIs.

## /algo
Provides a collection of frame manipulation and mathematical utilities.

## to-do
0.5 continue tidying
1. change GetFrames to use config like Order
2. backtesting
3. state management through `algo` subpackage for stateful algos
4. add built-in logging to client
5. expand MLP implementation

## usage
This trades on **BOLL(20, 2)** signals for **1h BTC/USD** using a **10%** fractional trade amount.

```go
package main

import (
	"github.com/haydenhigg/chrys"
	"github.com/haydenhigg/chrys/algo"
	"github.com/haydenhigg/chrys/client"
	"fmt"
	"os"
	"time"
)

func main() {
	// set up client
	c, err := client.NewKraken(os.Getenv("API_KEY"), os.Getenv("API_SECRET"))
	if err != nil {
		panic(err)
	}

	pair := chrys.NewPair("BTC", "USD").SetBalancePair("XBT.F", "ZUSD")
	feed := chrys.NewFeed(pair.String(), time.Hour)

	var orderConfig *client.OrderConfig

	// set up pipeline
	pipeline := chrys.NewPipeline().AddStage(func(now time.Time) error {
		frames, err := c.GetFrames(feed, now, 20)
		if err != nil {
			return err
		}

		zScore := algo.ZScore(algo.Closes(frames))

		orderConfig = &client.OrderConfig{
			Pair:    pair,
			Percent: 0.1,
		}

		if zScore < -2 {
			orderConfig.Side = client.MARKET_BUY
		} else if zScore > 2 {
			orderConfig.Side = client.MARKET_SELL
		} else {
			return nil
		}

		if err := c.PlaceOrder(orderConfig, now); err != nil {
			return err
		}

		return nil
	})

	// run
	if err := pipeline.Run(time.Now()); err != nil {
		panic(err)
	}
}
```

For more complex use cases, you'll potentially want to split the signal and order logic into separate stages. You can pass data down through the stages chain like so:

```go
pipeline := chrys.NewPipeline()
pipeline.AddStage(func(now time.Time) error {
	frames, err := c.GetFrames(feed, now, 20)
	if err != nil {
		return err
	}

	closes := algo.Closes(frames)

	pipeline.Set("bb", algo.ZScore(closes))
	pipeline.Set("bb-1/2", algo.ZScore(closes[10:]))

	return nil
})
pipeline.AddStage(func(now time.Time) error {
	fmt.Println(pipeline.Data) // map[bb:... bb-1/2:...]

	// ...

	if pipeline.Get("bb") < -2 && pipeline.Get("bb-1/2") < -2 {
		order.Side = "buy"
	} else if pipeline.Get("bb") > 2 && pipeline.Get("bb-1/2") > 2 {
		order.Side = "sell"
	}

	// ...
})
```

Stages are processed sequentially in the order that they are defined.
