package event

import (
	"sort"
	"sync"
)

// Dispatcher manages event listeners and dispatches events to them.
type Dispatcher interface {
	// AddListener adds a listener for the specified event.
	// The priority determines the order of execution of listeners, higher values mean earlier execution.
	AddListener(eventName string, listener Listener, priority ...int)

	// HasListener checks if a listener is registered for the specified event.
	HasListener(eventName string, listener Listener) bool

	// RemoveListener removes a listener from the specified event.
	RemoveListener(eventName string, listener Listener)

	// Dispatch dispatches an event to all registered listeners.
	Dispatch(event Event) Event
}

// EventDispatcher is the default implementation of the Dispatcher interface.
type EventDispatcher struct {
	listeners map[string]EventListeners
	mu        sync.RWMutex
}

// NewDispatcher creates a new event dispatcher.
func NewDispatcher() *EventDispatcher {
	return &EventDispatcher{
		listeners: make(map[string]EventListeners),
	}
}

// AddListener adds a listener for the specified event.
func (d *EventDispatcher) AddListener(eventName string, listener Listener, priority ...int) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Default priority is 0
	p := 0
	if len(priority) > 0 {
		p = priority[0]
	}

	// Initialize the listener slice if it doesn't exist
	if _, ok := d.listeners[eventName]; !ok {
		d.listeners[eventName] = make(EventListeners, 0)
	}

	// Add the listener with its priority
	d.listeners[eventName] = append(d.listeners[eventName], ListenerPriority{
		Listener: listener,
		Priority: p,
	})

	// Sort listeners by priority (higher first)
	sort.Sort(d.listeners[eventName])
}

// HasListener checks if a listener is registered for the specified event.
func (d *EventDispatcher) HasListener(eventName string, listener Listener) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if eventListeners, ok := d.listeners[eventName]; ok {
		for _, registered := range eventListeners {
			if registered.Listener == listener {
				return true
			}
		}
	}

	return false
}

// RemoveListener removes a listener from the specified event.
func (d *EventDispatcher) RemoveListener(eventName string, listener Listener) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if eventListeners, ok := d.listeners[eventName]; ok {
		newListeners := make(EventListeners, 0, len(eventListeners))

		for _, registered := range eventListeners {
			if registered.Listener != listener {
				newListeners = append(newListeners, registered)
			}
		}

		d.listeners[eventName] = newListeners
	}
}

// Dispatch dispatches an event to all registered listeners.
func (d *EventDispatcher) Dispatch(event Event) Event {
	d.mu.RLock()
	eventListeners, ok := d.listeners[event.Name()]
	d.mu.RUnlock()

	if !ok {
		return event
	}

	// Make a copy to avoid concurrent modification issues
	listenersCopy := make(EventListeners, len(eventListeners))
	copy(listenersCopy, eventListeners)

	// Call each listener in priority order
	for _, l := range listenersCopy {
		l.Listener.Handle(event)

		// Stop if propagation is stopped
		if event.IsPropagationStopped() {
			break
		}
	}

	return event
}
