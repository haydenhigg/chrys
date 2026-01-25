package store

import (
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

func NewFrames(api FrameAPI) *FrameStore {
	return &FrameStore{
		api:   api,
		cache: FrameCache{},
	}
}

// binary search a sorted frame slice for the first frame that starts at/after t
func searchFrames(frames []*frame.Frame, t time.Time) int {
	low, high := 0, len(frames)-1
	epochs := 0

	for high-low > 0 {
		epochs++
		if frames[low].Time.Equal(t) {
			return low
		} else if frames[high].Time.Equal(t) {
			return high
		}

		if high-low == 1 {
			if frames[low].Time.Before(t) && frames[high].Time.After(t) {
				return high
			} else {
				break
			}
		}

		midIndex := (low + high) / 2

		if frames[midIndex].Time.Before(t) {
			low = midIndex
		} else if frames[midIndex].Time.After(t) {
			high = midIndex
		} else {
			return midIndex
		}
	}

	return -1
}

func (store *FrameStore) getCachedSince(
	pair string,
	interval time.Duration,
	t time.Time,
) ([]*frame.Frame, bool) {
	// check if pair is in cache
	if _, ok := store.cache[pair]; !ok {
		return nil, false
	}

	// check if interval is in pair's partial cache
	frames, ok := store.cache[pair][interval]
	if !ok || len(frames) == 0 {
		return nil, false
	}

	// chop off older frames
	index := searchFrames(frames, t)
	if index > -1 {
		return frames[index:], true
	}

	return nil, false
}

func (store *FrameStore) GetSince(
	pair string,
	interval time.Duration,
	t time.Time,
) ([]*frame.Frame, error) {
	// check cache
	if frames, ok := store.getCachedSince(pair, interval, t); ok {
		return frames, nil
	}

	// retrieve from data source
	frames, err := store.api.FetchFramesSince(pair, interval, t)
	if err != nil {
		return nil, err
	}

	// cache retrieved data
	store.Set(pair, interval, frames)

	return frames, nil
}

func (store *FrameStore) GetNBefore(
	pair string,
	interval time.Duration,
	n int,
	t time.Time,
) ([]*frame.Frame, error) {
	t = t.Truncate(interval).Add(time.Duration(-n) * interval)

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
	// check all cached intervals to find a Close price for the given time
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

// merge and deduplicate two sorted frame slices
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
) *FrameStore {
	// check if pair is in cache
	if _, ok := store.cache[pair]; !ok {
		store.cache[pair] = PartialFrameCache{interval: frames}
		return store
	}

	// check if interval is in pair's partial cache {
	oldFrames, ok := store.cache[pair][interval]
	if !ok {
		store.cache[pair][interval] = frames
		return store
	}

	// merge frames
	store.cache[pair][interval] = mergeFrames(frames, oldFrames)

	return store
}
