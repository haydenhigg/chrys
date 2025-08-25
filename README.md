# chrys
lightweight algorithmic trading framework

## to-do
0.5 rename to `chrys` and `Candle` to `Frame` and restructure, move non-caching client logic out and make as many things as possible top-level
1. change GetCandles to use config like Order
2. backtesting
3. state management through `algo` subpackage for stateful algos
4. add built-in logging to client
5. expand MLP implementation

## usage
This trades on **BOLL(20, 2)** signals for **1h BTC/USD** using a **10%** fractional trade amount.

```go
package main

import (
	"github.com/haydenhigg/chrys/algo"
	"github.com/haydenhigg/chrys/client"
	"github.com/haydenhigg/chrys/pipeline"
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

	// set up pipeline
	p := pipeline.New().AddStage(func(now time.Time) error {
		candles, err := c.GetCandles("BTC/USD", time.Hour, now, 20)
		if err != nil {
			return err
		}

		zScore := algo.ZScore(algo.Closes(candles))

		order := &client.OrderConfig{
			Pair:            "BTC/USD",
			BaseBalanceKey:  "XBT.F",
			QuoteBalanceKey: "ZUSD",
			Percent:         0.1,
		}

		if zScore < -2 {
			order.Side = client.MARKET_BUY
		} else if zScore > 2 {
			order.Side = client.MARKET_SELL
		} else {
			return nil
		}

		if err := c.PlaceOrder(order, now); err != nil {
			return err
		}

		return nil
	})

	// run
	if err := p.Run(time.Now()); err != nil {
		panic(err)
	}
}
```

For more complex use cases, you'll potentially want to split the signal and order logic into separate stages. You can pass data down through the stages chain like so:

```go
p := pipeline.New()
p.AddStage(func(now time.Time) error {
	candles, err := c.GetCandles("BTC/USD", time.Hour, now, 20)
	if err != nil {
		return err
	}

	closes := algo.Closes(candles)

	p.Set("bb", algo.ZScore(closes))
	p.Set("bb-1/2", algo.ZScore(closes[10:]))

	return nil
})
p.AddStage(func(now time.Time) error {
	fmt.Println(p.Data) // map[bb:... bb-1/2:...]

	// ...

	if p.Get("bb") < -2 && p.Get("bb-1/2") < -2 {
		order.Side = "buy"
	} else if p.Get("bb") > 2 && p.Get("bb-1/2") > 2 {
		order.Side = "sell"
	}

	// ...
})
```

Stages are processed sequentially in the order that they are defined.
