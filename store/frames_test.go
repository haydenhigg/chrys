package store

import (
	"github.com/haydenhigg/chrys/frame"
	"testing"
	"time"
)

// mock
type MockFrameAPI struct {
	callback func()
}

func (api MockFrameAPI) FetchFramesSince(
	pair string,
	interval time.Duration,
	since time.Time,
) ([]*frame.Frame, error) {
	if api.callback != nil {
		api.callback()
	}

	frames := []*frame.Frame{}

	// all frames that start !Before(since) and end !After(now-interval)
	now := time.Now().Add(-interval)
	for t := since; !t.After(now); t = t.Add(interval) {
		if t.Equal(t.Truncate(interval)) {
			frames = append(frames, &frame.Frame{Time: t})
		} else {
			t = t.Truncate(interval)
		}
	}

	return frames, nil
}

// tests
func assertFrameSliceEqual(a, b []*frame.Frame, t *testing.T) {
	n := max(len(a), len(b))

	for i := range n {
		if i >= len(a) {
			t.Errorf(`a[%d] does not exist`, i)
		} else if i >= len(b) {
			t.Errorf(`b[%d] does not exist`, i)
		} else if a[i].Time != b[i].Time {
			t.Errorf(
				`a[%d].Time != b[%d].Time: %v != %v`,
				i,
				i,
				a[i].Time,
				b[i].Time,
			)
		}
	}
}

// tests -> GetSince -> Exact Time
func Test_GetSinceExactTimeUncached(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockFrameAPI{callback: func() { didUseAPI = true }}

	// GetSince()
	since := time.Now().Truncate(time.Hour).Add(-3 * time.Hour)
	frames, err := NewFrames(mockAPI).GetSince("BTC/USD", time.Hour, since)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	assertFrameSliceEqual(frames, []*frame.Frame{
		{Time: since},
		{Time: since.Add(time.Hour)},
		{Time: since.Add(2 * time.Hour)},
	}, t)

	if !didUseAPI {
		t.Errorf("cache was hit")
	}
}

func Test_GetSinceExactTimeCached(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockFrameAPI{callback: func() { didUseAPI = true }}

	// GetSince()
	since := time.Now().Truncate(time.Hour).Add(-10 * time.Hour)
	store := NewFrames(mockAPI)
	store.GetSince("BTC/USD", time.Hour, since)

	// reset didUseAPI
	didUseAPI = false

	// GetSince() again
	since = since.Add(7 * time.Hour) // time.Now().Add(-3 * time.Hour)
	frames, err := store.GetSince("BTC/USD", time.Hour, since)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	assertFrameSliceEqual(frames, []*frame.Frame{
		{Time: since},
		{Time: since.Add(time.Hour)},
		{Time: since.Add(2 * time.Hour)},
	}, t)

	if didUseAPI {
		t.Errorf("cache was not hit")
	}
}

func Test_GetSinceExactTimeCachedMiss(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockFrameAPI{callback: func() { didUseAPI = true }}

	// GetSince()
	since := time.Now().Truncate(time.Hour).Add(-time.Hour)
	store := NewFrames(mockAPI)
	store.GetSince("BTC/USD", time.Hour, since)

	// reset didUseAPI
	didUseAPI = false

	// GetSince() again
	since = since.Add(-2 * time.Hour)
	frames, err := store.GetSince("BTC/USD", time.Hour, since)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	assertFrameSliceEqual(frames, []*frame.Frame{
		{Time: since},
		{Time: since.Add(time.Hour)},
		{Time: since.Add(2 * time.Hour)},
	}, t)

	if !didUseAPI {
		t.Errorf("cache was hit")
	}
}

// tests -> GetSince -> Inexact Time
func Test_GetSinceInexactTimeUncached(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockFrameAPI{callback: func() { didUseAPI = true }}

	// GetSince()
	since := time.Now().Add(-3 * time.Hour)
	frames, err := NewFrames(mockAPI).GetSince("BTC/USD", time.Hour, since)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	truncatedSince := since.Truncate(time.Hour)
	assertFrameSliceEqual(frames, []*frame.Frame{
		{Time: truncatedSince.Add(time.Hour)},
		{Time: truncatedSince.Add(2 * time.Hour)},
	}, t)

	if !didUseAPI {
		t.Errorf("cache was hit")
	}
}

func Test_GetSinceInexactTimeCached(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockFrameAPI{callback: func() { didUseAPI = true }}

	// GetSince()
	since := time.Now().Add(-10 * time.Hour)
	store := NewFrames(mockAPI)
	store.GetSince("BTC/USD", time.Hour, since)

	// reset didUseAPI
	didUseAPI = false

	// GetSince() again
	since = since.Add(7 * time.Hour) // time.Now().Add(-3 * time.Hour)
	frames, err := store.GetSince("BTC/USD", time.Hour, since)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	truncatedSince := since.Truncate(time.Hour)
	assertFrameSliceEqual(frames, []*frame.Frame{
		{Time: truncatedSince.Add(time.Hour)},
		{Time: truncatedSince.Add(2 * time.Hour)},
	}, t)

	if didUseAPI {
		t.Errorf("cache was not hit")
	}
}

// tests -> GetSince -> Inexact Time
