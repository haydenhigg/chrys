package chrys

import (
	"math"
	"time"
)

type Backtest struct {
	Start   time.Time
	End     time.Time
	Step    time.Duration
	Values  []float64
	Returns []float64
}

const YEAR = time.Hour * 8760

func (test *Backtest) CAGR() float64 {
	totalReturn := test.Values[len(test.Values)-1] / test.Values[0]
	duration := test.End.Sub(test.Start)

	return math.Pow(totalReturn, float64(YEAR)/float64(duration)) - 1
}

func geometricMean(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}

	product := 1.
	for _, x := range xs {
		product *= x
	}

	return math.Pow(product, 1/float64(len(xs)))
}

func (test *Backtest) AverageReturn() float64 {
	onePlusReturns := make([]float64, len(test.Returns))
	for i, r := range test.Returns {
		onePlusReturns[i] = 1 + r
	}

	return geometricMean(onePlusReturns) - 1
}

// // print metrics
// returns := make([]float64, len(values)-1)
// for i, value := range values[1:] {
// 	returns[i] = value/values[i] - 1
// }

// avgReturn := algo.Mean(returns)
// volatility := algo.StandardDeviation(returns, avgReturn)
// fmt.Printf("return: %+.2f%%\n", 100*(values[len(values)-1]/values[0]-1))
// fmt.Printf("annualized avg. return: %+.2f%%\n", 100*(math.Pow(1+avgReturn, 8760)-1))
// fmt.Printf("annualized volatility: %.2f%%\n", 100*volatility*math.Sqrt(8760))

// btcReturns := make([]float64, len(btcValues)-1)
// excessReturns := []float64{}
// for i, value := range btcValues[1:] {
// 	btcReturns[i] = value/btcValues[i] - 1
// 	excessReturns = append(excessReturns, returns[i]-btcReturns[i])
// }

// btcAvgReturn := algo.Mean(btcReturns)
// btcVolatility := algo.StandardDeviation(btcReturns, btcAvgReturn)
// fmt.Printf("baseline return: %+.2f%%\n", 100*(btcValues[len(btcValues)-1]/btcValues[0]-1))
// fmt.Printf("baseline annualized avg. return: %+.2f%%\n", 100*(math.Pow(1+btcAvgReturn, 8760)-1))
// fmt.Printf("baseline annualized volatility: %.2f%%\n", 100*btcVolatility*math.Sqrt(8760))

// avgExcessReturn := algo.Mean(excessReturns)
// excessVolatility := algo.StandardDeviation(excessReturns, avgExcessReturn)

// fmt.Printf("sharpe ratio: %.2f\n", ((avgReturn-btcAvgReturn)/excessVolatility)*math.Sqrt(8760))

// fmt.Println("buys:", buys)
// fmt.Println("sells:", sells)
