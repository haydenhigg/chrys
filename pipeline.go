package chrys

import "time"

type Stage = func(now time.Time) error

type Pipeline struct {
	Data   map[string]float64
	Stages []Stage
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		Data:   map[string]float64{},
		Stages: []Stage{},
	}
}

func (pipeline *Pipeline) Get(k string) float64 {
	return pipeline.Data[k]
}

func (pipeline *Pipeline) Set(k string, v float64) *Pipeline {
	pipeline.Data[k] = v
	return pipeline
}

func (pipeline *Pipeline) AddStage(handler Stage) *Pipeline {
	pipeline.Stages = append(pipeline.Stages, handler)
	return pipeline
}

func (pipeline *Pipeline) Run(t time.Time) error {
	for _, handler := range pipeline.Stages {
		if err := handler(t); err != nil {
			return err
		}
	}

	return nil
}

// type BacktestReport {
// 	StartEquity float64
// 	EndEquity   float64
// 	Return      float64
// }

// func (pipeline *Pipeline) RunBacktest(
// 	start,
// 	end time.Time,
// 	out *Asset,
// ) (*BacktestReport, error) {
// 	return &BacktestReport{}, nil
// }
