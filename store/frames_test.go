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
	return frames, nil
}

func Test_GetSinceUncachedExactTime(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockFrameAPI{
		callback: func() {
			didUseAPI = true
		},
	}

	// GetSince()
	since := time.Now().Truncate(time.Hour).Add(-3 * time.Hour)
	frames, err := NewFrames(mockAPI).GetSince("BTC/USD", time.Hour, since)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	if len(frames) != 3 {
		t.Errorf("len(frames) != 3: %d", len(frames))
	}

	if !didUseAPI {
		t.Errorf("cache was hit")
	}
}
