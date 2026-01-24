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
