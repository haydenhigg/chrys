package chrys

import (
	"math"
	"testing"
	"time"
)

func assertSlicesEqual(a, b []float64, t *testing.T) {
	for i, va := range a {
		if i >= len(b) {
			t.Errorf("b[%d] does not exist", i)
		} else if vb := b[i]; math.Abs(va-vb) > 1e-6 {
			t.Errorf("a[%d] != b[%d]: %v != %v", i, i, va, vb)
		}
	}

	for i := range b {
		if i >= len(b) {
			t.Errorf("a[%d] does not exist", i)
		}
	}
}

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

func Test_RunBacktest(t *testing.T) {
	// create Scheduler
	scheduler := NewScheduler()

	// mock evaluator
	i := 0
	x := 1.1111111111111111
	firstTime, lastTime := time.Time{}, time.Time{}
	evaluator := func(now time.Time) (float64, error) {
		if i == 0 {
			firstTime = now
		}

		lastTime = now

		// -10% once, then +10% four times, then repeat
		if i%5 == 0 {
			x *= 0.9
		} else {
			x *= 1.1
		}

		i++

		return x, nil
	}

	// RunBacktest()
	start, _ := time.Parse(time.RFC822, "30 Jun 24 12:23 EST")
	end, _ := time.Parse(time.RFC822, "30 Jun 25 8:49 EST")
	test, err := scheduler.RunBacktest(start, end, time.Hour, evaluator)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	expectedFirstTime, _ := time.Parse(time.RFC822, "30 Jun 24 12:00 EST")
	if !firstTime.Equal(expectedFirstTime) {
		t.Errorf("first time != %v: %v", expectedFirstTime, firstTime)
	}

	expectedLastTime, _ := time.Parse(time.RFC822, "30 Jun 25 7:00 EST")
	if !lastTime.Equal(expectedLastTime) {
		t.Errorf("last time != %v: %v", expectedLastTime, lastTime)
	}

	assertSlicesEqual(test.Values[:6], []float64{
		1., 1.1, 1.21, 1.331, 1.4641, 1.31769,
	}, t)

	assertSlicesEqual(test.Returns[:5], []float64{
		.1, .1, .1, .1, -.1,
	}, t)
}
