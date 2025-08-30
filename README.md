# chrys
lightweight algorithmic trading framework

## concepts
- **chrys.Frame**: a frame of TOHLCV data
- **chrys.Pair**: a tradeable pair with customizable asset codes
- **chrys.Series**: a `chrys.Pair` and an interval
- **chrys.Pipeline**: a stateful function pipeline

## to-do
- save time in `Client` so you can do client.AtTime(t)...

## upcoming
1. tidying and API improvements
2. backtesting components
3. algo state management components
4. add built-in logging to client
5. expand MLP implementation

## usage
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
	client := chrys.NewClient(c).SetFee(0.004)

	// set up strategy-specific data objects
	pair := chrys.NewPair("BTC", "USD").SetCodes("XBT.F", "ZUSD")

	series := chrys.NewSeries(pair, time.Hour)
	order := chrys.NewOrder(pair, 0.10).SetLive(true) // ±10% live

	// set up pipeline
	pipeline := chrys.NewPipeline().AddStage(func(now time.Time) error {
		frames, err := client.GetFrames(series, now, 20)
		if err != nil {
			return err
		}

		zScore := algo.ZScore(algo.Closes(frames))
		fmt.Println("BB =", zScore)

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

For more complex use cases, you'll potentially want to split the signal and order logic into separate stages. You can pass data down through the stages chain like so:

```go
pipeline := chrys.NewPipeline()
pipeline.AddStage(func(now time.Time) error {
	frames, err := client.GetFrames(series, now, 20)
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
		orderConfig.Side = "buy"
	} else if pipeline.Get("bb") > 2 && pipeline.Get("bb-1/2") > 2 {
		orderConfig.Side = "sell"
	}

	// ...
})
```

Stages are processed sequentially in the order that they are defined.
