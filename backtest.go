package chrys

import (
	"github.com/haydenhigg/chrys/algo"
	"math"
	"time"
)

type Backtest struct {
	Step    time.Duration
	Values  [][]float64
	Returns [][]float64
}

func NewBacktest(step time.Duration) *Backtest {
	return &Backtest{
		Step:    step,
		Values:  [][]float64{},
		Returns: [][]float64{},
	}
}

func (backtest *Backtest) Record(values ...float64) *Backtest {
	backtest.Values = append(backtest.Values, values)

	if len(backtest.Values) > 1 {
		oldValues := backtest.Values[len(backtest.Values)-2]

		m := min(len(oldValues), len(values))
		returns := make([]float64, m)

		for i := range m {
			returns[i] = values[i]/oldValues[i] - 1
		}

		backtest.Returns = append(backtest.Returns, returns)
	}

	return backtest
}

const YEAR float64 = 3.1536e+16

func (backtest *Backtest) Return() []float64 {
	n := len(backtest.Values)
	if n <= 1 {
		return []float64{}
	}

	duration := float64(n) * float64(backtest.Step)
	annualizationPower := YEAR / duration

	m := len(backtest.Values[0])
	returns := make([]float64, m)
	for j := range m {
		returns[j] = math.Pow(
			backtest.Values[n-1][j]/backtest.Values[0][j],
			annualizationPower,
		) - 1
	}

	return returns
}

// func (backtest *Backtest) MaxDrawdown() []float64 {
// 	peak := test.Values[0]
// 	maxDrawdown := 0.

// 	var drawdown float64
// 	for _, value := range test.Values {
// 		if value > peak {
// 			peak = value
// 		} else if drawdown = value/peak - 1; drawdown < maxDrawdown {
// 			maxDrawdown = drawdown
// 		}
// 	}

// 	return maxDrawdown
// }

func (backtest *Backtest) returnsColumn(j int) []float64 {
	n := len(backtest.Returns)
	series := make([]float64, n)
	for i, returnsRow := range backtest.Returns {
		if j >= len(returnsRow) {
			break
		}

		series[i] = returnsRow[j]
	}

	return series
}

func (backtest *Backtest) Volatility() []float64 {
	n := len(backtest.Returns)
	if n <= 1 {
		return []float64{}
	}

	m := len(backtest.Values[0])
	vols := make([]float64, m)

	annualization := math.Sqrt(YEAR / float64(backtest.Step))

	for j := range m {
		series := backtest.returnsColumn(j)

		mean := algo.Mean(series)
		vols[j] = algo.StandardDeviation(series, mean) * annualization
	}

	return vols
}

func (backtest *Backtest) SharpeRatio(riskFreeReturn float64) []float64 {
	n := len(backtest.Returns)
	if n <= 1 {
		return []float64{}
	}

	m := len(backtest.Values[0])
	sharpes := make([]float64, m)

	periodsPerYear := YEAR / float64(backtest.Step)
	periodicRiskFreeReturn := riskFreeReturn / periodsPerYear
	annualization := math.Sqrt(periodsPerYear)

	for j := range m {
		series := backtest.returnsColumn(j)

		mean := algo.Mean(series)
		vol := algo.StandardDeviation(series, mean)

		if vol == 0 {
			sharpes[j] = 0
			continue
		}

		sharpes[j] = ((mean - periodicRiskFreeReturn) / vol) * annualization
	}

	return sharpes
}

// func (test *Backtest) Sortino(annualRiskFreeReturn float64) float64 {
// 	periodicRiskFreeReturn := annualRiskFreeReturn / (YEAR / float64(test.Step))

// 	downside := make([]float64, len(test.Returns))
// 	for i, ret := range test.Returns {
// 		if ret < periodicRiskFreeReturn {
// 			downside[i] = ret - periodicRiskFreeReturn
// 		} else {
// 			downside[i] = 0
// 		}
// 	}

// 	meanReturn := algo.Mean(test.Returns)
// 	downsideVolatility := algo.StandardDeviation(downside, algo.Mean(downside))

// 	sortino := (meanReturn - periodicRiskFreeReturn) / downsideVolatility

// 	return sortino * math.Sqrt(YEAR/float64(test.Step))
// }

// func (test *Backtest) Skew() float64 {
// 	meanReturn := algo.Mean(test.Returns)
// 	volatility := algo.StandardDeviation(test.Returns, meanReturn)

// 	skew := 0.
// 	for _, ret := range test.Returns {
// 		skew += math.Pow((ret-meanReturn)/volatility, 3)
// 	}

// 	return skew / float64(len(test.Returns))
// }
