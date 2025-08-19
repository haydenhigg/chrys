package client

import (
	"github.com/haydenhigg/clover/candle"
	"time"
)

type Store struct {
	Candles  map[string]map[time.Duration][]*candle.Candle
	Balances map[string]float64
}

func (s *Store) TryGetCandlesSince(
	pair string,
	interval time.Duration,
	since time.Time,
) ([]*candle.Candle, bool) {
	since = since.Truncate(interval).Add(-time.Second)

	if _, ok := s.Candles[pair]; !ok {
		s.Candles[pair] = map[time.Duration][]*candle.Candle{}
	}

	candles, ok := s.Candles[pair][interval]
	if !ok || len(candles) == 0 || since.Add(interval).Before(candles[0].Time) {
		return nil, false
	}

	for i, candle := range candles {
		if !candle.Time.Before(since) {
			return candles[i:], true
		}
	}

	return nil, false
}

func (s *Store) TryGetBalances() (map[string]float64, bool) {
	if s.Balances == nil || len(s.Balances) == 0 {
		return nil, false
	}

	return s.Balances, true
}
