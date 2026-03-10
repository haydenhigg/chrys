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

func (ranker Ranker) transposeFactors() [][]float64 {
	factors := [][]float64{}

	for _, row := range ranker {
		for j, factor := range row.Factors {
			if j >= len(factors) {
				factors = append(factors, []float64{factor})
			} else {
				factors[j] = append(factors[j], factor)
			}
		}
	}

	return factors
}

func (ranker Ranker) Score() map[string]float64 {
	factors := ranker.transposeFactors()
	if len(factors) == 0 {
		return nil
	}

	means := make([]float64, len(factors))
	stddevs := make([]float64, len(factors))
	for j, factorValues := range factors {
		means[j] = algo.Mean(factorValues)
		stddevs[j] = algo.StandardDeviation(factorValues, means[j])
	}

	scores := make(map[string]float64, len(ranker))
	for i, row := range ranker {
		for j, factor := range row.Factors {
			if stddevs[j] > 1e-8 {
				ranker[i].Factors[j] = (factor - means[j]) / stddevs[j]
			} else {
				ranker[i].Factors[j] = 0.
			}
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
