package chrys

import (
	"github.com/haydenhigg/chrys/algo"
	"slices"
)

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

func (ranker Ranker) Score() map[string]float64 {
	minNumFactors := 0
	for _, row := range ranker {
		numFactors := len(row.Factors)
		if numFactors > minNumFactors {
			minNumFactors = numFactors
		}
	}

	scores := make(map[string]float64, len(ranker))
	for _, row := range ranker {
		scores[row.Key] = 0
	}

	if minNumFactors == 0 {
		return scores
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

		scores[row.Key] = algo.Mean(ranker[i].Factors)
	}

	slices.SortFunc(ranker, func(a, b *RankerRow) int {
		if scores[a.Key] > scores[b.Key] {
			return -1
		} else if scores[a.Key] < scores[b.Key] {
			return 1
		} else {
			return 0
		}
	})

	return scores
}

func (ranker Ranker) Top(quantile float64) Ranker {
	return ranker[:int(quantile*float64(len(ranker)))]
}

func (ranker Ranker) Bottom(quantile float64) Ranker {
	return ranker[len(ranker)-int(quantile*float64(len(ranker))):]
}
