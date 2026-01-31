package chrys

import (
	// "fmt"
	"math"
	"testing"
	"time"
)

// helpers
func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= 1e-6
}

func assertSlicesEqual(a, b []float64, t *testing.T) {
	for i, va := range a {
		if i >= len(b) {
			t.Errorf("b[%d] does not exist", i)
		} else if vb := b[i]; !almostEqual(va, vb) {
			t.Errorf("a[%d] != b[%d]: %v != %v", i, i, va, vb)
		}
	}
}

// tests
// tests -> setters
func Test_SetStep(t *testing.T) {
	backtest := &Backtest{}
	backtest.SetStep(time.Minute * 1337)

	if backtest.Step != time.Minute*1337 {
		t.Errorf("Step != 1337min: %v", backtest.Step)
	}
}

func Test_SetStepZero(t *testing.T) {
	backtest := &Backtest{}
	backtest.SetStep(time.Time{}.Sub(time.Time{}))

	if backtest.Step != time.Nanosecond {
		t.Errorf("Step != 1: %v", backtest.Step)
	}
}

// tests -> Update
func Test_UpdateFirst(t *testing.T) {
	// create Backtest
	backtest := NewBacktest(time.Hour)

	// Update()
	backtest.Update(10)

	// assert
	if backtest.N != 1 {
		t.Errorf("N != 1: %d", backtest.N)
	}
	if backtest.FirstValue != 10 {
		t.Errorf("FirstValue != 10: %f", backtest.FirstValue)
	}
	if backtest.LastValue != 10 {
		t.Errorf("LastValue != 10: %f", backtest.LastValue)
	}
	assertSlicesEqual(backtest.Returns, []float64{}, t)
	if backtest.peakValue != 10 {
		t.Errorf("peakValue != 10: %f", backtest.peakValue)
	}
	if backtest.maxDrawdown != 0 {
		t.Errorf("maxDrawdown != 0: %f", backtest.maxDrawdown)
	}
	if backtest.meanReturn != 0 {
		t.Errorf("meanReturn != 0: %f", backtest.meanReturn)
	}
}

func Test_UpdateSubsequent(t *testing.T) {
	// create Backtest
	backtest := NewBacktest(time.Hour)

	// Update()
	backtest.Update(10)
	backtest.Update(13)

	// assert
	if backtest.N != 2 {
		t.Errorf("N != 2: %d", backtest.N)
	}
	if backtest.FirstValue != 10 {
		t.Errorf("FirstValue != 10: %f", backtest.FirstValue)
	}
	if backtest.LastValue != 13 {
		t.Errorf("LastValue != 13: %f", backtest.LastValue)
	}
	assertSlicesEqual(backtest.Returns, []float64{.3}, t)
	if backtest.peakValue != 13 {
		t.Errorf("peakValue != 13: %f", backtest.peakValue)
	}
	if backtest.maxDrawdown != 0 {
		t.Errorf("maxDrawdown != 0: %f", backtest.maxDrawdown)
	}
	if !almostEqual(backtest.meanReturn, 0.3) {
		t.Errorf("meanReturn != 0.3: %f", backtest.meanReturn)
	}

	// Update() again
	backtest.Update(11)

	// assert
	if backtest.N != 3 {
		t.Errorf("N != 3: %d", backtest.N)
	}
	if backtest.FirstValue != 10 {
		t.Errorf("FirstValue != 10: %f", backtest.FirstValue)
	}
	if backtest.LastValue != 11 {
		t.Errorf("LastValue != 11: %f", backtest.LastValue)
	}
	assertSlicesEqual(backtest.Returns, []float64{.3, -.1538462}, t)
	if backtest.peakValue != 13 {
		t.Errorf("peakValue != 13: %f", backtest.peakValue)
	}
	if !almostEqual(backtest.maxDrawdown, -.1538462) {
		t.Errorf("maxDrawdown != -.1538462: %f", backtest.maxDrawdown)
	}
	if !almostEqual(backtest.meanReturn, .0730769) {
		t.Errorf("meanReturn != 0.0730769: %f", backtest.meanReturn)
	}
}

// tests -> metrics
func Test_MaxDrawdown(t *testing.T) {
	// create Backtest
	backtest := NewBacktest(time.Hour)

	// Update()
	backtest.Update(100)
	backtest.Update(120)
	backtest.Update(105)
	backtest.Update(110)

	// assert
	if !almostEqual(backtest.MaxDrawdown(), -.125) {
		t.Errorf("MaxDrawdown() != -.125: %f", backtest.MaxDrawdown())
	}
}

func Test_Return(t *testing.T) {
	// create Backtest
	backtest := NewBacktest(24 * time.Hour)

	// Update()
	backtest.Update(100)
	backtest.Update(120)
	backtest.Update(105)
	backtest.Update(110)
	backtest.Update(101)

	// assert
	if !almostEqual(backtest.Return(), 1.0675703) {
		t.Errorf("Return() != 1.0675703: %f", backtest.Return())
	}
}

func Test_Volatility(t *testing.T) {
	// create Backtest
	backtest := NewBacktest(24 * time.Hour)

	// Update()
	backtest.Update(100)
	backtest.Update(120)
	backtest.Update(105)
	backtest.Update(110)
	backtest.Update(101)

	// assert
	if !almostEqual(backtest.Volatility(), 2.7941769) {
		t.Errorf("Volatility() != 2.7941769: %f", backtest.Volatility())
	}
}

func Test_Sharpe(t *testing.T) {
	// create Backtest
	backtest := NewBacktest(24 * time.Hour)

	// Update()
	backtest.Update(100)
	backtest.Update(120)
	backtest.Update(105)
	backtest.Update(110)
	backtest.Update(101)

	// assert
	if !almostEqual(backtest.Sharpe(.01), 1.3288633) {
		t.Errorf("Sharpe(0.01) != 1.3288633: %f", backtest.Sharpe(.01))
	}
}

func Test_Sortino(t *testing.T) {
	// create Backtest
	backtest := NewBacktest(24 * time.Hour)

	// Update()
	backtest.Update(100)
	backtest.Update(120) // +0.2
	backtest.Update(105) // -0.125
	backtest.Update(110) // +0.04761904762
	backtest.Update(101) // -0.08181818182

	// assert
	if !almostEqual(backtest.Sortino(0.01), 2.6011659) {
		t.Errorf("Sortino() != 2.6011659: %f", backtest.Sortino(0.01))
	}
}

func Test_Omega(t *testing.T) {
	// create Backtest
	backtest := NewBacktest(24 * time.Hour)

	// Update()
	backtest.Update(100)
	backtest.Update(120) // +0.2
	backtest.Update(105) // -0.125
	backtest.Update(110) // +0.04761904762
	backtest.Update(101) // -0.08181818182

	// assert
	if !almostEqual(backtest.Omega(0.01), 1.1966969) {
		t.Errorf("Omega() != 1.1966969: %f", backtest.Omega(0.01))
	}
}
