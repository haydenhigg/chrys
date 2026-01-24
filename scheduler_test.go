package chrys

import (
	"testing"
	"time"
)

func Test_Add(t *testing.T) {
	// create Scheduler
	scheduler := NewScheduler()

	// Add()
	scheduler.Add(time.Minute, func(now time.Time) error { return nil })

	// assert
	if blocks, ok := scheduler[time.Minute]; !ok {
		t.Errorf("scheduler[time.Minute] does not exist")
	} else {
		if len(blocks) != 1 {
			t.Errorf("len(scheduler[time.Minute]) != 1: %d", len(blocks))
		}
	}
}

func Test_Run(t *testing.T) {
	// create Scheduler
	scheduler := NewScheduler()

	// mock
	didRunMin, didRun5Min, didRun15Min := false, false, false
	scheduler.Add(time.Minute, func(now time.Time) error {
		didRunMin = true
		return nil
	}).Add(5*time.Minute, func(now time.Time) error {
		didRun5Min = true
		return nil
	}).Add(15*time.Minute, func(now time.Time) error {
		didRun15Min = true
		return nil
	})

	// Run()
	now := time.Now().Truncate(15 * time.Minute).Add(5 * time.Minute)
	if err := scheduler.Run(now); err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	if !didRunMin {
		t.Errorf("time.Minute Block did not run")
	}

	if !didRun5Min {
		t.Errorf("5 * time.Minute Block did not run")
	}

	if didRun15Min {
		t.Errorf("15 * time.Minute Block did run")
	}
}
