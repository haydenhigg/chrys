package client

import (
	"github.com/haydenhigg/chrys"
	"time"
	"fmt"
)

type IntervalFrameCache = map[time.Duration][]*chrys.Frame
type FrameCache map[string]IntervalFrameCache

type FrameConnector interface {
	FetchFramesSince(
		pair *chrys.Pair,
		interval time.Duration,
		since time.Time,
	) ([]*chrys.Frame, error)
}

type Frames struct {
	Connector FrameConnector
	Cache FrameCache
}

func NewFrames(connector FrameConnector) *Frames {
	return &Frames{
		Connector: connector,
		Cache: FrameCache{},
	}
}

func (store *Frames) GetCachedSince(
	pair *chrys.Pair,
	interval time.Duration,
	t time.Time,
) ([]*chrys.Frame, bool) {
	// check if pair is in cache
	if _, ok := store.Cache[pair.Name]; !ok {
		return nil, false
	}

	// assert that the frames exist and contain the time requested
	t = t.Truncate(interval)
	frames, ok := store.Cache[pair.Name][interval]
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

func (store *Frames) GetSince(
	pair *chrys.Pair,
	interval time.Duration,
	t time.Time,
) ([]*chrys.Frame, error) {
	t = t.Truncate(interval)

	// check cache
	if frames, ok := store.GetCachedSince(pair, interval, t); ok {
		fmt.Println("cache hit", pair.Name, interval, t)
		return frames, nil
	}

	// retrieve from data source
	frames, err := store.Connector.FetchFramesSince(pair, interval, t)
	if err != nil {
		return nil, err
	}

	// cache retrieved data
	fmt.Println("cache miss", pair.Name, interval, t)
	store.Set(pair, interval, frames)

	return frames, nil
}

func (store *Frames) GetNBefore(
	pair *chrys.Pair,
	interval time.Duration,
	n int,
	t time.Time,
) ([]*chrys.Frame, error) {
	t = t.Add(time.Duration(-n) * interval)

	frames, err := store.GetSince(pair, interval, t)
	if err != nil {
		return nil, err
	}

	return frames[:n], nil
}


func (store *Frames) GetCachedPriceAt(
	pair *chrys.Pair,
	t time.Time,
) (float64, bool) {
	// cycle through all cached intervals for a pair to see if any of them have
	// a price for the given time
	if intervalFrames, ok := store.Cache[pair.Name]; ok {
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

func (store *Frames) GetPriceAt(
	pair *chrys.Pair,
	t time.Time,
) (float64, error) {
	price, ok := store.GetCachedPriceAt(pair, t)
	if !ok {
		frames, err := store.GetNBefore(pair, time.Minute, 1, t)
		if err != nil {
			return 0, err
		}

		price = frames[len(frames)-1].Close
	}

	return price, nil
}

func (store *Frames) Set(
	pair *chrys.Pair,
	interval time.Duration,
	frames []*chrys.Frame,
) {
	// ensure pair is in cache
	if _, ok := store.Cache[pair.Name]; !ok {
		store.Cache[pair.Name] = IntervalFrameCache{}
	}

	store.Cache[pair.Name][interval] = frames
}
