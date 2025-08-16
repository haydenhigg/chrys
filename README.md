# clover
lightweight algorithmic trading framework

## goals
- small codebase and simple design
- works for ML4T pipelines

## to-do
1. add `time` parameter to engine.Run
2. pass `time` as parameter to signaler funcs
3. fetch candles using `time`
4. order placement from CachedClient
5. state management
6. expand MLP implementation

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
