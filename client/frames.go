package client

import (
	"fmt"
	"github.com/haydenhigg/chrys/frame"
	"time"
)

type FrameAPI interface {
	FetchFramesSince(
		pair string,
		interval time.Duration,
		since time.Time,
	) ([]*frame.Frame, error)
}

type PartialFrameCache = map[time.Duration][]*frame.Frame
type FrameCache = map[string]PartialFrameCache

type FrameStore struct {
	api   FrameAPI
	cache map[string]PartialFrameCache
}

func NewFrameStore(api FrameAPI) *FrameStore {
	return &FrameStore{
		api:   api,
		cache: FrameCache{},
	}
}

func (store *FrameStore) getCachedSince(
	pair string,
	interval time.Duration,
	t time.Time,
) ([]*frame.Frame, bool) {
	t = t.Truncate(interval)

	// check if pair is in cache
	if _, ok := store.cache[pair]; !ok {
		return nil, false
	}

	// check if interval is in pair's partial cache
	frames, ok := store.cache[pair][interval]
	if !ok {
		return nil, false
	}

	// check if frames contain necessary time
	// TODO: check if
	if !frames[0].Time.Before(t.Add(interval)) {
		return nil, false
	}

	// chop off older frames
	for i, frame := range frames {
		if !frame.Time.Before(t) {
			return frames[i:], true
		}
	}

	return nil, false
}

func (store *FrameStore) GetSince(
	pair string,
	interval time.Duration,
	t time.Time,
) ([]*frame.Frame, error) {
	t = t.Truncate(interval)

	// check cache
	if frames, ok := store.getCachedSince(pair, interval, t); ok {
		fmt.Println("cache hit", pair, interval, t)
		return frames, nil
	}

	fmt.Println("cache miss", pair, interval, t)

	// retrieve from data source
	frames, err := store.api.FetchFramesSince(pair, interval, t)
	if err != nil {
		return nil, err
	}

	// cache retrieved data
	if _, ok := store.cache[pair]; !ok {
		store.cache[pair] = PartialFrameCache{}
	}

	store.cache[pair][interval] = frames

	return frames, nil
}

func (store *FrameStore) GetNBefore(
	pair string,
	interval time.Duration,
	n int,
	t time.Time,
) ([]*frame.Frame, error) {
	// TODO: t.Truncate(interval).Add(...) ?
	t = t.Add(time.Duration(-n) * interval)

	frames, err := store.GetSince(pair, interval, t)
	if err != nil {
		return nil, err
	}

	return frames[:n], nil
}

func (store *FrameStore) getCachedPriceAt(
	pair string,
	t time.Time,
) (float64, bool) {
	// check all cached intervals for a pair to find a price for the given time
	if intervalFrames, ok := store.cache[pair]; ok {
		frameTime := t.Truncate(time.Minute)

		for interval, frames := range intervalFrames {
			priorFrameTime := frameTime.Add(-interval)

			for _, frame := range frames {
				if frame.Time.Equal(priorFrameTime) {
					return frame.Close, true
				}
			}
		}
	}

	return 0, false
}

func (store *FrameStore) GetPriceAt(pair string, t time.Time) (float64, error) {
	// check cache
	price, ok := store.getCachedPriceAt(pair, t)
	if ok {
		return price, nil
	}

	// retrieve and cache data
	frames, err := store.GetNBefore(pair, time.Minute, 1, t)
	if err != nil {
		return 0, err
	}

	return frames[len(frames)-1].Close, nil
}

func mergeFrames(a, b []*frame.Frame) []*frame.Frame {
	merged := make([]*frame.Frame, 0, len(a)+len(b))
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		if a[i].Time.Before(b[j].Time) {
			merged = append(merged, a[i])
			i++
		} else if b[j].Time.Before(a[i].Time) {
			merged = append(merged, b[j])
			j++
		} else {
			merged = append(merged, a[i])
			i++
			j++
		}
	}

	merged = append(merged, a[i:]...)
	merged = append(merged, b[j:]...)

	return merged
}

func (store *FrameStore) Set(
	pair string,
	interval time.Duration,
	frames []*frame.Frame,
) {
	// check if pair is in cache
	if _, ok := store.cache[pair]; !ok {
		store.cache[pair] = PartialFrameCache{
			interval: frames,
		}
		return
	}

	// check if interval is in pair's partial cache {
	oldFrames, ok := store.cache[pair][interval]
	if !ok {
		store.cache[pair][interval] = frames
		return
	}

	// merge frames
	store.cache[pair][interval] = mergeFrames(frames, oldFrames)
}
