package examples

import (
	"testing"

	"github.com/parsilver/event"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockListener is a mock implementation of the Listener interface for testing
type MockListener struct {
	mock.Mock
}

// Handle is the mock implementation of the Handle method
func (m *MockListener) Handle(e event.Event) bool {
	args := m.Called(e)
	return args.Bool(0)
}

// Test case for dispatcher with mocked listeners
func TestDispatcherWithMockedListeners(t *testing.T) {
	// Create a dispatcher
	dispatcher := event.NewDispatcher()
	
	// Create a mock listener
	mockListener := new(MockListener)
	
	// Set up expectations
	mockListener.On("Handle", mock.AnythingOfType("*event.BaseEvent")).Return(true).Once()
	
	// Register the listener
	dispatcher.AddListener("test.event", mockListener)
	
	// Dispatch an event
	e := event.NewEvent("test.event")
	dispatcher.Dispatch(e)
	
	// Verify the expectations were met
	mockListener.AssertExpectations(t)
}

// Test case for event propagation stopping
func TestDispatcherWithPropagationStopping(t *testing.T) {
	// Create a dispatcher
	dispatcher := event.NewDispatcher()
	
	// Create mock listeners
	firstListener := new(MockListener)
	secondListener := new(MockListener)
	
	// First listener will stop propagation
	firstListener.On("Handle", mock.AnythingOfType("*event.BaseEvent")).Run(func(args mock.Arguments) {
		e := args.Get(0).(event.Event)
		e.StopPropagation()
	}).Return(true).Once()
	
	// Second listener should not be called
	// We don't set any expectations on it, so if it's called, the test will fail
	
	// Register the listeners with priorities
	dispatcher.AddListener("test.event", firstListener, 100)   // higher priority, called first
	dispatcher.AddListener("test.event", secondListener, 50)   // lower priority, should not be called
	
	// Dispatch an event
	e := event.NewEvent("test.event")
	dispatcher.Dispatch(e)
	
	// Verify expectations
	firstListener.AssertExpectations(t)
	// Use assert to make sure the second listener wasn't called
	secondListener.AssertNotCalled(t, "Handle", mock.Anything)
	
	// Use assert package to verify event propagation was stopped
	assert.True(t, e.IsPropagationStopped())
}

// Test case for event with arguments
func TestEventWithArguments(t *testing.T) {
	// Create a dispatcher
	dispatcher := event.NewDispatcher()
	
	// Create a mock listener
	mockListener := new(MockListener)
	
	// Define arguments
	args := map[string]interface{}{
		"user_id": 123,
		"email":   "test@example.com",
	}
	
	// Set up expectations to verify arguments
	mockListener.On("Handle", mock.MatchedBy(func(e event.Event) bool {
		// Check that it's the right event name
		if e.Name() != "test.event" {
			return false
		}
		
		// Check that the arguments match
		eventArgs := e.Arguments()
		return eventArgs["user_id"] == args["user_id"] && 
			eventArgs["email"] == args["email"]
	})).Return(true).Once()
	
	// Register the listener
	dispatcher.AddListener("test.event", mockListener)
	
	// Create and dispatch an event with arguments
	e := event.NewEvent("test.event", args)
	dispatcher.Dispatch(e)
	
	// Verify the expectations were met
	mockListener.AssertExpectations(t)
}