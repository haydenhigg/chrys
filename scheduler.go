package chrys

import "time"

type Block = func(now time.Time) error
type Scheduler map[time.Duration][]Block

// initializer
func NewScheduler() Scheduler {
	return Scheduler{}
}

func (scheduler Scheduler) Add(interval time.Duration, block Block) Scheduler {
	if blocks, ok := scheduler[interval]; ok {
		scheduler[interval] = append(blocks, block)
	} else {
		scheduler[interval] = []Block{block}
	}

	return scheduler
}

func (scheduler Scheduler) Run(now time.Time) error {
	t := now.Truncate(time.Minute)

	for interval, blocks := range scheduler {
		if t.Truncate(interval).Equal(t) {
			for _, block := range blocks {
				if err := block(t); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (scheduler Scheduler) RunBacktest(
	start time.Time,
	end time.Time,
	step time.Duration,
	evaluator func(time.Time) (float64, error),
) (*Backtest, error) {
	n := int64(end.Sub(start)) / int64(step)

	test := &Backtest{
		Start:   start.Truncate(step),
		End:     end.Truncate(step),
		Step:    step,
		Values:  make([]float64, int(n)),
		Returns: make([]float64, int(n-1)),
	}

	var err error
	for i, t := int64(0), test.Start; i < n; i, t = i+1, t.Add(step) {
		// run
		err = scheduler.Run(t)
		if err != nil {
			return test, err
		}

		// evaluate
		test.Values[i], err = evaluator(t)
		if err != nil {
			return test, err
		}

		if i > 0 {
			test.Returns[i-1] = test.Values[i]/test.Values[i-1] - 1
		}
	}

	return test, nil
}
