package chrys

import "time"

type Block = func(now time.Time) error
type Pipeline struct {
	Blocks []Block
}

func NewPipeline() *Pipeline {
	return &Pipeline{Blocks: []Block{}}
}

func (pipeline *Pipeline) AddBlock(handler Block) *Pipeline {
	pipeline.Blocks = append(pipeline.Blocks, handler)
	return pipeline
}

func (pipeline *Pipeline) Run(t time.Time) error {
	for _, handler := range pipeline.Blocks {
		if err := handler(t); err != nil {
			return err
		}
	}

	return nil
}

type BacktestReport struct {
	Start, End           time.Time
	StartValue, EndValue float64
	EndEquity,
	Return float64
}

func (pipeline *Pipeline) RunBacktest(
	start, end time.Time,
	out []string,
) (BacktestReport, error) {
	return BacktestReport{}, nil
}
