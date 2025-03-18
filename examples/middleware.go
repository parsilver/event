package examples

import (
	"fmt"
	"log"
	"time"

	"github.com/parsilver/event"
)

// LoggingMiddleware is a listener that logs all events
type LoggingMiddleware struct{}

func (m *LoggingMiddleware) Handle(e event.Event) bool {
	log.Printf("[EVENT] %s occurred at %s", e.Name(), time.Now().Format(time.RFC3339))
	return true
}

// TimingMiddleware measures how long event processing takes
type TimingMiddleware struct{}

func (m *TimingMiddleware) Handle(e event.Event) bool {
	start := time.Now()
	
	// We return true to continue processing
	// The actual timing will be logged after all listeners have processed
	
	// Store the start time in the event arguments for later retrieval
	args := e.Arguments()
	args["_timing_start"] = start
	
	return true
}

// TimingCleanupMiddleware calculates and logs the time taken
type TimingCleanupMiddleware struct{}

func (m *TimingCleanupMiddleware) Handle(e event.Event) bool {
	args := e.Arguments()
	if startTime, ok := args["_timing_start"].(time.Time); ok {
		duration := time.Since(startTime)
		log.Printf("[TIMING] Event %s took %v to process", e.Name(), duration)
		
		// Clean up our internal timing data
		delete(args, "_timing_start")
	}
	return true
}

// OrderCreatedEvent is a custom event
type OrderCreatedEvent struct {
	*event.BaseEvent
	OrderID    string
	CustomerID string
	Amount     float64
}

// NewOrderCreatedEvent creates a new OrderCreatedEvent
func NewOrderCreatedEvent(orderID, customerID string, amount float64) *OrderCreatedEvent {
	return &OrderCreatedEvent{
		BaseEvent:  event.NewEvent("order.created"),
		OrderID:    orderID,
		CustomerID: customerID,
		Amount:     amount,
	}
}

// OrderProcessor processes orders
type OrderProcessor struct{}

func (p *OrderProcessor) Handle(e event.Event) bool {
	if order, ok := e.(*OrderCreatedEvent); ok {
		// Simulate some processing time
		time.Sleep(50 * time.Millisecond)
		fmt.Printf("Processing order %s for customer %s (Amount: $%.2f)\n",
			order.OrderID, order.CustomerID, order.Amount)
	}
	return true
}

// NotificationSubscriber sends notifications for various events
type NotificationSubscriber struct{}

func (s *NotificationSubscriber) OnOrderCreated(e event.Event) bool {
	if order, ok := e.(*OrderCreatedEvent); ok {
		// Simulate notification sending
		time.Sleep(30 * time.Millisecond)
		fmt.Printf("Sending order confirmation for Order #%s to customer %s\n",
			order.OrderID, order.CustomerID)
	}
	return true
}

func (s *NotificationSubscriber) GetSubscribedEvents() map[string][]event.SubscriberConfig {
	return map[string][]event.SubscriberConfig{
		"order.created": {
			{
				Method:   "OnOrderCreated",
				Priority: 10,
			},
		},
	}
}

func MiddlewareExample() {
	// Create a dispatcher
	dispatcher := event.NewDispatcher()
	
	// Add middleware with very high and very low priorities to wrap all other listeners
	dispatcher.AddListener("order.created", &TimingMiddleware{}, 1000)      // Run first
	dispatcher.AddListener("order.created", &LoggingMiddleware{}, 900)      // Run second
	dispatcher.AddListener("order.created", &TimingCleanupMiddleware{}, -1000) // Run last
	
	// Add business logic listeners
	dispatcher.AddListener("order.created", &OrderProcessor{}, 100)
	
	// Add subscribers
	event.RegisterSubscriber(dispatcher, &NotificationSubscriber{})
	
	// Create and dispatch an event
	orderEvent := NewOrderCreatedEvent("ORD-12345", "CUST-789", 99.99)
	dispatcher.Dispatch(orderEvent)
	
	// Create and dispatch another event
	anotherOrder := NewOrderCreatedEvent("ORD-67890", "CUST-456", 149.99)
	dispatcher.Dispatch(anotherOrder)
}