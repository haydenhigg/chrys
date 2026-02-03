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
	i := 1.

	// all frames that start !Before(since) and end !After(now-interval)
	now := time.Now().Add(-interval)
	for t := since; !t.After(now); t = t.Add(interval) {
		if t.Equal(t.Truncate(interval)) {
			frames = append(frames, &frame.Frame{Time: t, Close: i})
			i++
		} else {
			t = t.Truncate(interval)
		}
	}

	return frames, nil
}

// tests
func assertFrameTimesEqual(a, b []*frame.Frame, t *testing.T) {
	n := max(len(a), len(b))

	for i := range n {
		if i >= len(a) {
			t.Errorf(`a[%d] does not exist`, i)
		} else if i >= len(b) {
			t.Errorf(`b[%d] does not exist`, i)
		} else if a[i].Time != b[i].Time {
			t.Errorf(
				`a[%d].Time != b[%d].Time: %v != %v`,
				i, i, a[i].Time, b[i].Time,
			)
		}
	}
}

func assertFrameClosesEqual(a, b []*frame.Frame, t *testing.T) {
	n := max(len(a), len(b))

	for i := range n {
		if i >= len(a) {
			t.Errorf(`a[%d] does not exist`, i)
		} else if i >= len(b) {
			t.Errorf(`b[%d] does not exist`, i)
		} else if a[i].Close != b[i].Close {
			t.Errorf(
				`a[%d].Close != b[%d].Close: %v != %v`,
				i, i, a[i].Close, b[i].Close,
			)
		}
	}
}

// tests -> GetSince
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
	assertFrameTimesEqual(frames, []*frame.Frame{
		{Time: since},
		{Time: since.Add(time.Hour)},
		{Time: since.Add(2 * time.Hour)},
	}, t)

	if !didUseAPI {
		t.Error("cache was hit")
	}
}

func Test_GetSinceExactTimeCached(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockFrameAPI{callback: func() { didUseAPI = true }}

	// set up store
	store := NewFrames(mockAPI)

	// GetSince()
	since := time.Now().Truncate(time.Hour).Add(-10 * time.Hour)
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
	assertFrameTimesEqual(frames, []*frame.Frame{
		{Time: since},
		{Time: since.Add(time.Hour)},
		{Time: since.Add(2 * time.Hour)},
	}, t)

	if didUseAPI {
		t.Error("cache was not hit")
	}
}

func Test_GetSinceExactTimeCachedMiss(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockFrameAPI{callback: func() { didUseAPI = true }}

	// set up store
	store := NewFrames(mockAPI)

	// GetSince()
	since := time.Now().Truncate(time.Hour).Add(-time.Hour)
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
	assertFrameTimesEqual(frames, []*frame.Frame{
		{Time: since},
		{Time: since.Add(time.Hour)},
		{Time: since.Add(2 * time.Hour)},
	}, t)

	if !didUseAPI {
		t.Error("cache was hit")
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
	assertFrameTimesEqual(frames, []*frame.Frame{
		{Time: truncatedSince.Add(time.Hour)},
		{Time: truncatedSince.Add(2 * time.Hour)},
	}, t)

	if !didUseAPI {
		t.Error("cache was hit")
	}
}

func Test_GetSinceInexactTimeCached(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockFrameAPI{callback: func() { didUseAPI = true }}

	// set up store
	store := NewFrames(mockAPI)

	// GetSince()
	since := time.Now().Add(-10 * time.Hour)
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
	assertFrameTimesEqual(frames, []*frame.Frame{
		{Time: truncatedSince.Add(time.Hour)},
		{Time: truncatedSince.Add(2 * time.Hour)},
	}, t)

	if didUseAPI {
		t.Error("cache was not hit")
	}
}

func Test_GetSinceInexactTimeCachedMiss(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockFrameAPI{callback: func() { didUseAPI = true }}

	// set up store
	store := NewFrames(mockAPI)

	// GetSince()
	since := time.Now().Add(-time.Hour)
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
	truncatedSince := since.Truncate(time.Hour)
	assertFrameTimesEqual(frames, []*frame.Frame{
		{Time: truncatedSince.Add(time.Hour)},
		{Time: truncatedSince.Add(2 * time.Hour)},
	}, t)

	if !didUseAPI {
		t.Error("cache was hit")
	}
}

// tests -> GetNBefore
// tests -> GetNBefore -> Exact Time
func Test_GetNBeforeExactTime(t *testing.T) {
	// set up store
	store := NewFrames(MockFrameAPI{})

	// GetNBefore()
	now := time.Now().Truncate(time.Hour)
	frames, err := store.GetNBefore("BTC/USD", time.Hour, 5, now)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	assertFrameTimesEqual(frames, []*frame.Frame{
		{Time: now.Add(-5 * time.Hour)},
		{Time: now.Add(-4 * time.Hour)},
		{Time: now.Add(-3 * time.Hour)},
		{Time: now.Add(-2 * time.Hour)},
		{Time: now.Add(-time.Hour)},
	}, t)
}

// tests -> GetNBefore -> Inexact Time
func Test_GetNBeforeInexactTime(t *testing.T) {
	// set up store
	store := NewFrames(MockFrameAPI{})

	// GetNBefore()
	now := time.Now()
	frames, err := store.GetNBefore("BTC/USD", time.Minute, 5, now)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	truncatedNow := now.Truncate(time.Minute)
	assertFrameTimesEqual(frames, []*frame.Frame{
		{Time: truncatedNow.Add(-5 * time.Minute)},
		{Time: truncatedNow.Add(-4 * time.Minute)},
		{Time: truncatedNow.Add(-3 * time.Minute)},
		{Time: truncatedNow.Add(-2 * time.Minute)},
		{Time: truncatedNow.Add(-time.Minute)},
	}, t)
}

// tests -> GetPriceAt
// tests -> GetPriceAt -> Exact Time
func Test_GetPriceAtExactTimeUncached(t *testing.T) {
	// set up store
	store := NewFrames(MockFrameAPI{})

	// GetPriceAt()
	now := time.Now().Truncate(time.Minute)
	price, err := store.GetPriceAt("BTC/USD", now)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	expectedPrice := 1.
	if price != expectedPrice {
		t.Errorf("price != expectedPrice: %f != %f", price, expectedPrice)
	}

	assertFrameTimesEqual(store.Cache["BTC/USD"][time.Minute], []*frame.Frame{
		{Time: now.Add(-time.Minute)},
	}, t)
}

func Test_GetPriceAtExactTimeCached(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockFrameAPI{callback: func() { didUseAPI = true }}

	// set up store
	store := NewFrames(mockAPI)

	// GetNBefore()
	now := time.Now().Truncate(30 * time.Minute)
	store.GetNBefore("BTC/USD", 30*time.Minute, 5, now)

	// reset didUseAPI
	didUseAPI = false

	// GetPriceAt()
	price, err := store.GetPriceAt("BTC/USD", now)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	expectedPrice := 5.
	if price != expectedPrice {
		t.Errorf("price != expectedPrice: %f != %f", price, expectedPrice)
	}

	assertFrameTimesEqual(store.Cache["BTC/USD"][30*time.Minute], []*frame.Frame{
		{Time: now.Add(-5 * 30 * time.Minute)},
		{Time: now.Add(-4 * 30 * time.Minute)},
		{Time: now.Add(-3 * 30 * time.Minute)},
		{Time: now.Add(-2 * 30 * time.Minute)},
		{Time: now.Add(-30 * time.Minute)},
	}, t)

	if didUseAPI {
		t.Error("cache was not hit")
	}
}

func Test_GetPriceAtExactTimeCachedMiss(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockFrameAPI{callback: func() { didUseAPI = true }}

	// set up store
	store := NewFrames(mockAPI)

	// GetNBefore()
	now := time.Now().Truncate(30 * time.Minute).Add(time.Minute)
	store.GetNBefore("BTC/USD", 30*time.Minute, 5, now)

	// reset didUseAPI
	didUseAPI = false

	// GetPriceAt()
	price, err := store.GetPriceAt("BTC/USD", now)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	expectedPrice := 1.
	if price != expectedPrice {
		t.Errorf("price != expectedPrice: %f != %f", price, expectedPrice)
	}

	assertFrameTimesEqual(
		store.Cache["BTC/USD"][time.Minute][:1],
		[]*frame.Frame{
			{Time: now.Add(-time.Minute)},
		},
		t,
	)

	if !didUseAPI {
		t.Error("cache was hit")
	}
}

// tests -> GetPriceAt -> Inexact Time
func Test_GetPriceAtInexactTimeUncached(t *testing.T) {
	// set up store
	store := NewFrames(MockFrameAPI{})

	// GetPriceAt()
	now := time.Now()
	price, err := store.GetPriceAt("BTC/USD", now)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	expectedPrice := 1.
	if price != expectedPrice {
		t.Errorf("price != expectedPrice: %f != %f", price, expectedPrice)
	}

	assertFrameTimesEqual(store.Cache["BTC/USD"][time.Minute], []*frame.Frame{
		{Time: now.Truncate(time.Minute).Add(-time.Minute)},
	}, t)
}

func Test_GetPriceAtInexactTimeCached(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockFrameAPI{callback: func() { didUseAPI = true }}

	// set up store
	store := NewFrames(mockAPI)

	// GetNBefore()
	now := time.Now().Truncate(30 * time.Minute).Add(37 * time.Second)
	store.GetNBefore("BTC/USD", 30*time.Minute, 5, now)

	// reset didUseAPI
	didUseAPI = false

	// GetPriceAt()
	price, err := store.GetPriceAt("BTC/USD", now)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	expectedPrice := 5.
	if price != expectedPrice {
		t.Errorf("price != expectedPrice: %f != %f", price, expectedPrice)
	}

	truncatedNow := now.Truncate(time.Minute)
	assertFrameTimesEqual(store.Cache["BTC/USD"][30*time.Minute], []*frame.Frame{
		{Time: truncatedNow.Add(-5 * 30 * time.Minute)},
		{Time: truncatedNow.Add(-4 * 30 * time.Minute)},
		{Time: truncatedNow.Add(-3 * 30 * time.Minute)},
		{Time: truncatedNow.Add(-2 * 30 * time.Minute)},
		{Time: truncatedNow.Add(-30 * time.Minute)},
	}, t)

	if didUseAPI {
		t.Error("cache was not hit")
	}
}

func Test_GetPriceAtInexactTimeCachedMiss(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockFrameAPI{callback: func() { didUseAPI = true }}

	// set up store
	store := NewFrames(mockAPI)

	// GetNBefore()
	now := time.Now().Truncate(30 * time.Minute).
		Add(time.Minute).
		Add(37 * time.Second)
	store.GetNBefore("BTC/USD", 30*time.Minute, 5, now)

	// reset didUseAPI
	didUseAPI = false

	// GetPriceAt()
	price, err := store.GetPriceAt("BTC/USD", now)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	expectedPrice := 1.
	if price != expectedPrice {
		t.Errorf("price != expectedPrice: %f != %f", price, expectedPrice)
	}

	assertFrameTimesEqual(
		store.Cache["BTC/USD"][time.Minute][:1],
		[]*frame.Frame{
			{Time: now.Truncate(time.Minute).Add(-time.Minute)},
		},
		t,
	)

	if !didUseAPI {
		t.Error("cache was hit")
	}
}

// tests -> Set
func Test_SetInitial(t *testing.T) {
	// set up store
	store := NewFrames(MockFrameAPI{})

	// Set()
	now := time.Now().Truncate(time.Hour)
	newFrames, _ := store.api.FetchFramesSince(
		"BTC/USD",
		time.Hour,
		now.Add(-5*time.Hour),
	)

	store.Set("BTC/USD", time.Hour, newFrames)

	// assert
	expectedFrames := []*frame.Frame{
		{Time: now.Add(-5 * time.Hour), Close: 1.},
		{Time: now.Add(-4 * time.Hour), Close: 2.},
		{Time: now.Add(-3 * time.Hour), Close: 3.},
		{Time: now.Add(-2 * time.Hour), Close: 4.},
		{Time: now.Add(-time.Hour), Close: 5.},
	}

	assertFrameTimesEqual(store.Cache["BTC/USD"][time.Hour], expectedFrames, t)
	assertFrameClosesEqual(store.Cache["BTC/USD"][time.Hour], expectedFrames, t)
}

func Test_SetMerge(t *testing.T) {
	// set up store
	store := NewFrames(MockFrameAPI{})

	// Set()
	now := time.Now().Truncate(time.Hour)
	newFrames, _ := store.api.FetchFramesSince(
		"BTC/USD",
		time.Hour,
		now.Add(-5*time.Hour),
	)

	store.Set("BTC/USD", time.Hour, newFrames)

	// Set() again
	newFrames, _ = store.api.FetchFramesSince(
		"BTC/USD",
		time.Hour,
		now.Add(-2*time.Hour),
	)

	store.Set("BTC/USD", time.Hour, newFrames)

	// assert
	expectedFrames := []*frame.Frame{
		{Time: now.Add(-5 * time.Hour), Close: 1.},
		{Time: now.Add(-4 * time.Hour), Close: 2.},
		{Time: now.Add(-3 * time.Hour), Close: 3.},
		{Time: now.Add(-2 * time.Hour), Close: 1.},
		{Time: now.Add(-time.Hour), Close: 2.},
	}

	assertFrameTimesEqual(store.Cache["BTC/USD"][time.Hour], expectedFrames, t)
	assertFrameClosesEqual(store.Cache["BTC/USD"][time.Hour], expectedFrames, t)
}
