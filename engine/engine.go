package engine

type Signals = map[string]float64

type Signaler = func() (float64, error)
type Handler = func(signals Signals) error

type Engine struct {
	Signalers map[string]Signaler
	Handlers  []Handler
}

func New() *Engine {
	return &Engine{
		Signalers: map[string]Signaler{},
		Handlers:  []Handler{},
	}
}

func (e *Engine) Signal(key string, signaler Signaler) *Engine {
	e.Signalers[key] = signaler
	return e
}

func (e *Engine) Handle(handler Handler) *Engine {
	e.Handlers = append(e.Handlers, handler)
	return e
}

func (e *Engine) Run() []error {
	signals := make(Signals, len(e.Signalers))
	errors := []error{}

	// evaluate signals
	for key, signaler := range e.Signalers {
		if signal, err := signaler(); err == nil {
			signals[key] = signal
		} else {
			errors = append(errors, err)
		}
	}

	// run handlers
	for _, handler := range e.Handlers {
		if err := handler(signals); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
