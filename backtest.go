package chrys

import (
	"github.com/haydenhigg/chrys/algo"
	"math"
	"time"
)

type Backtest struct {
	Step    time.Duration
	N       int
	Values  []float64
	Returns []float64

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
	backtest := &Backtest{
		Values:  []float64{},
		Returns: []float64{},
	}

	return backtest.SetStep(step)
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
	// append value
	backtest.Values = append(backtest.Values, value)

	if backtest.N > 0 {
		// append return
		r := value/backtest.Values[backtest.N-1] - 1
		backtest.Returns = append(backtest.Returns, r)

		// update meanReturn
		n := float64(backtest.N)
		backtest.meanReturn = (backtest.meanReturn*(n-1) + r) / n
	}

	// increment n
	backtest.N++
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
	if backtest.N <= 1 {
		return 0
	}

	growthFactor := backtest.Values[backtest.N-1] / backtest.Values[0]
	annualizationPower := YEAR / (float64(backtest.N-1) * backtest.step)

	return math.Pow(growthFactor, annualizationPower) - 1
}

func (backtest *Backtest) Volatility() float64 {
	if len(backtest.Returns) <= 1 {
		return 0
	}

	vol := algo.StandardDeviation(backtest.Returns, backtest.meanReturn)
	annualizationCoef := math.Sqrt(YEAR / backtest.step)
	// Sharpe and Sortino are also annualized this way. In theory, this is not
	// a robust way to annualize unless returns are i.i.d., which is likely
	// untrue. However, in practice, this tends to produce stable results
	// across time scales more reliably than an estimated Hurst exponent.

	return vol * annualizationCoef
}

func (backtest *Backtest) Sharpe(minReturn float64) float64 {
	vol := algo.StandardDeviation(backtest.Returns, backtest.meanReturn)
	if vol == 0 {
		return 0
	}

	periodsPerYear := YEAR / backtest.step
	periodicMinReturn := math.Pow(1+minReturn, 1/periodsPerYear) - 1

	sharpe := (backtest.meanReturn - periodicMinReturn) / vol
	annualizationCoef := math.Sqrt(periodsPerYear)

	return sharpe * annualizationCoef
}

func (backtest *Backtest) Sortino(minReturn float64) float64 {
	periodsPerYear := YEAR / backtest.step
	periodicMinReturn := math.Pow(1+minReturn, 1/periodsPerYear) - 1

	downside := make([]float64, len(backtest.Returns))
	for i, r := range backtest.Returns {
		downside[i] = math.Pow(min(0, r-periodicMinReturn), 2)
	}

	downsideVol := math.Sqrt(algo.Mean(downside))
	if downsideVol == 0 {
		return 0
	}

	sortino := (backtest.meanReturn - periodicMinReturn) / downsideVol
	annualizationCoef := math.Sqrt(periodsPerYear)

	return sortino * annualizationCoef
}

func (backtest *Backtest) Omega(minReturn float64) float64 {
	periodicMinReturn := math.Pow(1+minReturn, backtest.step/YEAR) - 1

	var sumGain, sumLoss float64
	for _, r := range backtest.Returns {
		if r > periodicMinReturn {
			sumGain += r - periodicMinReturn
		} else if r < periodicMinReturn {
			sumLoss += periodicMinReturn - r
		}
	}

	return sumGain / sumLoss
}

func (backtest *Backtest) Martin(minReturn float64) float64 {
	var peak, sum float64
	for _, v := range backtest.Values {
		if v > peak {
			peak = v
		} else {
			sum += math.Pow((peak-v)/peak, 2)
		}
	}

	ulcerIndex := math.Sqrt(sum / float64(backtest.N))

	return (backtest.Return() - minReturn) / ulcerIndex
}
