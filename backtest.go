package chrys

import "time"

type BacktestReport struct {
	Start         time.Time
	End           time.Time
	StartValue    float64
	EndValue      float64
	TotalReturn   float64
	AverageReturn float64
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
