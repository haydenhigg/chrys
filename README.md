# clover
lightweight algorithmic trading framework

## to-do
1. add `time` parameter to engine.Run
2. pass `time` as parameter to signaler funcs
3. fetch candles using `time`
4. order placement from CachedClient
5. state management for stateful algos
6. backtesting
7. add `precedence` to signalers so there can be arbitrarily many steps that depend on previous steps (right now, the only steps are Signalers -> Handlers)
8. expand MLP implementation

## usage
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
	e := engine.New()

	e.Signal("zScore", func() (float64, error) {
		candles, err := c.Fetch("BTC/USD", time.Hour, 20)
		if err != nil {
			return 0, err
		}

		return algo.ZScore(algo.Closes(candles)), nil
	})

	e.Handle(func(signals engine.Signals) error {
		if signals["zScore"] < -2 {
			fmt.Println("buy!")
		} else if signals["zScore"] > 2 {
			fmt.Println("sell!")
		}
	})

	// run
	if err := e.Run(); err != nil {
		panic(err)
	}
}
````
