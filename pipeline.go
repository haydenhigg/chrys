package chrys

import "time"

type Stage = func(now time.Time) error

type Pipeline struct {
	Stages []Stage
	Data   map[string]float64
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		Stages: []Stage{},
		Data:   map[string]float64{},
	}
}

func (p *Pipeline) AddStage(handler Stage) *Pipeline {
	p.Stages = append(p.Stages, handler)
	return p
}

func (p *Pipeline) Get(k string) float64 {
	return p.Data[k]
}

func (p *Pipeline) Set(k string, v float64) *Pipeline {
	p.Data[k] = v
	return p
}

func (p *Pipeline) Run(t time.Time) error {
	for _, handler := range p.Stages {
		if err := handler(t); err != nil {
			return err
		}
	}

	return nil
}
