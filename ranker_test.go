package chrys

import "testing"

func Test_Score(t *testing.T) {
	// create Ranker
	ranker := Ranker{
		NewRankerRow("BTC", 1, 10, 3), // => -0.199085
		NewRankerRow("ETH", 3, 8, 4),  // => 0.013022
		NewRankerRow("SOL", -2, 9, 4), // => -0.529420
		NewRankerRow("BCH", 0, 11, 5), // => 0.715483
	}

	// Score()
	scores := ranker.Score()

	// assert
	for i, row := range ranker {
		var (
			expectedKey     string
			expectedFactors []float64
			expectedScore   float64
		)

		switch i {
		case 3:
			expectedKey = "SOL"
			expectedFactors = []float64{-1.200961, -0.387298, 0}
			expectedScore = -0.529420
		case 2:
			expectedKey = "BTC"
			expectedFactors = []float64{0.2401922, 0.387298, -1.224745}
			expectedScore = -0.199085
		case 1:
			expectedKey = "ETH"
			expectedFactors = []float64{1.200961, -1.161895, 0}
			expectedScore = 0.013022
		case 0:
			expectedKey = "BCH"
			expectedFactors = []float64{-0.2401922, 1.161895, 1.224745}
			expectedScore = 0.715483
		}

		if row.Key != expectedKey {
			t.Errorf("ranker[%d].Key != %s: %s", i, row.Key, expectedKey)
		}

		rowScore := scores[row.Key]
		if !almostEqual(rowScore, expectedScore) {
			t.Errorf("scores[%s] != %f: %f", row.Key, expectedScore, rowScore)
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
