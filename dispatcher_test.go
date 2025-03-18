package event_test

import (
	"testing"

	"github.com/parsilver/event"
	"github.com/stretchr/testify/assert"
)

type TestListener struct {
	called        bool
	eventReceived event.Event
}

func (l *TestListener) Handle(e event.Event) bool {
	l.called = true
	l.eventReceived = e
	return true
}

func TestDispatcher_AddListener(t *testing.T) {
	dispatcher := event.NewDispatcher()
	listener := &TestListener{}
	
	dispatcher.AddListener("user.created", listener)
	
	assert.True(t, dispatcher.HasListener("user.created", listener))
}

func TestDispatcher_RemoveListener(t *testing.T) {
	dispatcher := event.NewDispatcher()
	listener := &TestListener{}
	
	dispatcher.AddListener("user.created", listener)
	assert.True(t, dispatcher.HasListener("user.created", listener))
	
	dispatcher.RemoveListener("user.created", listener)
	assert.False(t, dispatcher.HasListener("user.created", listener))
}

func TestDispatcher_Dispatch(t *testing.T) {
	dispatcher := event.NewDispatcher()
	listener := &TestListener{}
	
	dispatcher.AddListener("user.created", listener)
	
	e := event.NewEvent("user.created")
	dispatcher.Dispatch(e)
	
	assert.True(t, listener.called)
	assert.Equal(t, e, listener.eventReceived)
}

func TestDispatcher_DispatchStopPropagation(t *testing.T) {
	dispatcher := event.NewDispatcher()
	
	// First listener stops propagation
	listener1 := &TestListener{}
	dispatcher.AddListener("user.created", listener1, 100)
	
	// Second listener should not be called
	listener2 := &TestListener{}
	dispatcher.AddListener("user.created", listener2, 90)
	
	// Create custom listener that stops propagation
	stopListener := event.ListenerFunc(func(e event.Event) bool {
		e.StopPropagation()
		return true
	})
	
	dispatcher.AddListener("user.created", stopListener, 100)
	
	e := event.NewEvent("user.created")
	dispatcher.Dispatch(e)
	
	assert.True(t, listener1.called)
	assert.False(t, listener2.called)
}

func TestDispatcher_Priority(t *testing.T) {
	dispatcher := event.NewDispatcher()
	
	var callOrder []int
	
	listener1 := event.ListenerFunc(func(e event.Event) bool {
		callOrder = append(callOrder, 1)
		return true
	})
	
	listener2 := event.ListenerFunc(func(e event.Event) bool {
		callOrder = append(callOrder, 2)
		return true
	})
	
	listener3 := event.ListenerFunc(func(e event.Event) bool {
		callOrder = append(callOrder, 3)
		return true
	})
	
	// Add listeners with different priorities
	dispatcher.AddListener("test.event", listener1, 10)  // lowest priority, called last
	dispatcher.AddListener("test.event", listener2, 30)  // middle priority
	dispatcher.AddListener("test.event", listener3, 50)  // highest priority, called first
	
	e := event.NewEvent("test.event")
	dispatcher.Dispatch(e)
	
	assert.Equal(t, []int{3, 2, 1}, callOrder)
}