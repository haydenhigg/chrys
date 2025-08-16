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

func (e *Engine) Run() error {
	signals := make(Signals, len(e.Signalers))

	// evaluate signals
	for key, signaler := range e.Signalers {
		signal, err := signaler()
		if err != nil {
			return err
		}

		signals[key] = signal
	}

	// run handlers
	for _, handler := range e.Handlers {
		if err := handler(signals); err != nil {
			return err
		}
	}

	return nil
}
