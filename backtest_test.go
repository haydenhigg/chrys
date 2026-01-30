package chrys

import (
	"fmt"
	"math"
	"testing"
	"time"
)

// helpers
func assertMatricesEqual(a, b [][]float64, t *testing.T) {
	for i := range a {
		if i >= len(b) {
			t.Errorf("b[%d] does not exist", i)
		} else {
			for j, va := range a[i] {
				if j >= len(b[i]) {
					t.Errorf("b[%d][%d] does not exist", i, j)
				} else if vb := b[i][j]; math.Abs(va-vb) > 1e-6 {
					t.Errorf(
						"a[%d][%d] != b[%d][%d]: %v != %v",
						i, j, i, j, va, vb,
					)
				}
			}
		}
	}

	for i := range b {
		if i >= len(a) {
			t.Errorf("a[%d] does not exist", i)
		}
	}
}

func assertSlicesEqual(a, b []float64, t *testing.T) {
	for i, va := range a {
		if i >= len(b) {
			t.Errorf("b[%d] does not exist", i)
		} else if vb := b[i]; math.Abs(va-vb) > 1e-6 {
			t.Errorf("a[%d] != b[%d]: %v != %v", i, i, va, vb)
		}
	}
}

// tests
func Test_Record(t *testing.T) {
	// create Backtest
	backtest := NewBacktest(time.Hour)

	// Record()
	backtest.Record(10, 100)
	backtest.Record(11, 110)
	backtest.Record(12, 121)
	backtest.Record(11.4, 115)
	backtest.Record(11.97, 121.9)

	// assert
	assertMatricesEqual(backtest.Values, [][]float64{
		{10, 100},
		{11, 110},
		{12, 121},
		{11.4, 115},
		{11.97, 121.9},
	}, t)

	assertMatricesEqual(backtest.Returns, [][]float64{
		{.1, .1},
		{.0909091, .1},
		{-.05, -.04958677685},
		{.05, .06},
	}, t)
}

func Test_Return(t *testing.T) {
	// create Backtest
	backtest := NewBacktest(time.Hour * 24)

	// Record()
	backtest.Record(10, 50)
	backtest.Record(9, 45)
	backtest.Record(10.5, 52.5)
	backtest.Record(11, 55)
	backtest.Record(10.2, 51.1)

	// Return()
	returns := backtest.Return()

	// assert
	assertSlicesEqual(returns, []float64{3.2443632, 3.8968341}, t)
}

func Test_Volatility(t *testing.T) {
	// create Backtest
	backtest := NewBacktest(time.Hour * 24)

	// Record()
	backtest.Record(10, 100)
	backtest.Record(11, 110)
	backtest.Record(12, 121)
	backtest.Record(11.4, 115)
	backtest.Record(11.97, 121.9)

	// Volatility()
	vols := backtest.Volatility()

	// assert
	assertSlicesEqual(vols, []float64{1.3122253, 1.350494}, t)
}

func Test_SharpeRatio(t *testing.T) {
	// create Backtest
	backtest := NewBacktest(time.Hour * 24)

	// Record()
	backtest.Record(10, 100)
	backtest.Record(11, 110)
	backtest.Record(12, 121)
	backtest.Record(11.4, 115)
	backtest.Record(11.97, 121.9)

	fmt.Println(backtest.returnsColumn(1))

	// SharpeRatio()
	sharpes := backtest.SharpeRatio(0.04)

	// assert
	assertSlicesEqual(sharpes, []float64{13.2450234, 14.1875538}, t)
}
