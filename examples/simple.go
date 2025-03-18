package examples

import (
	"fmt"
	"time"

	"github.com/parsilver/event"
)

// Define a custom event
type UserCreatedEvent struct {
	*event.BaseEvent
	UserID   int
	Username string
	Email    string
}

// NewUserCreatedEvent creates a new UserCreatedEvent.
func NewUserCreatedEvent(userID int, username, email string) *UserCreatedEvent {
	return &UserCreatedEvent{
		BaseEvent: event.NewEvent("user.created"),
		UserID:    userID,
		Username:  username,
		Email:     email,
	}
}

// Simple listener implementation
type UserLogger struct{}

func (l *UserLogger) Handle(e event.Event) bool {
	if userCreated, ok := e.(*UserCreatedEvent); ok {
		fmt.Printf("[%s] User created: %s (ID: %d, Email: %s)\n",
			time.Now().Format(time.RFC3339),
			userCreated.Username,
			userCreated.UserID,
			userCreated.Email)
	}
	return true
}

// EmailSubscriber is a subscriber that listens to multiple events.
type EmailSubscriber struct{}

func (s *EmailSubscriber) OnUserCreated(e event.Event) bool {
	if userCreated, ok := e.(*UserCreatedEvent); ok {
		fmt.Printf("Sending welcome email to %s\n", userCreated.Email)
	}
	return true
}

func (s *EmailSubscriber) GetSubscribedEvents() map[string][]event.SubscriberConfig {
	return map[string][]event.SubscriberConfig{
		"user.created": {
			{
				Method:   "OnUserCreated",
				Priority: 5, // Lower priority, will be called after UserLogger
			},
		},
	}
}

func SimpleExample() {
	// Create a dispatcher
	dispatcher := event.NewDispatcher()

	// Register a simple listener
	userLogger := &UserLogger{}
	dispatcher.AddListener("user.created", userLogger, 10) // Higher priority

	// Register a subscriber
	emailSubscriber := &EmailSubscriber{}
	event.RegisterSubscriber(dispatcher, emailSubscriber)

	// Create and dispatch an event
	userEvent := NewUserCreatedEvent(1, "johndoe", "john@example.com")
	dispatcher.Dispatch(userEvent)

	// Using ListenerFunc for a one-off listener
	dispatcher.AddListener("user.created", event.ListenerFunc(func(e event.Event) bool {
		fmt.Println("User creation process completed")
		return true
	}), 1) // Lowest priority

	// Dispatch another event
	anotherUser := NewUserCreatedEvent(2, "janedoe", "jane@example.com")
	dispatcher.Dispatch(anotherUser)

	// Example of stopping propagation
	stoppingListener := event.ListenerFunc(func(e event.Event) bool {
		fmt.Println("This listener will stop propagation")
		e.StopPropagation()
		return true
	})

	dispatcher.AddListener("user.created", stoppingListener, 15) // Highest priority

	// This event will only trigger the stoppingListener
	thirdUser := NewUserCreatedEvent(3, "bobsmith", "bob@example.com")
	dispatcher.Dispatch(thirdUser)
}