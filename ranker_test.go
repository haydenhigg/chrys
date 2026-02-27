package chrys

import "testing"

func Test_Rank(t *testing.T) {
	// create Ranker
	ranker := Ranker{
		NewRankerRow("BTC", 1, 10, 3), // => [0.6, 0.666_, 0.0] => 1.266_
		NewRankerRow("ETH", 3, 8, 4),  // => [1.0, 0.0, 0.5] => 1.5
		NewRankerRow("SOL", -2, 9, 4), // => [0.0, 0.333_, 0.5] => 0.833_
		NewRankerRow("BCH", 0, 11, 5), // => [0.4, 1.0, 1.0] => 2.4
	}

	// Rank()
	ranker.Rank()

	// assert
	for i, row := range ranker {
		var (
			expectedKey     string
			expectedFactors []float64
		)

		switch i {
		case 0:
			expectedKey = "BCH"
			expectedFactors = []float64{0.4, 1, 1}
		case 1:
			expectedKey = "ETH"
			expectedFactors = []float64{1.0, 0, 0.5}
		case 2:
			expectedKey = "BTC"
			expectedFactors = []float64{0.6, 0.6666667, 0}
		case 3:
			expectedKey = "SOL"
			expectedFactors = []float64{0, 0.3333333, 0.5}
		}

		if row.Key != expectedKey {
			t.Errorf("ranker[%d].Key != %s: %s", i, row.Key, expectedKey)
		}

		assertSlicesEqual(row.Factors, expectedFactors, t)
	}
}

func Test_Top(t *testing.T) {
	// create Ranker
	ranker := Ranker{
		NewRankerRow("BCH", 0.4, 1, 1),
		NewRankerRow("ETH", 1.0, 0, 0.5),
		NewRankerRow("BTC", 0.6, 0.6666667, 0),
		NewRankerRow("SOL", 0, 0.3333333, 0.5),
	}

	// Top()
	top := ranker.Top(0.25)

	// assert
	if len(top) != 1 {
		t.Errorf("len(top) != 1: %d", len(top))
	}

	if top[0].Key != "BCH" {
		t.Errorf("top[0].Key != BCH: %s", top[0].Key)
	}

	assertSlicesEqual(top[0].Factors, []float64{0.4, 1, 1}, t)
}

func Test_Bottom(t *testing.T) {
	// create Ranker
	ranker := Ranker{
		NewRankerRow("BCH", 0.4, 1, 1),
		NewRankerRow("ETH", 1.0, 0, 0.5),
		NewRankerRow("BTC", 0.6, 0.6666667, 0),
		NewRankerRow("SOL", 0, 0.3333333, 0.5),
	}

	// Bottom()
	bottom := ranker.Bottom(0.25)

	// assert
	if len(bottom) != 1 {
		t.Errorf("len(bottom) != 1: %d", len(bottom))
	}

	if bottom[0].Key != "SOL" {
		t.Errorf("bottom[0].Key != SOL: %s", bottom[0].Key)
	}

	assertSlicesEqual(bottom[0].Factors, []float64{0, 0.333333, 0.5}, t)
}
