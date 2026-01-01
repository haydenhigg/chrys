package chrys

import "time"

type IntervalFrameCache = map[time.Duration][]*Frame
type FrameCache map[string]IntervalFrameCache

func (cache FrameCache) GetSince(
	pair *Pair,
	interval time.Duration,
	t time.Time,
) ([]*Frame, bool) {
	// check if pair is in cache
	if _, ok := cache[pair.Name]; !ok {
		return nil, false
	}

	// assert that the frames exist and contain the time requested
	t = t.Truncate(interval)
	frames, ok := cache[pair.Name][interval]
	if !ok || !frames[0].Time.Before(t.Add(interval)) {
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

func (cache FrameCache) GetPriceAt(pair *Pair, t time.Time) (float64, bool) {
	// cycle through all cached intervals for a pair to see if any of them have
	// a price for the given time
	if intervalFrames, ok := cache[pair.Name]; ok {
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

func (cache FrameCache) Set(
	pair *Pair,
	interval time.Duration,
	frames []*Frame,
) {
	// ensure pair is in cache
	if _, ok := cache[pair.Name]; !ok {
		cache[pair.Name] = IntervalFrameCache{}
	}

	cache[pair.Name][interval] = frames
}
