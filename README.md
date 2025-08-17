# clover
lightweight algorithmic trading framework

## to-do
1. order placement through Client
2. state management through `algo` subpackage for stateful algos
3. backtesting
4. expand MLP implementation

## usage
This generates **BOLL(20, 2)** signals for **1h BTC/USD**.

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

		if zScore < -2 {
			fmt.Println("buy")
		} else if zScore > 2 {
			fmt.Println("sell")
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

	if e.Get("bb") < -2 && e.Get("bb-1/2") < -2 {
		fmt.Println("buy")
	} else if e.Get("bb") > 2 && e.Get("bb-1/2") > 2 {
		fmt.Println("sell")
	}

	return nil
})
// ...
```

Handlers are processed sequentially in the order that they are defined.
