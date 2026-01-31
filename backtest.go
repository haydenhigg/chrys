package chrys

import (
	"github.com/haydenhigg/chrys/algo"
	"math"
	"time"
)

type Backtest struct {
	Step       time.Duration
	N          int
	FirstValue float64
	LastValue  float64
	Returns    []float64

	// input step as float64 for easier manipulation
	step float64

	// O(1) max drawdown calculation
	peakValue   float64
	maxDrawdown float64

	// O(1) arithmetic mean return calculation
	meanReturn float64
}

// initializer
func NewBacktest(step time.Duration) *Backtest {
	return (&Backtest{Returns: []float64{}}).SetStep(step)
}

// setters
func (backtest *Backtest) SetStep(step time.Duration) *Backtest {
	// having a zero step will cause problems for metrics
	if int64(step) > 0 {
		backtest.Step = step
	} else {
		backtest.Step = time.Duration(1)
	}

	backtest.step = float64(backtest.Step)

	return backtest
}

// Update
func (backtest *Backtest) update(value float64) {
	if backtest.N == 0 {
		backtest.FirstValue = value
	} else {
		// append return
		r := value/backtest.LastValue - 1
		backtest.Returns = append(backtest.Returns, r)

		// update meanReturn
		n := float64(backtest.N)
		backtest.meanReturn = (backtest.meanReturn*(n-1) + r) / n
	}

	backtest.N++
	backtest.LastValue = value
}

func (backtest *Backtest) updateDrawdown(value float64) {
	if value > backtest.peakValue {
		backtest.peakValue = value
		return
	}

	drawdown := value/backtest.peakValue - 1
	if drawdown < backtest.maxDrawdown {
		backtest.maxDrawdown = drawdown
	}
}

func (backtest *Backtest) Update(value float64) *Backtest {
	backtest.update(value)
	backtest.updateDrawdown(value)

	return backtest
}

// metrics
const YEAR float64 = 3.1536e+16

func (backtest *Backtest) MaxDrawdown() float64 {
	return backtest.maxDrawdown
}

func (backtest *Backtest) Return() float64 {
	if backtest.FirstValue == 0 {
		return 0
	}

	growthFactor := backtest.LastValue / backtest.FirstValue
	annualizationPower := YEAR / (float64(backtest.N) * backtest.step)

	return math.Pow(growthFactor, annualizationPower) - 1
}

func (backtest *Backtest) Volatility() float64 {
	if len(backtest.Returns) <= 1 {
		return 0
	}

	vol := algo.StandardDeviation(backtest.Returns, backtest.meanReturn)
	annualizationCoef := math.Sqrt(YEAR / backtest.step)

	return vol * annualizationCoef
}

func (backtest *Backtest) Sharpe(riskFreeReturn float64) float64 {
	vol := algo.StandardDeviation(backtest.Returns, backtest.meanReturn)
	if vol == 0 {
		return 0
	}

	periodsPerYear := YEAR / backtest.step
	periodicRiskFreeReturn := riskFreeReturn / periodsPerYear

	sharpe := (backtest.meanReturn - periodicRiskFreeReturn) / vol
	annualizationCoef := math.Sqrt(periodsPerYear)

	return sharpe * annualizationCoef
}

func (backtest *Backtest) GainLoss() float64 {
	var (
		sumGain, sumLoss float64
		nGain, nLoss     int
	)

	for _, r := range backtest.Returns {
		if r > 0 {
			sumGain += r
			nGain++
		} else if r < 0 {
			sumLoss -= r
			nLoss++
		}
	}

	meanGain := sumGain / float64(nGain)
	meanLoss := sumLoss / float64(nLoss)

	return (meanGain - meanLoss) / (meanGain + meanLoss)
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
