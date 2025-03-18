package event

// Event represents an event that can be dispatched and listened to.
type Event interface {
	// Name returns the name of the event.
	Name() string
	
	// Arguments returns the arguments for the event.
	Arguments() map[string]interface{}
	
	// StopPropagation stops the propagation of the event to further listeners.
	StopPropagation()
	
	// IsPropagationStopped returns whether the propagation of this event has been stopped.
	IsPropagationStopped() bool
}

// BaseEvent provides a basic implementation of the Event interface.
type BaseEvent struct {
	name              string
	arguments         map[string]interface{}
	propagationStopped bool
}

// NewEvent creates a new event with the given name and optional arguments.
func NewEvent(name string, arguments ...map[string]interface{}) *BaseEvent {
	var args map[string]interface{}
	if len(arguments) > 0 {
		args = arguments[0]
	} else {
		args = make(map[string]interface{})
	}
	
	return &BaseEvent{
		name:      name,
		arguments: args,
	}
}

// Name returns the name of the event.
func (e *BaseEvent) Name() string {
	return e.name
}

// Arguments returns the arguments for the event.
func (e *BaseEvent) Arguments() map[string]interface{} {
	return e.arguments
}

// StopPropagation stops the propagation of the event to further listeners.
func (e *BaseEvent) StopPropagation() {
	e.propagationStopped = true
}

// IsPropagationStopped returns whether the propagation of this event has been stopped.
func (e *BaseEvent) IsPropagationStopped() bool {
	return e.propagationStopped
}