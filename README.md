# clover
lightweight algorithmic trading framework

## to-do
1. backtesting
2. state management through `algo` subpackage for stateful algos
3. expand MLP implementation

## usage
This trades on **BOLL(20, 2)** signals for **1h BTC/USD** using a **10%** fractional trade amount.

```go
package main

import (
	"github.com/haydenhigg/clover/algo"
	"github.com/haydenhigg/clover/client"
	"github.com/haydenhigg/clover/engine"
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

	// set up engine
	e := engine.New().Handle(func(now time.Time) error {
		candles, err := c.GetCandles("BTC/USD", time.Hour, now, 20)
		if err != nil {
			return err
		}

		zScore := algo.ZScore(algo.Closes(candles))

		order := &client.OrderConfig{
			Pair:      "BTC/USD",
			BaseCode:  "XBT.F",
			QuoteCode: "ZUSD",
			Percent:   0.1,
		}

		if zScore < -2 {
			order.Side = "buy"
		} else if zScore > 2 {
			order.Side = "sell"
		}

		if err := c.MarketOrder(order); err != nil {
			return err
		}

		return nil
	})

	// run
	if err := e.Run(time.Now()); err != nil {
		panic(err)
	}
}
```

For more complex use cases, you'll likely want to split the signal and order logic into separate handlers. You can pass data down through the handler chain like so:

```go
// ...
e := engine.New()
e.Handle(func(now time.Time) error {
	candles, err := c.GetCandles("BTC/USD", time.Hour, now, 20)
	if err != nil {
		return err
	}

	closes := algo.Closes(candles)

	e.Set("bb", algo.ZScore(closes))
	e.Set("bb-1/2", algo.ZScore(closes[10:]))

	return nil
})
e.Handle(func(now time.Time) error {
	fmt.Println(e.Data) // map[bb:... bb-1/2:...]

	order := &client.OrderConfig{
		Pair:      "BTC/USD",
		BaseCode:  "XBT.F",
		QuoteCode: "ZUSD",
		Percent:   0.1,
	}

	if e.Get("bb") < -2 && e.Get("bb-1/2") < -2 {
		order.Side = "buy"
	} else if e.Get("bb") > 2 && e.Get("bb-1/2") > 2 {
		order.Side = "sell"
	}

	if err := c.MarketOrder(order); err != nil {
		return err
	}

	return nil
})
// ...
```

Handlers are processed sequentially in the order that they are defined.
