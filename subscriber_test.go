package event_test

import (
	"testing"

	"github.com/parsilver/event"
	"github.com/stretchr/testify/assert"
)

type TestSubscriber struct {
	calledEvents map[string]bool
}

func NewTestSubscriber() *TestSubscriber {
	return &TestSubscriber{
		calledEvents: make(map[string]bool),
	}
}

func (s *TestSubscriber) OnUserCreated(e event.Event) bool {
	s.calledEvents["user.created"] = true
	return true
}

func (s *TestSubscriber) OnUserUpdated(e event.Event) bool {
	s.calledEvents["user.updated"] = true
	return true
}

func (s *TestSubscriber) GetSubscribedEvents() map[string][]event.SubscriberConfig {
	return map[string][]event.SubscriberConfig{
		"user.created": {
			{
				Method:   "OnUserCreated",
				Priority: 100,
			},
		},
		"user.updated": {
			{
				Method:   "OnUserUpdated",
				Priority: 90,
			},
		},
	}
}

func TestSubscriber_Registration(t *testing.T) {
	dispatcher := event.NewDispatcher()
	subscriber := NewTestSubscriber()

	// Register subscriber
	event.RegisterSubscriber(dispatcher, subscriber)

	// Dispatch events
	dispatcher.Dispatch(event.NewEvent("user.created"))
	dispatcher.Dispatch(event.NewEvent("user.updated"))

	// Verify both methods were called
	assert.True(t, subscriber.calledEvents["user.created"])
	assert.True(t, subscriber.calledEvents["user.updated"])
}

func TestSubscriber_PriorityRespected(t *testing.T) {
	dispatcher := event.NewDispatcher()

	var callOrder []string

	// Create listeners directly with different priorities
	listener1 := event.ListenerFunc(func(e event.Event) bool {
		callOrder = append(callOrder, "listener1")
		return true
	})

	listener2 := event.ListenerFunc(func(e event.Event) bool {
		callOrder = append(callOrder, "listener2")
		return true
	})

	// Add listeners with different priorities
	dispatcher.AddListener("test.event", listener1, 10)  // Lower priority
	dispatcher.AddListener("test.event", listener2, 100) // Higher priority

	// Dispatch event
	dispatcher.Dispatch(event.NewEvent("test.event"))

	// Verify call order based on priority (higher priority first)
	assert.Equal(t, []string{"listener2", "listener1"}, callOrder)
}
