package algo

import (
	"math"
	"testing"
)

// tests
func Test_ZScore(t *testing.T) {
	// ZScore
	zScore := ZScore([]float64{1, 3, 4, 7, 5})

	// assert
	if math.Abs(zScore-0.447214) > 1e-6 {
		t.Errorf("zScore != 0.447214: %f", zScore)
	}
}
