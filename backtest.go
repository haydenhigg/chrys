package chrys

import (
	"github.com/haydenhigg/chrys/algo"
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

const YEAR float64 = 3.1536e+16

func (test *Backtest) CAGR() float64 {
	totalReturn := test.Values[len(test.Values)-1] / test.Values[0]
	duration := test.End.Sub(test.Start)

	return math.Pow(totalReturn, YEAR/float64(duration)) - 1
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
	growthFactors := make([]float64, len(test.Returns))
	for i, ret := range test.Returns {
		growthFactors[i] = 1 + ret
	}

	return geometricMean(growthFactors) - 1
}

func (test *Backtest) Volatility() float64 {
	volatility := algo.StandardDeviation(test.Returns, algo.Mean(test.Returns))

	return volatility * math.Sqrt(YEAR/float64(test.Step))
}

func (test *Backtest) Sharpe(annualRiskFreeReturn float64) float64 {
	periodicRiskFreeReturn := annualRiskFreeReturn / (YEAR / float64(test.Step))

	meanReturn := algo.Mean(test.Returns)
	volatility := algo.StandardDeviation(test.Returns, meanReturn)

	sharpe := (meanReturn - periodicRiskFreeReturn) / volatility

	return sharpe * math.Sqrt(YEAR/float64(test.Step))
}

func (test *Backtest) Sortino(annualRiskFreeReturn float64) float64 {
	periodicRiskFreeReturn := annualRiskFreeReturn / (YEAR / float64(test.Step))

	downside := make([]float64, len(test.Returns))
	for i, ret := range test.Returns {
		if ret < periodicRiskFreeReturn {
			downside[i] = ret - periodicRiskFreeReturn
		} else {
			downside[i] = 0
		}
	}

	meanReturn := algo.Mean(test.Returns)
	downsideVolatility := algo.StandardDeviation(downside, algo.Mean(downside))

	sortino := (meanReturn - periodicRiskFreeReturn) / downsideVolatility

	return sortino * math.Sqrt(YEAR/float64(test.Step))
}

func (test *Backtest) Skew() float64 {
	meanReturn := algo.Mean(test.Returns)
	volatility := algo.StandardDeviation(test.Returns, meanReturn)

	skew := 0.
	for _, ret := range test.Returns {
		skew += math.Pow((ret-meanReturn)/volatility, 3)
	}

	return skew / float64(len(test.Returns))
}

func (test *Backtest) MaxDrawdown() float64 {
	peak := test.Values[0]
	maxDrawdown := 0.

	var drawdown float64
	for _, value := range test.Values {
		if value > peak {
			peak = value
		} else if drawdown = value/peak - 1; drawdown < maxDrawdown {
			maxDrawdown = drawdown
		}
	}

	return maxDrawdown
}
