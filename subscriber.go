package event

import (
	"fmt"
	"reflect"
)

// SubscriberConfig represents the configuration for a subscriber method.
type SubscriberConfig struct {
	// Method is the name of the method to call on the subscriber.
	Method string
	
	// Priority is the priority of the listener.
	Priority int
}

// Subscriber is the interface that must be implemented by event subscribers.
type Subscriber interface {
	// GetSubscribedEvents returns a map of event names to subscriber configurations.
	GetSubscribedEvents() map[string][]SubscriberConfig
}

// SubscriberFunc is a function type that implements the Subscriber interface.
type SubscriberFunc func() map[string][]SubscriberConfig

// GetSubscribedEvents implements the Subscriber interface for SubscriberFunc.
func (f SubscriberFunc) GetSubscribedEvents() map[string][]SubscriberConfig {
	return f()
}

// RegisterSubscriber registers a subscriber with the dispatcher.
func RegisterSubscriber(dispatcher Dispatcher, subscriber Subscriber) {
	subscribedEvents := subscriber.GetSubscribedEvents()
	
	for eventName, configs := range subscribedEvents {
		for _, config := range configs {
			// Create a listener for each method
			listener := createListenerFromSubscriber(subscriber, config.Method)
			dispatcher.AddListener(eventName, listener, config.Priority)
		}
	}
}

// RegisterListener registers a method on a target object as a listener.
// This is used internally by RegisterSubscriber but can also be used directly.
func RegisterListener(dispatcher Dispatcher, target interface{}, methodName string, handler func(Event) bool) {
	// Create a listener from the handler function
	listener := ListenerFunc(handler)
	
	// Get the subscriber's registered events
	if subscriber, ok := target.(Subscriber); ok {
		subscribedEvents := subscriber.GetSubscribedEvents()
		
		// Find all event names that use this method
		for eventName, configs := range subscribedEvents {
			for _, config := range configs {
				if config.Method == methodName {
					dispatcher.AddListener(eventName, listener, config.Priority)
				}
			}
		}
	}
}

// createListenerFromSubscriber creates a listener from a subscriber object and method name.
func createListenerFromSubscriber(subscriber interface{}, methodName string) Listener {
	return ListenerFunc(func(event Event) bool {
		// Get the method by name using reflection
		method := reflect.ValueOf(subscriber).MethodByName(methodName)
		if !method.IsValid() {
			panic(fmt.Sprintf("Method %s does not exist on subscriber %T", methodName, subscriber))
		}
		
		// Call the method with the event
		results := method.Call([]reflect.Value{reflect.ValueOf(event)})
		
		// Check if the method returned a boolean
		if len(results) > 0 && results[0].Kind() == reflect.Bool {
			return results[0].Bool()
		}
		
		// Default to true if the method doesn't return a boolean
		return true
	})
}