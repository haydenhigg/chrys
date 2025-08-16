# clover
lightweight algorithmic trading framework

## usage
```go
package main

import (
	"github.com/haydenhigg/clover/algo"
	"github.com/haydenhigg/clover/client"
	"github.com/haydenhigg/clover/engine"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

func readConfig(fileName string) (map[string]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config map[string]string
	if err = json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	// set up client
	config, err := readConfig("config.json")
	if err != nil {
		panic(err)
	}

	c, err := client.NewKraken(config["API_KEY"], config["API_SECRET"])
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
	}).Handle(func(signals engine.Signals) error {
		if signals["zScore"] < -2 {
			fmt.Println("buy!")
		} else if signals["zScore"] > 2 {
			fmt.Println("sell!")
		}
	})

	// run
	if errs := e.Run(); len(errs) > 0 {
		fmt.Println(errs)
	}
}
````
