package chrys

import "time"

type Stage = func(now time.Time) error

type Pipeline struct {
	Stages []Stage
	Data   map[string]float64
}

func New() *Pipeline {
	return &Pipeline{Stages: []Stage{}, Data: map[string]float64{}}
}

func (p *Pipeline) AddStage(handler Stage) *Pipeline {
	p.Stages = append(p.Stages, handler)
	return p
}

func (p *Pipeline) Get(k string) float64 {
	return p.Data[k]
}

func (p *Pipeline) Set(k string, v float64) {
	p.Data[k] = v
}

func (p *Pipeline) Run(now time.Time) error {
	for _, handler := range p.Stages {
		if err := handler(now); err != nil {
			return err
		}
	}

	return nil
}
