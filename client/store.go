package client

import (
	"chrys"
	"time"
)

type Pair = string

type Store struct {
	Frames  map[Pair]map[time.Duration][]*chrys.Frame
	Balances map[string]float64
}

func (s *Store) TryGetFramesSince(
	pair string,
	interval time.Duration,
	since time.Time,
) ([]*chrys.Frame, bool) {
	since = since.Truncate(interval)

	if _, ok := s.Frames[pair]; !ok {
		s.Frames[pair] = map[time.Duration][]*chrys.Frame{}
	}

	frames, ok := s.Frames[pair][interval]
	if !ok || !frames[0].Time.Before(since.Add(interval)) {
		return nil, false
	}

	for i, frame := range frames {
		if !frame.Time.Before(since) {
			return frames[i:], true
		}
	}

	return nil, false
}

func (s *Store) TryGetPriceAt(pair string, now time.Time) (float64, bool) {
	price := 0.
	ok := false

	if intervalFrames, ok := s.Frames[pair]; ok {
		for interval, frames := range intervalFrames {
			latestFrameTime := now.Truncate(interval).Add(-interval)

			for _, frame := range frames {
				if frame.Time.Equal(latestFrameTime) {
					price = frame.Close
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
