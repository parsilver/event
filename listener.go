package event

// Listener is the interface that must be implemented by event listeners.
type Listener interface {
	// Handle handles the given event.
	//
	// The implementation should return true if the event was handled successfully,
	// or false if there was an error handling the event.
	Handle(Event) bool
}

// ListenerFunc is a function that implements the Listener interface.
type ListenerFunc func(Event) bool

// Handle implements the Listener interface for ListenerFunc.
func (f ListenerFunc) Handle(e Event) bool {
	return f(e)
}

// ListenerPriority represents a listener with a priority.
type ListenerPriority struct {
	Listener Listener
	Priority int
}

// EventListeners represents a collection of listeners for an event.
type EventListeners []ListenerPriority

// Less is used for sorting listeners by priority in descending order.
func (el EventListeners) Less(i, j int) bool {
	return el[i].Priority > el[j].Priority
}

// Swap swaps two listeners.
func (el EventListeners) Swap(i, j int) {
	el[i], el[j] = el[j], el[i]
}

// Len returns the number of listeners.
func (el EventListeners) Len() int {
	return len(el)
}
