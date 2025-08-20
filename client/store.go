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
	since = since.Truncate(interval)

	if _, ok := s.Candles[pair]; !ok {
		s.Candles[pair] = map[time.Duration][]*candle.Candle{}
	}

	candles, ok := s.Candles[pair][interval]
	if !ok || !candles[0].Time.Before(since.Add(interval)) {
		return nil, false
	}

	for i, candle := range candles {
		if !candle.Time.Before(since) {
			return candles[i:], true
		}
	}

	return nil, false
}

func (s *Store) TryGetPrice(pair string, now time.Time) (float64, bool) {
	price := 0.
	ok := false

	if intervalCandles, ok := s.Candles[pair]; ok {
		for interval, candles := range intervalCandles {
			latestCandleTime := now.Truncate(interval).Add(-interval)

			for _, candle := range candles {
				if candle.Time.Equal(latestCandleTime) {
					price = candle.Close
					ok = true

					break
				}
			}

			if ok {
				break
			}
		}
	}

	return price, ok
}

func (s *Store) TryGetBalances() (map[string]float64, bool) {
	if s.Balances == nil || len(s.Balances) == 0 {
		return nil, false
	}

	return s.Balances, true
}
