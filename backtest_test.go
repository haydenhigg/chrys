package chrys

import (
	"math"
	"testing"
)

// func Test_TotalReturn(t *testing.T) {
// 	// mock Backtest
// 	test := &Backtest{Values: []float64{1, 1.5, 2.25, 3.125, 4}}

// 	// TotalReturn()
// 	totalReturn := test.TotalReturn()

// 	// assert
// 	if math.Abs(totalReturn-3.0) > 10e-6 {
// 		t.Errorf("total return != 3.0: %f", totalReturn)
// 	}
// }

func Test_geometricMean(t *testing.T) {
	// mock
	xs := []float64{3, 6, 8, 2}

	// geometricMean()
	mean := geometricMean(xs)

	// assert
	if math.Abs(mean-4.1195343) > 10e-6 {
		t.Errorf("geometric mean != 4.1195343: %f", mean)
	}
}

func Test_AverageReturn(t *testing.T) {
	// mock Backtest
	test := &Backtest{Returns: []float64{.05, -.01, .01, .02, -.04}}

	// AverageReturn()
	averageReturn := test.AverageReturn()

	// assert
	if math.Abs(averageReturn-0.0055495) > 10e-6 {
		t.Errorf("average return != 0.0055495: %f", averageReturn)
	}
}
