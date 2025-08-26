package client

import (
	"github.com/haydenhigg/chrys"
	"time"
)

type Store struct {
	Frames   map[string]map[time.Duration][]*chrys.Frame
	Balances map[string]float64
}

func (store *Store) TryGetFramesSince(
	series *chrys.Series,
	since time.Time,
) ([]*chrys.Frame, bool) {
	since = since.Truncate(series.Interval)

	if _, ok := store.Frames[series.Pair.String()]; !ok {
		store.Frames[series.Pair.String()] = map[time.Duration][]*chrys.Frame{}
	}

	frames, ok := store.Frames[series.Pair.String()][series.Interval]
	if !ok || !frames[0].Time.Before(since.Add(series.Interval)) {
		return nil, false
	}

	for i, frame := range frames {
		if !frame.Time.Before(since) {
			return frames[i:], true
		}
	}

	return nil, false
}

func (store *Store) TryGetPriceAt(pair string, now time.Time) (float64, bool) {
	price := 0.
	ok := false

	if intervalFrames, ok := store.Frames[pair]; ok {
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

func (store *Store) TryGetBalances() (map[string]float64, bool) {
	if store.Balances == nil || len(store.Balances) == 0 {
		return nil, false
	}

	return store.Balances, true
}
