package algo

import (
	"math"
	"testing"
)

func Test_Covariance(t *testing.T) {
	// create xs and ys
	xs := []float64{1, 2, 3}
	ys := []float64{10, 9, 12}

	// Covariance()
	cov := Covariance(xs, ys, Mean(xs), Mean(ys))

	// assert
	if cov != 1.0 {
		t.Errorf("covariance != 1.0: %f", cov)
	}
}

func Test_Correlation(t *testing.T) {
	// create xs and ys
	xs := []float64{1, 2, 3}
	ys := []float64{10, 9, 12}

	// Correlation()
	corr := Correlation(xs, ys, Mean(xs), Mean(ys))

	// assert
	if math.Abs(corr-0.6546537) > 1e-6 {
		t.Errorf("correlation != 0.6546537: %f", corr)
	}
}
