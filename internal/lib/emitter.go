package lib

// Emitter provides a type that can have callback functions attached to string-based events.
type Emitter struct {
	handlers map[string][]func(Event)
}

// NewEmitter creates a new Emitter.
func NewEmitter() *Emitter {
	return &Emitter{
		handlers: make(map[string][]func(Event)),
	}
}

// Emit emits the given event with the provided data, calling all registered handlers.
func (e *Emitter) Emit(event string, data Event) {
	if h, ok := e.handlers[event]; ok {
		for _, f := range h {
			(f)(data)
		}
	}
}

// On adds an event handler to a given event string.
func (e *Emitter) On(event string, cb func(Event)) {
	if _, ok := e.handlers[event]; !ok {
		e.handlers[event] = make([]func(Event), 0)
	}
	e.handlers[event] = append(e.handlers[event], cb)
}
