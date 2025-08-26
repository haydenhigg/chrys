package client

import (
	"github.com/haydenhigg/chrys"
	"time"
)

type Store struct {
	// TODO: use chrys.Feed as the key for Frames
	Frames   map[string]map[time.Duration]*chrys.Frame
	Balances map[string]float64
}

func (store *Store) TryGetFramesSince(
	feed chrys.Feed,
	since time.Time,
) ([]*chrys.Frame, bool) {
	since = since.Truncate(feed.Interval)

	if _, ok := store.Frames[feed.Symbol]; !ok {
		store.Frames[feed.Symbol] = map[time.Duration][]*chrys.Frame{}
	}

	frames, ok := store.Frames[feed.Symbol][feed.Interval]
	if !ok || !frames[0].Time.Before(since.Add(feed.Interval)) {
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
