package algo

import "github.com/haydenhigg/clover/candle"

func Closes(candles []*candle.Candle) []float64 {
	closes := make([]float64, len(candles))
	for i, candle := range candles {
		closes[i] = candle.Close
	}

	return closes
}
