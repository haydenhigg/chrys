package chrys

import "slices"

type RankerRow struct {
	Key     string
	Factors []float64
}

func NewRankerRow(key string, factors ...float64) *RankerRow {
	return &RankerRow{key, factors}
}

type Ranker []*RankerRow

func NewRanker(capacity int) Ranker {
	return make(Ranker, 0, capacity)
}

func (ranker Ranker) Rank() Ranker {
	minNumFactors := 0
	for _, row := range ranker {
		numFactors := len(row.Factors)
		if numFactors > minNumFactors {
			minNumFactors = numFactors
		}
	}

	if minNumFactors == 0 {
		return ranker
	}

	mins := make([]float64, minNumFactors)
	maxes := make([]float64, minNumFactors)
	isInitialized := false

	for _, row := range ranker {
		for j, factor := range row.Factors {
			if factor < mins[j] || !isInitialized {
				mins[j] = factor
			}
			if factor > maxes[j] || !isInitialized {
				maxes[j] = factor
			}
		}

		isInitialized = true
	}

	for i, row := range ranker {
		for j, factor := range row.Factors {
			ranker[i].Factors[j] = (factor - mins[j]) / (maxes[j] - mins[j])
		}
	}

	slices.SortFunc(ranker, func(a, b *RankerRow) int {
		aScore, bScore := 0., 0.
		for i, aFactor := range a.Factors {
			aScore += aFactor
			bScore += b.Factors[i]
		}

		if aScore < bScore {
			return 1
		} else if aScore > bScore {
			return -1
		} else {
			return 0
		}
	})

	return ranker
}

func (ranker Ranker) Top(quantile float64) Ranker {
	return ranker[:int(quantile*float64(len(ranker)))]
}

func (ranker Ranker) Bottom(quantile float64) Ranker {
	return ranker[len(ranker)-int(quantile*float64(len(ranker))):]
}
