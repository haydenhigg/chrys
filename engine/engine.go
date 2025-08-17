package engine

import "time"

type Handler = func(now time.Time) error

type Engine struct {
	Handlers []Handler
	Data     map[string]float64
}

func New() *Engine {
	return &Engine{Handlers: []Handler{}, Data: map[string]float64{}}
}

func (e *Engine) Handle(handler Handler) *Engine {
	e.Handlers = append(e.Handlers, handler)
	return e
}

func (e *Engine) Get(k string) float64 {
	return e.Data[k]
}

func (e *Engine) Set(k string, v float64) {
	e.Data[k] = v
}

func (e *Engine) Run(now time.Time) error {
	for _, handler := range e.Handlers {
		if err := handler(now); err != nil {
			return err
		}
	}

	return nil
}
